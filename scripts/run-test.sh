#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BACKEND_ROOT="$PROJECT_ROOT"
TEST_DIR="$PROJECT_ROOT/tests"
REPORTS_DIR="$TEST_DIR/reports"
SCRIPTS_DIR="$PROJECT_ROOT/scripts"

# 确保报告目录存在
mkdir -p "$REPORTS_DIR"

# 测试类型
TEST_TYPE="api"
VIEW_REPORT=false
OPERATION="test"

# 显示标题
show_title() {
    clear
    echo -e "${BLUE}===================================================${NC}"
    echo -e "${CYAN}               测试执行工具 v1.0                   ${NC}"
    echo -e "${BLUE}===================================================${NC}"
    echo ""
    echo -e "${YELLOW}当前选择:${NC}"
    echo -e "  测试类型: ${GREEN}$TEST_TYPE${NC}"
    echo -e "  查看报告: ${VIEW_REPORT_TEXT}${NC}"
    echo ""
}

# 显示菜单
show_menu() {
    VIEW_REPORT_TEXT="${RED}否${NC}"
    if [ "$VIEW_REPORT" = true ]; then
        VIEW_REPORT_TEXT="${GREEN}是${NC}"
    fi
    
    show_title
    
    echo -e "${CYAN}请选择操作:${NC}"
    echo -e "${YELLOW}1)${NC} 选择测试类型"
    echo -e "${YELLOW}2)${NC} 切换是否查看报告 (当前: $VIEW_REPORT_TEXT)"
    echo -e "${YELLOW}3)${NC} 执行测试"
    echo -e "${YELLOW}4)${NC} 仅查看报告"
    echo -e "${YELLOW}0)${NC} 退出"
    echo ""
    echo -e "${BLUE}===================================================${NC}"
    echo -e "${YELLOW}输入选项编号:${NC} "
}

# 选择测试类型
select_test_type() {
    show_title
    echo -e "${CYAN}请选择测试类型:${NC}"
    echo -e "${YELLOW}1)${NC} API测试"
    echo -e "${YELLOW}2)${NC} 单元测试"
    echo -e "${YELLOW}3)${NC} 完整测试"
    echo -e "${YELLOW}4)${NC} 覆盖率测试"
    echo -e "${YELLOW}5)${NC} 集成测试"
    echo -e "${YELLOW}6)${NC} 性能测试"
    echo -e "${YELLOW}7)${NC} 所有测试"
    echo -e "${YELLOW}0)${NC} 返回主菜单"
    echo ""
    echo -e "${BLUE}===================================================${NC}"
    echo -e "${YELLOW}输入选项编号:${NC} "
    
    read -r option
    case $option in
        1) TEST_TYPE="api" ;;
        2) TEST_TYPE="basic" ;;
        3) TEST_TYPE="full" ;;
        4) TEST_TYPE="coverage" ;;
        5) TEST_TYPE="integration" ;;
        6) TEST_TYPE="performance" ;;
        7) TEST_TYPE="all" ;;
        0) return ;;
        *) echo -e "${RED}无效选项!${NC}"; sleep 1 ;;
    esac
}

# 切换是否查看报告
toggle_view_report() {
    if [ "$VIEW_REPORT" = true ]; then
        VIEW_REPORT=false
    else
        VIEW_REPORT=true
    fi
}

# 运行测试
run_test() {
    local timestamp=$(date "+%Y%m%d_%H%M%S")
    local report_dir="$REPORTS_DIR/$timestamp"
    
    # 创建报告目录
    mkdir -p "$report_dir"
    
    echo -e "${CYAN}正在执行 $TEST_TYPE 测试...${NC}"
    
    # 运行测试并捕获输出
    cd "$BACKEND_ROOT"
    "$SCRIPTS_DIR/test.sh" "$TEST_TYPE" > "$report_dir/test_output.log" 2>&1
    
    local status=$?
    
    # 创建latest链接 (在报告目录内部使用相对路径)
    cd "$REPORTS_DIR" && ln -sf $(basename "$report_dir") latest && cd "$BACKEND_ROOT"
    
    if [ $status -eq 0 ]; then
        echo -e "${GREEN}测试成功完成!${NC}"
    else
        echo -e "${RED}测试执行失败，查看日志获取详细信息。${NC}"
        cat "$report_dir/test_output.log"
    fi
    
    # 根据用户选择是否查看报告
    if [ "$VIEW_REPORT" = true ]; then
        view_report
    else
        echo -e "${YELLOW}按Enter键继续...${NC}"
        read -r
    fi
}

# 启动报告服务器
start_report_server() {
    cd "$BACKEND_ROOT"
    "$SCRIPTS_DIR/test.sh" coverage serve > /dev/null 2>&1 &
    SERVER_PID=$!
    
    echo -e "${GREEN}报告服务器已启动 (PID: $SERVER_PID)${NC}"
    
    # 等待服务器启动
    sleep 2
}

# 查看测试报告
view_report() {
    # 检查报告服务器是否已启动
    if [ -z "$SERVER_PID" ] || ! ps -p $SERVER_PID > /dev/null; then
        echo -e "${YELLOW}启动报告服务器...${NC}"
        start_report_server
    fi
    
    # 打开浏览器查看报告
    echo -e "${CYAN}打开浏览器查看报告...${NC}"
    
    # 根据操作系统选择打开命令
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # MacOS
        open "http://localhost:8089"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        xdg-open "http://localhost:8089" || firefox "http://localhost:8089" || google-chrome "http://localhost:8089"
    else
        echo -e "${YELLOW}请手动打开浏览器访问: http://localhost:8089${NC}"
    fi
    
    echo -e "${YELLOW}按Enter键继续...${NC}"
    read -r
}

# 清理函数
cleanup() {
    echo -e "\n${YELLOW}正在清理...${NC}"
    if [ ! -z "$SERVER_PID" ]; then
        kill $SERVER_PID > /dev/null 2>&1
        echo -e "${GREEN}报告服务器已关闭 (PID: $SERVER_PID)${NC}"
    fi
    exit 0
}

# 注册清理函数
trap cleanup INT TERM

# 处理命令行参数
parse_args() {
    for arg in "$@"; do
        case $arg in
            --type=*)
                TEST_TYPE="${arg#*=}"
                ;;
            --view)
                VIEW_REPORT=true
                ;;
            --help)
                echo "使用方法: ./run-test.sh [选项]"
                echo "选项:"
                echo "  --type=TYPE    设置测试类型 (api, basic, full, coverage, integration, performance, all)"
                echo "  --view         执行测试后查看报告"
                echo "  --help         显示此帮助信息"
                exit 0
                ;;
        esac
    done
}

# 主函数
main() {
    # 处理命令行参数
    if [ $# -gt 0 ]; then
        parse_args "$@"
        run_test
        cleanup
    fi
    
    # 交互模式
    while true; do
        show_menu
        read -r option
        
        case $option in
            1) select_test_type ;;
            2) toggle_view_report ;;
            3) run_test ;;
            4) view_report ;;
            0) cleanup ;;
            *) echo -e "${RED}无效选项!${NC}"; sleep 1 ;;
        esac
    done
}

# 执行主函数
main "$@" 