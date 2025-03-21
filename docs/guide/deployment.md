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

# 部署指南

本文档提供了KY-Admin系统的详细部署步骤，包括开发环境、测试环境和生产环境的部署方法。

## 部署前提

在部署KY-Admin系统前，请确保满足以下条件：

- Docker 20.10+ 和 Docker Compose 2.0+
- Git 客户端
- 4GB+ RAM, 2+ CPU核心
- 20GB+ 磁盘空间
- 可访问互联网（用于拉取镜像）

## 环境准备

### 1. 克隆代码仓库

```bash
# 克隆项目代码
git clone https://github.com/yourusername/ky-admin.git
cd ky-admin
```

### 2. 配置环境变量

```bash
# 进入Docker配置目录
cd docker

# 复制环境变量模板
cp .env.example .env

# 编辑环境变量
vi .env
```

需要配置的主要环境变量：

- 数据库连接信息（主机、端口、用户名、密码等）
- Redis连接信息
- 应用配置（端口、密钥等）
- 域名配置

## 开发环境部署

### 使用本地Go环境

```bash
# 进入项目根目录
cd kysion.com/backend/ky-admin

# 安装依赖
go mod download

# 创建并编辑配置
cp configs/config.example.yaml configs/config.yaml
vi configs/config.yaml

# 启动服务
go run main.go
```

### 使用Docker部署开发环境

```bash
# 进入Docker目录
cd kysion.com/backend/docker

# 使用部署脚本启动开发环境
./deploy.sh up

# 查看服务状态
./deploy.sh status

# 查看日志
./deploy.sh logs
```

## 测试环境部署

测试环境建议使用Docker Compose进行部署，并启用YAPI接口测试平台。

```bash
# 进入Docker目录
cd kysion.com/backend/docker

# 启动测试环境
docker-compose -f docker-compose.test.yml up -d

# 启动YAPI平台
docker-compose -f docker-compose.yapi.yml up -d

# 查看服务状态
docker-compose -f docker-compose.test.yml ps
docker-compose -f docker-compose.yapi.yml ps
```

访问测试环境:

- 后端API：`http://localhost:8080/api/v1`
- Swagger文档：`http://localhost:8080/swagger/index.html`
- YAPI平台：`http://localhost:3000`

## 生产环境部署

### 手动部署方式

1. 服务器环境准备：

```bash
# 安装Docker和Docker Compose
curl -fsSL https://get.docker.com | sh
apt-get install -y docker-compose
```

2. 代码拉取与配置：

```bash
# 克隆代码到服务器
git clone https://github.com/yourusername/ky-admin.git
cd ky-admin/docker

# 配置环境变量
cp .env.example .env
vi .env
```

3. 启动服务：

```bash
# 添加执行权限
chmod +x deploy.sh

# 启动所有服务
./deploy.sh up
```

4. 配置Nginx：

```bash
# 复制SSL配置模板
cp nginx/conf.d/default-ssl.conf.example nginx/conf.d/default-ssl.conf

# 编辑配置文件，添加SSL证书路径
vi nginx/conf.d/default-ssl.conf

# 重启Nginx服务
./deploy.sh restart nginx
```

### 使用CI/CD自动部署

1. 配置GitHub Actions:
   - 在GitHub仓库中设置必要的Secrets:
     - `DOCKER_USERNAME`: DockerHub用户名
     - `DOCKER_PASSWORD`: DockerHub密码
     - `DEV_SSH_HOST`: 开发服务器IP
     - `DEV_SSH_USERNAME`: 开发服务器用户名
     - `DEV_SSH_KEY`: 开发服务器SSH私钥
     - `PROD_SSH_HOST`: 生产服务器IP
     - `PROD_SSH_USERNAME`: 生产服务器用户名
     - `PROD_SSH_KEY`: 生产服务器SSH私钥

2. 自动部署流程:
   - 提交代码到`dev`分支，自动部署到开发环境
   - 提交代码到`main`分支，自动部署到生产环境

3. 监控部署状态:
   - 在GitHub Actions页面查看部署进度和日志
   - 部署完成后，使用`./deploy.sh status`检查服务状态

## 部署脚本使用指南

`deploy.sh`脚本提供了多种功能来管理部署环境。主要命令包括：

```bash
# 启动服务
./deploy.sh up

# 停止服务
./deploy.sh down

# 重启服务
./deploy.sh restart [服务名]

# 查看服务状态
./deploy.sh status

# 查看日志
./deploy.sh logs [服务名]

# 备份数据库
./deploy.sh backup

# 恢复数据库
./deploy.sh restore [备份文件]

# 清理数据
./deploy.sh clean

# 显示帮助信息
./deploy.sh help
```

## 多环境部署配置

KY-Admin支持多种环境配置，通过不同的配置文件进行管理：

### 开发环境

使用`docker-compose.dev.yml`配置文件，特点是：

- 开启调试模式
- 使用较小资源限制
- 启用热重载

### 测试环境

使用`docker-compose.test.yml`配置文件，特点是：

- 配置测试数据库
- 启用测试覆盖率收集
- 启用测试报告生成

### 生产环境

使用`docker-compose.yml`配置文件，特点是：

- 优化性能配置
- 启用资源限制
- 配置健康检查
- 开启负载均衡

## 容器健康监控

KY-Admin系统内置了健康检查端点和监控配置：

- 健康检查API：`/health`, `/livez`, `/readyz`
- 容器健康检查：每个容器都配置了健康检查命令
- 监控系统：支持与Prometheus和Grafana集成

## 备份策略

建议采用以下备份策略：

1. **数据库备份**:
   - 使用`./deploy.sh backup`每天进行一次全量备份
   - 保留最近7天的每日备份和每月一次的月度备份

2. **配置备份**:
   - 将环境变量和配置文件纳入版本控制
   - 关键配置变更时进行备份

3. **日志备份**:
   - 配置日志轮转
   - 重要日志存档到外部存储

## 故障排除

### 常见问题解决方案

1. **容器启动失败**:
   - 检查日志：`./deploy.sh logs [服务名]`
   - 检查环境变量配置
   - 检查磁盘空间和权限

2. **数据库连接问题**:
   - 检查数据库容器状态：`./deploy.sh status mysql`
   - 验证数据库凭据
   - 检查网络连接

3. **API访问错误**:
   - 检查Nginx配置
   - 验证防火墙设置
   - 检查SSL证书

4. **性能问题**:
   - 检查资源使用情况：`docker stats`
   - 优化数据库查询
   - 检查日志中的慢查询

## 升级指南

### 小版本升级

```bash
# 拉取最新代码
git pull

# 重新构建并启动服务
./deploy.sh restart
```

### 大版本升级

```bash
# 备份数据
./deploy.sh backup

# 拉取新版本代码
git fetch
git checkout v2.0.0

# 应用数据库迁移
./deploy.sh migrate

# 重新构建并启动服务
./deploy.sh up
```

## 安全建议

1. **网络安全**:
   - 使用防火墙限制端口访问
   - 为生产环境配置HTTPS
   - 使用固定IP白名单

2. **数据安全**:
   - 定期备份数据
   - 加密敏感数据
   - 实施最小权限原则

3. **系统安全**:
   - 定期更新Docker镜像
   - 扫描容器安全漏洞
   - 监控异常访问
