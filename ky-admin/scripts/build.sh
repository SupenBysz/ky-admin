#!/bin/bash

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
ARTIFACTS_DIR="../../artifacts"

# 创建构建目录
mkdir -p "$PROJECT_ROOT/$ARTIFACTS_DIR/bin"
mkdir -p "$PROJECT_ROOT/$ARTIFACTS_DIR/logs"

echo "开始构建 ky-admin..."

# 编译项目
cd "$PROJECT_ROOT"
go build -o "$ARTIFACTS_DIR/bin/ky-admin" ./cmd/main.go

if [ $? -eq 0 ]; then
    echo "编译成功! 二进制文件位于: $ARTIFACTS_DIR/bin/ky-admin"
    chmod +x "$ARTIFACTS_DIR/bin/ky-admin"
else
    echo "编译失败!"
    exit 1
fi

echo "构建完成!" 