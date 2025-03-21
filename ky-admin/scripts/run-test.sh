#!/bin/bash

# 脚本设置
set -e

# 目录定义
BACKEND_ROOT="/Volumes/DataDocument/CodeSpace/AI Solution/kysion.com/backend"
PROJECT_ROOT="$BACKEND_ROOT/ky-admin"
REPORT_ROOT="$BACKEND_ROOT/temp/test-reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_DIR="$REPORT_ROOT/$TIMESTAMP"

# 创建报告目录
mkdir -p "$REPORT_DIR"

# 测试类型
test_type=""
view_report=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --type=*)
      test_type="${key#*=}"
      shift
      ;;
    --view)
      view_report=true
      shift
      ;;
    *)
      echo "未知选项: $key"
      exit 1
      ;;
  esac
done

# 如果没有指定测试类型，显示交互式菜单
if [ -z "$test_type" ]; then
  echo "请选择测试类型:"
  echo "1) API测试"
  echo "2) 覆盖率测试"
  echo "3) 所有测试"
  read -p "选择 [1-3]: " choice
  case $choice in
    1) test_type="api" ;;
    2) test_type="coverage" ;;
    3) test_type="all" ;;
    *) echo "无效选择"; exit 1 ;;
  esac
fi

# 显示测试信息
echo "============================================"
echo "开始执行测试 - $(date)"
echo "测试类型: $test_type"
echo "报告目录: $REPORT_DIR"
echo "============================================"

# 执行测试
cd "$PROJECT_ROOT"

case $test_type in
  "api")
    echo "执行API测试..."
    go test ./tests/api -v | tee "$REPORT_DIR/api_test.log"
    ;;
  "coverage")
    echo "执行覆盖率测试..."
    go test ./... -coverprofile="$REPORT_DIR/coverage.out" -covermode=atomic
    go tool cover -html="$REPORT_DIR/coverage.out" -o "$REPORT_DIR/coverage.html"
    go tool cover -func="$REPORT_DIR/coverage.out" | tee "$REPORT_DIR/coverage.txt"
    ;;
  "all")
    echo "执行单元测试..."
    go test ./... -v | tee "$REPORT_DIR/unit_test.log"
    
    echo "执行API测试..."
    go test ./tests/api -v | tee -a "$REPORT_DIR/api_test.log"
    
    echo "生成覆盖率报告..."
    go test ./... -coverprofile="$REPORT_DIR/coverage.out" -covermode=atomic
    go tool cover -html="$REPORT_DIR/coverage.out" -o "$REPORT_DIR/coverage.html"
    go tool cover -func="$REPORT_DIR/coverage.out" | tee "$REPORT_DIR/coverage.txt"
    ;;
  *)
    echo "未知测试类型: $test_type"
    exit 1
    ;;
esac

# 创建最新报告的符号链接
ln -sf "$REPORT_DIR" "$REPORT_ROOT/latest"

# 生成测试摘要
echo "============================================"
echo "测试完成 - $(date)"
echo "测试结果摘要:"

if [ -f "$REPORT_DIR/coverage.txt" ]; then
  total_coverage=$(grep "total:" "$REPORT_DIR/coverage.txt" | awk '{print $3}')
  echo "代码覆盖率: $total_coverage"
fi

if [ -f "$REPORT_DIR/unit_test.log" ]; then
  passed=$(grep -c "PASS" "$REPORT_DIR/unit_test.log" || echo "0")
  failed=$(grep -c "FAIL" "$REPORT_DIR/unit_test.log" || echo "0")
  echo "单元测试: 通过 $passed, 失败 $failed"
fi

if [ -f "$REPORT_DIR/api_test.log" ]; then
  passed=$(grep -c "PASS" "$REPORT_DIR/api_test.log" || echo "0")
  failed=$(grep -c "FAIL" "$REPORT_DIR/api_test.log" || echo "0")
  echo "API测试: 通过 $passed, 失败 $failed"
fi

echo "报告目录: $REPORT_DIR"
echo "============================================"

# 如果指定了查看报告，并且生成了HTML报告，则打开它
if [ "$view_report" = true ] && [ -f "$REPORT_DIR/coverage.html" ]; then
  echo "打开覆盖率报告..."
  open "$REPORT_DIR/coverage.html"
fi

# 退出
exit 0 