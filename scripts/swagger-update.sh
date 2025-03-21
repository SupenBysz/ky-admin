#!/bin/bash

# Swagger文档更新脚本
# 用于自动生成Swagger文档并推送到YAPI

set -e  # 发生错误时退出

# 定义颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # 无颜色

echo -e "${GREEN}开始更新Swagger文档...${NC}"

# 检查swag命令是否存在
if ! command -v swag &> /dev/null; then
    echo -e "${RED}swag命令不存在，正在安装...${NC}"
    go install github.com/swaggo/swag/cmd/swag@latest
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# 脚本目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# 项目根目录
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
# Node.js依赖目录
NODE_DIR="$SCRIPT_DIR/node"

# 创建Node.js依赖目录（如果不存在）
mkdir -p "$NODE_DIR"

# 检查axios是否已安装
if [ ! -d "$NODE_DIR/node_modules/axios" ]; then
    echo -e "${YELLOW}axios未安装，正在安装...${NC}"
    # 切换到Node.js依赖目录并安装axios
    cd "$NODE_DIR"
    if [ ! -f "package.json" ]; then
        echo '{"name":"swagger-sync","private":true}' > package.json
    fi
    npm install axios --no-package-lock
    cd "$PROJECT_DIR"
fi

# Swagger输出目录
SWAGGER_OUTPUT_DIR="$PROJECT_DIR/ky-admin/swagger"

echo -e "${GREEN}使用swag生成Swagger文档...${NC}"
cd "$PROJECT_DIR/ky-admin"
swag init -g cmd/main.go -o "$SWAGGER_OUTPUT_DIR"

# 判断是否需要同步到YAPI
if [ -n "$SYNC_TO_YAPI" ] && [ "$SYNC_TO_YAPI" = "true" ]; then
    echo -e "${GREEN}同步Swagger文档到YAPI...${NC}"
    # 使用NODE_PATH环境变量使脚本能找到node_modules
    NODE_PATH="$NODE_DIR/node_modules" node "$SCRIPT_DIR/sync-swagger-to-yapi.js"
else
    echo -e "${YELLOW}跳过同步到YAPI (设置SYNC_TO_YAPI=true以启用)${NC}"
fi

echo -e "${GREEN}Swagger文档更新完成!${NC}"
echo -e "文档位置: $SWAGGER_OUTPUT_DIR"
echo -e "可通过 http://localhost:8080/swagger/index.html 访问Swagger UI" 