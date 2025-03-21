#!/bin/bash

# 测试辅助脚本
# 用法: ./test.sh [命令] [参数]

# 目录定义
BACKEND_ROOT="/Volumes/DataDocument/CodeSpace/AI Solution/kysion.com/backend"
PROJECT_ROOT="$BACKEND_ROOT/ky-admin"
REPORT_ROOT="$BACKEND_ROOT/temp/test-reports"

# 命令和参数
cmd=$1
subCmd=$2

# 显示使用帮助
show_help() {
  echo "测试辅助脚本"
  echo "用法: ./test.sh [命令] [参数]"
  echo ""
  echo "可用命令:"
  echo "  coverage [serve|open|latest] - 覆盖率报告操作"
  echo "  report  [serve|open|latest]  - 测试报告操作"
  echo "  clean   [all|reports]        - 清理测试数据"
  echo "  help                         - 显示帮助信息"
  echo ""
  echo "示例:"
  echo "  ./test.sh coverage serve     - 启动覆盖率报告服务器"
  echo "  ./test.sh coverage open      - 打开最新覆盖率报告"
  echo "  ./test.sh coverage latest    - 显示最新覆盖率报告路径"
  echo "  ./test.sh clean reports      - 清理所有测试报告"
}

# 启动覆盖率报告服务器
serve_coverage() {
  if [ ! -d "$REPORT_ROOT/latest" ]; then
    echo "错误: 找不到测试报告目录 $REPORT_ROOT/latest"
    return 1
  fi
  
  echo "启动覆盖率报告服务器..."
  echo "报告URL: http://localhost:8089/"
  cd "$REPORT_ROOT/latest" && python3 -m http.server 8089
}

# 打开最新的覆盖率报告
open_coverage() {
  if [ ! -f "$REPORT_ROOT/latest/coverage.html" ]; then
    echo "错误: 找不到覆盖率报告 $REPORT_ROOT/latest/coverage.html"
    return 1
  }
  
  echo "打开覆盖率报告..."
  open "$REPORT_ROOT/latest/coverage.html"
}

# 获取最新覆盖率报告路径
latest_coverage() {
  if [ ! -f "$REPORT_ROOT/latest/coverage.html" ]; then
    echo "错误: 找不到覆盖率报告 $REPORT_ROOT/latest/coverage.html"
    return 1
  }
  
  echo "最新覆盖率报告: $REPORT_ROOT/latest/coverage.html"
}

# 清理测试报告
clean_reports() {
  echo "清理测试报告..."
  rm -rf "$REPORT_ROOT"/*
  mkdir -p "$REPORT_ROOT"
  echo "测试报告已清理"
}

# 清理所有测试数据
clean_all() {
  echo "清理所有测试数据..."
  clean_reports
  echo "所有测试数据已清理"
}

# 主逻辑
case $cmd in
  "coverage")
    case $subCmd in
      "serve") serve_coverage ;;
      "open") open_coverage ;;
      "latest") latest_coverage ;;
      *) echo "未知的覆盖率命令: $subCmd"; show_help; exit 1 ;;
    esac
    ;;
  "report")
    case $subCmd in
      "serve") serve_coverage ;;  # 目前复用覆盖率报告服务器
      "open") open_coverage ;;    # 目前复用覆盖率报告打开
      "latest") latest_coverage ;;  # 目前复用覆盖率报告路径
      *) echo "未知的报告命令: $subCmd"; show_help; exit 1 ;;
    esac
    ;;
  "clean")
    case $subCmd in
      "reports") clean_reports ;;
      "all") clean_all ;;
      *) echo "未知的清理命令: $subCmd"; show_help; exit 1 ;;
    esac
    ;;
  "help")
    show_help
    ;;
  *)
    echo "未知命令: $cmd"
    show_help
    exit 1
    ;;
esac

exit 0 