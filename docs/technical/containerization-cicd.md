# 容器化与CI/CD文档

本文档详细描述了KY-Admin项目的容器化配置和CI/CD流程。

## 目录结构

```
docker/
├── .env.example              # 环境变量示例文件
├── Dockerfile                # 多阶段构建Dockerfile
├── README.md                 # Docker使用说明
├── deploy.sh                 # 部署和管理脚本
├── docker-compose.yml        # 生产环境Docker Compose配置
├── init-scripts/             # 初始化脚本目录
│   └── mysql/
│       └── 01-init.sql       # MySQL初始化SQL脚本
└── nginx/                    # Nginx配置目录
    └── conf.d/
        ├── default.conf      # HTTP配置
        └── default-ssl.conf.example # HTTPS配置示例
```

## 容器化配置

### Dockerfile

多阶段构建的Dockerfile，包括：

1. 构建阶段：
   - 基于`golang:1.20-alpine`镜像
   - 设置Go环境变量
   - 复制代码并编译应用

2. 运行阶段：
   - 基于`alpine:3.16`镜像
   - 配置时区和必要环境
   - 从构建阶段复制编译好的应用
   - 设置健康检查

### docker-compose.yml

生产环境Docker Compose配置，包括：

1. 服务：
   - `ky-admin`：后端应用服务
   - `mysql`：MySQL数据库服务
   - `redis`：Redis缓存服务
   - `nginx`：Nginx反向代理服务

2. 网络：
   - `ky-network`：内部网络

3. 数据卷：
   - 用于持久化数据和配置

### 环境变量

`.env.example`文件中包含的关键环境变量：

- 数据库配置（主机、端口、用户名、密码等）
- Redis配置
- 应用配置（端口、密钥等）
- Nginx配置

### 初始化脚本

MySQL初始化脚本`01-init.sql`包含：

- 创建用户和角色表
- 创建权限表
- 设置初始管理员和测试用户
- 分配初始权限

### Nginx配置

1. HTTP配置(`default.conf`)：
   - API请求代理
   - 静态文件服务
   - 健康检查端点

2. HTTPS配置(`default-ssl.conf.example`)：
   - SSL/TLS配置示例
   - HTTPS重定向设置

## CI/CD流程

### GitHub Actions配置

`.github/workflows/ci-cd.yml`文件定义了完整的CI/CD流程：

1. 测试阶段：
   - 启动MySQL和Redis服务容器
   - 运行单元测试和集成测试
   - 生成测试报告

2. 构建阶段：
   - 构建Docker镜像
   - 推送镜像到Docker仓库
   - 使用分支名和版本号标记镜像

3. 部署阶段：
   - 开发环境部署（`dev`分支）
   - 生产环境部署（`main`分支）
   - 通过SSH执行远程部署命令

### 部署脚本

`deploy.sh`脚本提供以下功能：

- 部署服务：`./deploy.sh up`
- 停止服务：`./deploy.sh down`
- 重启服务：`./deploy.sh restart`
- 查看服务状态：`./deploy.sh status`
- 查看日志：`./deploy.sh logs [服务名]`
- 清理数据：`./deploy.sh clean`
- 备份数据库：`./deploy.sh backup`
- 恢复数据库：`./deploy.sh restore [备份文件]`
- 显示帮助信息：`./deploy.sh help`

## 使用指南

### 本地开发环境部署

1. 复制环境变量示例文件：

   ```bash
   cp .env.example .env
   ```

2. 编辑环境变量文件，设置必要参数：

   ```bash
   vim .env
   ```

3. 使用部署脚本启动服务：

   ```bash
   ./deploy.sh up
   ```

### CI/CD流程触发

1. 开发环境部署：
   - 推送或合并代码到`dev`分支
   - GitHub Actions自动构建和部署到开发服务器

2. 生产环境部署：
   - 推送或合并代码到`main`分支
   - GitHub Actions自动构建和部署到生产服务器

## 最佳实践

1. 环境变量管理：
   - 敏感信息（密码、密钥等）使用环境变量
   - 使用`.env`文件在本地开发，GitHub Secrets在CI/CD环境

2. 镜像标记策略：
   - 使用语义化版本号
   - 使用分支名作为开发版本标记

3. 数据备份：
   - 定期使用`./deploy.sh backup`备份数据
   - 重要变更前备份数据

4. 监控和日志：
   - 使用`./deploy.sh logs`查看服务日志
   - 设置健康检查监控服务状态
