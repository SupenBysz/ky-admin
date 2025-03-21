#!/bin/bash

set -e

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 获取脚本所在的真实目录，无论从哪里调用
SCRIPT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
KY_ADMIN_DIR="$PROJECT_ROOT/ky-admin"
TEST_DIR="$PROJECT_ROOT/tests"
TEST_REPORT_DIR="$TEST_DIR/reports"

# 打印当前路径信息
echo "脚本目录: $SCRIPT_DIR"
echo "项目根目录: $PROJECT_ROOT"
echo "KY-Admin目录: $KY_ADMIN_DIR"
echo "测试目录: $TEST_DIR"
echo "报告目录: $TEST_REPORT_DIR"

# 切换到项目根目录
cd "$PROJECT_ROOT"

# 创建测试报告目录（如果不存在）
mkdir -p "$TEST_REPORT_DIR"

# 测试类型和操作
TEST_TYPE=$1
OPERATION=$2
COVERAGE_THRESHOLD=80
TEST_LOG_FILE="test_output.log"

# 检查测试类型是否有效
is_valid_test_type() {
    local type=$1
    case $type in
        basic|full|coverage|api|integration|performance|all)
            return 0;;
        *)
            return 1;;
    esac
}

# 检查操作是否有效
is_valid_operation() {
    local op=$1
    case $op in
        report|serve|"")
            return 0;;
        *)
            return 1;;
    esac
}

# 检查测试覆盖率是否达到要求
check_coverage() {
    local coverage_file=$1
    local threshold=$2
    
    # 提取总覆盖率百分比
    local coverage=$(grep -o "total:.*statements" "$coverage_file" | grep -o "[0-9]\+\.[0-9]\+")
    
    echo -e "${YELLOW}测试覆盖率: $coverage%${NC}"
    
    # 比较覆盖率与阈值
    if (( $(echo "$coverage < $threshold" | bc -l) )); then
        echo -e "${RED}❌ 测试覆盖率不足 $threshold%${NC}"
        return 1
    else
        echo -e "${GREEN}✅ 测试覆盖率达到 $threshold%${NC}"
        return 0
    fi
}

# 运行测试并生成覆盖率报告
run_coverage_test() {
    local pkgs=$1
    local output_dir=$2
    
    # 创建输出目录
    mkdir -p "$output_dir"
    
    echo -e "${YELLOW}正在运行测试并生成覆盖率报告...${NC}"
    
    # 运行测试并生成覆盖率数据，保存完整输出到日志文件
    go test $pkgs -coverprofile="$output_dir/coverage.out" -v | tee "$output_dir/$TEST_LOG_FILE"
    
    # 生成覆盖率HTML报告
    go tool cover -html="$output_dir/coverage.out" -o "$output_dir/coverage.html"
    
    # 生成覆盖率文本报告
    go tool cover -func="$output_dir/coverage.out" > "$output_dir/coverage.txt"
    
    # 生成测试摘要报告
    echo -e "${YELLOW}正在生成测试摘要报告...${NC}"
    go run "$SCRIPT_DIR/test_report_tool.go" -mode=summary -input "$output_dir/$TEST_LOG_FILE" -dir "$output_dir"
    
    echo -e "${GREEN}报告已生成:${NC}"
    echo -e "  ${YELLOW}覆盖率报告:${NC} $output_dir/coverage.html"
    echo -e "  ${YELLOW}测试摘要:${NC} $output_dir/test_summary.html"
    
    # 检查覆盖率
    check_coverage "$output_dir/coverage.txt" $COVERAGE_THRESHOLD
    return $?
}

# 启动测试报告服务器
start_report_server() {
    echo -e "${YELLOW}启动测试报告服务器...${NC}"
    go run "$SCRIPT_DIR/test_report_tool.go" -mode=server -dir="$TEST_REPORT_DIR"
}

# 显示帮助信息
show_help() {
    echo "使用: $0 [测试类型] [操作]"
    echo ""
    echo "测试类型:"
    echo "  basic       - 运行基本单元测试"
    echo "  full        - 运行所有单元测试"
    echo "  coverage    - 运行测试并生成覆盖率报告"
    echo "  api         - 运行API测试"
    echo "  integration - 运行集成测试"
    echo "  performance - 运行性能测试"
    echo "  all         - 运行所有测试"
    echo ""
    echo "操作:"
    echo "  report      - 只生成报告，不运行测试 (与 coverage 类型一起使用)"
    echo "  serve       - 启动测试报告服务器"
    echo ""
    echo "示例:"
    echo "  $0 basic                - 运行基本单元测试"
    echo "  $0 coverage             - 生成测试覆盖率报告"
    echo "  $0 all                  - 运行所有测试并生成报告"
    echo "  $0 coverage report      - 使用最新的测试输出生成报告"
    echo "  $0 coverage serve       - 启动测试报告服务器"
}

# 如果没有参数或参数无效，显示帮助信息
if [ -z "$TEST_TYPE" ] || ! is_valid_test_type "$TEST_TYPE" || ! is_valid_operation "$OPERATION"; then
    show_help
    exit 1
fi

# 如果是启动报告服务器
if [ "$OPERATION" = "serve" ]; then
    start_report_server
    exit 0
fi

echo "开始测试..."

# 根据测试类型和操作运行不同的测试
case $TEST_TYPE in
    basic)
        echo "运行基本测试..."
        REPORT_DIR="$TEST_REPORT_DIR/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$REPORT_DIR"
        
        # 保存所有输出到日志文件
        {
            echo "运行测试: 配置模块"
            go test ./pkg/config -v || { echo -e "${RED}✗ 失败: 测试 配置模块${NC}"; exit 1; }
            
            echo "运行测试: 健康检查API"
            go test ./test/api/health_test.go -v || { echo -e "${RED}✗ 失败: 测试 健康检查API${NC}"; exit 1; }
            
            echo -e "${GREEN}✓ 成功: 所有基本测试通过${NC}"
        } | tee "$REPORT_DIR/$TEST_LOG_FILE"
        
        # 生成测试摘要报告
        echo -e "${YELLOW}正在生成测试摘要报告...${NC}"
        go run "$SCRIPT_DIR/test_report_tool.go" -mode=summary -input "$REPORT_DIR/$TEST_LOG_FILE" -dir "$REPORT_DIR" || true
        
        echo -e "${GREEN}报告已生成:${NC} $REPORT_DIR/test_summary.html"
        ;;
        
    full)
        echo "运行所有单元测试..."
        REPORT_DIR="$TEST_REPORT_DIR/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$REPORT_DIR"
        
        # 保存输出到日志文件
        go test ./... -v | tee "$REPORT_DIR/$TEST_LOG_FILE" || { echo -e "${RED}✗ 失败: 单元测试${NC}"; exit 1; }
        
        # 生成测试摘要报告
        echo -e "${YELLOW}正在生成测试摘要报告...${NC}"
        go run "$SCRIPT_DIR/test_report_tool.go" -mode=summary -input "$REPORT_DIR/$TEST_LOG_FILE" -dir "$REPORT_DIR" || true
        
        echo -e "${GREEN}✓ 成功: 所有单元测试通过${NC}"
        echo -e "${GREEN}报告已生成:${NC} $REPORT_DIR/test_summary.html"
        ;;
        
    coverage)
        if [ "$OPERATION" = "report" ]; then
            echo "使用最新的测试输出生成报告..."
            LATEST_DIR=$(ls -td "$TEST_REPORT_DIR/"* | head -1)
            
            if [ -z "$LATEST_DIR" ]; then
                echo -e "${RED}✗ 失败: 没有找到测试报告目录${NC}"
                exit 1
            fi
            
            if [ -f "$LATEST_DIR/$TEST_LOG_FILE" ]; then
                # 生成测试摘要报告
                echo -e "${YELLOW}正在生成测试摘要报告...${NC}"
                go run "$SCRIPT_DIR/test_report_tool.go" -mode=summary -input "$LATEST_DIR/$TEST_LOG_FILE" -dir "$LATEST_DIR" || true
                
                echo -e "${GREEN}报告已生成:${NC}"
                echo -e "  ${YELLOW}覆盖率报告:${NC} $LATEST_DIR/coverage.html"
                echo -e "  ${YELLOW}测试摘要:${NC} $LATEST_DIR/test_summary.html"
            else
                echo -e "${RED}✗ 失败: 没有找到测试日志文件${NC}"
                exit 1
            fi
        else
            echo "生成测试覆盖率报告..."
            # 创建报告目录
            REPORT_DIR="$TEST_REPORT_DIR/$(date +%Y%m%d_%H%M%S)"
            mkdir -p "$REPORT_DIR"
            
            # 运行覆盖率测试
            run_coverage_test "./..." "$REPORT_DIR"
            if [ $? -ne 0 ]; then
                echo -e "${RED}✗ 失败: 测试覆盖率未达标${NC}"
                exit 1
            fi
            echo -e "${GREEN}✓ 成功: 测试覆盖率达标${NC}"
        fi
        ;;
        
    api)
        echo "运行API测试..."
        # 指定完整路径
        REPORT_DIR="${TEST_REPORT_DIR}/$(date +%Y%m%d_%H%M%S)"
        # 确保目录存在
        mkdir -p "${REPORT_DIR}"
        
        # 保存输出到日志文件，使用完整路径
        go test ./test/api -v | tee "${REPORT_DIR}/${TEST_LOG_FILE}" || { echo -e "${RED}✗ 失败: API测试${NC}"; exit 1; }
        
        # 确认日志文件已创建
        if [ -f "${REPORT_DIR}/${TEST_LOG_FILE}" ]; then
            # 生成测试摘要报告
            echo -e "${YELLOW}正在生成测试摘要报告...${NC}"
            go run "$SCRIPT_DIR/test_report_tool.go" -mode=summary -input "${REPORT_DIR}/${TEST_LOG_FILE}" -dir "${REPORT_DIR}" || true
            
            echo -e "${GREEN}✓ 成功: 所有API测试通过${NC}"
            echo -e "${GREEN}报告已生成:${NC} ${REPORT_DIR}/test_summary.html"
        else
            echo -e "${YELLOW}警告: 未能创建日志文件，可能是路径问题${NC}"
            echo -e "${GREEN}✓ 成功: 所有API测试通过${NC}"
        fi
        ;;
        
    integration)
        echo "运行集成测试..."
        REPORT_DIR="$TEST_REPORT_DIR/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$REPORT_DIR"
        
        # 保存输出到日志文件
        go test ./test/integration -v | tee "$REPORT_DIR/$TEST_LOG_FILE" || { echo -e "${RED}✗ 失败: 集成测试${NC}"; exit 1; }
        
        # 生成测试摘要报告
        echo -e "${YELLOW}正在生成测试摘要报告...${NC}"
        go run "$SCRIPT_DIR/test_report_tool.go" -mode=summary -input "$REPORT_DIR/$TEST_LOG_FILE" -dir "$REPORT_DIR" || true
        
        echo -e "${GREEN}✓ 成功: 所有集成测试通过${NC}"
        echo -e "${GREEN}报告已生成:${NC} $REPORT_DIR/test_summary.html"
        ;;
        
    performance)
        echo "运行性能测试..."
        REPORT_DIR="$TEST_REPORT_DIR/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$REPORT_DIR"
        
        # 保存输出到日志文件
        go test ./test/performance -v -bench=. | tee "$REPORT_DIR/$TEST_LOG_FILE" || { echo -e "${RED}✗ 失败: 性能测试${NC}"; exit 1; }
        
        # 生成测试摘要报告
        echo -e "${YELLOW}正在生成测试摘要报告...${NC}"
        go run "$SCRIPT_DIR/test_report_tool.go" -mode=summary -input "$REPORT_DIR/$TEST_LOG_FILE" -dir "$REPORT_DIR" || true
        
        echo -e "${GREEN}✓ 成功: 所有性能测试通过${NC}"
        echo -e "${GREEN}报告已生成:${NC} $REPORT_DIR/test_summary.html"
        ;;
        
    all)
        echo "运行所有测试..."
        
        # 创建报告目录
        REPORT_DIR="$TEST_REPORT_DIR/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$REPORT_DIR"
        
        {
            # 运行单元测试
            echo "1. 运行单元测试"
            go test ./... -v || { echo -e "${RED}✗ 失败: 单元测试${NC}"; exit 1; }
            
            # 运行API测试
            echo "2. 运行API测试"
            go test ./test/api -v || { echo -e "${RED}✗ 失败: API测试${NC}"; exit 1; }
            
            # 运行集成测试
            echo "3. 运行集成测试"
            go test ./test/integration -v || { echo -e "${RED}✗ 失败: 集成测试${NC}"; exit 1; }
            
            # 运行性能测试
            echo "4. 运行性能测试"
            go test ./test/performance -v -bench=. || { echo -e "${RED}✗ 失败: 性能测试${NC}"; exit 1; }
            
            # 生成测试覆盖率报告
            echo "5. 生成测试覆盖率报告"
            go test ./... -coverprofile=$REPORT_DIR/coverage.out
            go tool cover -html=$REPORT_DIR/coverage.out -o $REPORT_DIR/coverage.html
            go tool cover -func=$REPORT_DIR/coverage.out > $REPORT_DIR/coverage.txt
        } | tee $REPORT_DIR/$TEST_LOG_FILE
        
        # 生成测试摘要报告
        echo -e "${YELLOW}正在生成测试摘要报告...${NC}"
        go run "$SCRIPT_DIR/test_report_tool.go" -mode=summary -input "$REPORT_DIR/$TEST_LOG_FILE" -dir "$REPORT_DIR" || true
        
        # 检查覆盖率
        check_coverage "$REPORT_DIR/coverage.txt" $COVERAGE_THRESHOLD
        
        echo -e "${GREEN}✓ 成功: 所有测试通过${NC}"
        echo -e "${GREEN}报告已生成:${NC}"
        echo -e "  ${YELLOW}覆盖率报告:${NC} $REPORT_DIR/coverage.html"
        echo -e "  ${YELLOW}测试摘要:${NC} $REPORT_DIR/test_summary.html"
        ;;
esac

echo "测试完成."
echo -e "${YELLOW}提示: 运行 '$0 coverage serve' 可启动测试报告服务器查看测试报告。${NC}"
exit 0
