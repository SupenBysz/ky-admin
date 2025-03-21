.PHONY: build run clean test lint help test-api test-coverage test-all test-report

# 项目名称和主包位置
PROJECT_NAME := ky-admin
MAIN_PKG := ./cmd/main.go

# 版本号和构建信息
VERSION := 0.1.0
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建输出目录
BUILD_DIR := ../../artifacts/bin

# Go命令
GO := go
GOBUILD := $(GO) build
GORUN := $(GO) run
GOTEST := $(GO) test
GOMOD := $(GO) mod
GOLINT := golangci-lint

# 构建参数
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

help:
	@echo "Makefile for $(PROJECT_NAME)"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - 构建项目"
	@echo "  make run      - 运行项目"
	@echo "  make clean    - 清理构建文件"
	@echo "  make test     - 运行测试"
	@echo "  make lint     - 代码规范检查"
	@echo "  make mod      - 更新依赖"
	@echo "  make help     - 显示帮助信息"

build:
	@echo "Building $(PROJECT_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) $(MAIN_PKG)
	@echo "Build complete. Binary at $(BUILD_DIR)/$(PROJECT_NAME)"

run:
	@echo "Running $(PROJECT_NAME)..."
	$(GORUN) $(MAIN_PKG)

clean:
	@echo "Cleaning build files..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.txt
	@echo "Clean complete"

test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "Tests complete"

# 测试相关目标
test-api:
	@echo "运行API测试..."
	@cd tests && make test-api

test-coverage:
	@echo "运行覆盖率测试..."
	@cd tests && make test-coverage

test-all:
	@echo "运行所有测试..."
	@cd tests && make test-all

test-report:
	@echo "生成并查看测试报告..."
	@cd tests && make test-report

lint:
	@echo "Running linter..."
	$(GOLINT) run ./...
	@echo "Lint complete"

mod:
	@echo "Updating dependencies..."
	$(GOMOD) tidy
	@echo "Dependencies updated"

# 设置默认目标
default: help 