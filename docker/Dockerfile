# 第一阶段：构建阶段
FROM golang:1.20-alpine AS builder

# 设置构建环境
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY ./ky-admin/go.mod ./ky-admin/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY ./ky-admin ./

# 构建应用
RUN go build -ldflags="-s -w" -o /app/bin/ky-admin ./cmd/main.go

# 第二阶段：运行阶段
FROM alpine:3.16

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从构建阶段复制构建好的应用
COPY --from=builder /app/bin/ky-admin /app/

# 复制配置文件
COPY ./ky-admin/configs /app/configs/

# 创建数据目录
RUN mkdir -p /app/data /app/logs

# 设置环境变量
ENV APP_ENV=production \
    APP_PORT=8080

# 暴露端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O- http://localhost:8080/health || exit 1

# 设置容器启动命令
CMD ["/app/ky-admin"] 