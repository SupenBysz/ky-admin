.PHONY: test-api test-coverage test-all test-report serve-report clean help

# 设置项目根目录
BACKEND_ROOT := /Volumes/DataDocument/CodeSpace/AI\ Solution/kysion.com/backend
TEST_DIR := $(BACKEND_ROOT)/tests
REPORTS_DIR := $(TEST_DIR)/reports
SCRIPTS_DIR := $(BACKEND_ROOT)/scripts

# 测试目标
test-api:
	@echo "运行API测试..."
	@$(SCRIPTS_DIR)/run-test.sh --type=api
	@echo "API测试完成"

test-coverage:
	@echo "运行覆盖率测试..."
	@$(SCRIPTS_DIR)/run-test.sh --type=coverage
	@echo "覆盖率测试完成"

test-all:
	@echo "运行所有测试..."
	@$(SCRIPTS_DIR)/run-test.sh --type=all
	@echo "所有测试完成"

# 测试报告
test-report:
	@echo "生成并查看测试报告..."
	@$(SCRIPTS_DIR)/run-test.sh --type=coverage --view
	@echo "测试报告已生成并打开"

# 查看测试报告
serve-report:
	@echo "启动测试报告服务器..."
	@cd $(BACKEND_ROOT) && $(SCRIPTS_DIR)/test.sh coverage serve

# 编译测试工具
build-tools:
	@echo "编译测试工具..."
	@cd $(SCRIPTS_DIR)/tools && go build -o ../bin/test-report test_report_tool.go

# 清理测试报告
clean:
	@echo "清理测试报告..."
	@rm -rf $(REPORTS_DIR)/*

# 帮助信息
help:
	@echo "测试命令:"
	@echo "  make test-api        - 执行API测试"
	@echo "  make test-coverage   - 执行覆盖率测试"
	@echo "  make test-all        - 执行所有测试"
	@echo "  make test-report     - 生成并查看测试报告"
	@echo "  make serve-report    - 启动测试报告服务器"
	@echo "  make build-tools     - 编译测试工具"
	@echo "  make clean           - 清理测试报告"
	@echo "  make help            - 显示帮助信息"
	@echo ""
	@echo "或者直接使用交互式脚本:"
	@echo "  $(SCRIPTS_DIR)/run-test.sh"

# 默认目标
default: help 