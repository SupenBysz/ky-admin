#!/bin/bash

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无颜色

# 项目根目录
PROJECT_ROOT=$(pwd)

# 函数：打印带颜色的消息
print_message() {
  local color=$1
  local message=$2
  echo -e "${color}${message}${NC}"
}

# 函数：检查上一个命令的执行结果
check_result() {
  if [ $? -ne 0 ]; then
    print_message "$RED" "✗ 失败: $1"
    exit 1
  else
    print_message "$GREEN" "✓ 成功: $1"
  fi
}

# 函数：运行测试
run_tests() {
  local test_path=$1
  local test_name=$2
  local coverage_file="coverage_$(echo $test_name | tr '/' '_').out"
  
  print_message "$BLUE" "运行测试: $test_name"
  
  go test -v -cover -coverprofile=$coverage_file $test_path
  check_result "测试 $test_name"
  
  go tool cover -html=$coverage_file -o coverage_$(echo $test_name | tr '/' '_').html
  check_result "生成覆盖率报告"
}

# 函数：运行基本测试
run_basic_tests() {
  print_message "$YELLOW" "运行基本测试..."
  
  # 测试配置模块
  run_tests "./pkg/config" "配置模块"
  
  # 测试日志模块
  run_tests "./pkg/logger" "日志模块"
  
  # 测试数据库模块
  run_tests "./pkg/database" "数据库模块"
}

# 函数：运行API测试
run_api_tests() {
  print_message "$YELLOW" "运行API测试..."
  
  # 测试API响应模块
  run_tests "./pkg/api/response" "API响应模块"
  
  # 测试健康检查
  run_tests "./test/api" "API健康检查"
}

# 函数：运行集成测试
run_integration_tests() {
  print_message "$YELLOW" "运行集成测试..."
  
  # 启动集成测试环境
  print_message "$BLUE" "启动集成测试环境..."
  docker-compose -f scripts/docker-compose.test.yml up -d
  check_result "启动测试容器"
  
  # 等待服务就绪
  sleep 5
  
  # 运行集成测试
  run_tests "./test/integration" "集成测试"
  
  # 关闭集成测试环境
  print_message "$BLUE" "关闭集成测试环境..."
  docker-compose -f scripts/docker-compose.test.yml down
  check_result "关闭测试容器"
}

# 函数：生成总覆盖率报告
generate_coverage_report() {
  print_message "$YELLOW" "生成总覆盖率报告..."
  
  # 合并覆盖率文件
  echo "mode: set" > coverage.out
  grep -h -v "mode: set" coverage_*.out >> coverage.out
  check_result "合并覆盖率文件"
  
  # 生成HTML报告
  go tool cover -html=coverage.out -o coverage.html
  check_result "生成HTML覆盖率报告"
  
  # 计算总覆盖率
  COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
  print_message "$GREEN" "总测试覆盖率: $COVERAGE"
}

# 主函数
main() {
  print_message "$YELLOW" "开始测试..."
  
  # 清理旧的覆盖率文件
  rm -f coverage*.out coverage*.html
  
  # 根据参数选择测试类型
  case "$1" in
    "basic")
      run_basic_tests
      ;;
    "api")
      run_api_tests
      ;;
    "integration")
      run_integration_tests
      ;;
    *)
      run_basic_tests
      run_api_tests
      run_integration_tests
      ;;
  esac
  
  # 生成总覆盖率报告
  generate_coverage_report
  
  print_message "$GREEN" "测试完成!"
}

# 执行主函数
main "$@"
