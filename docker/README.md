# Docker 配置文件

本目录包含项目所需的Docker相关配置文件和数据。

## 目录结构

```
docker/
├── data/                     # 数据目录
│   └── mongo/                # MongoDB数据
├── init-scripts/             # 初始化脚本
│   └── mysql/                # MySQL初始化脚本
├── nginx/                    # Nginx配置
│   ├── conf.d/               # 配置文件
│   └── ssl/                  # SSL证书
├── backups/                  # 数据库备份
├── .env.example              # 环境变量示例
├── deploy.sh                 # 部署脚本
├── Dockerfile                # 应用Dockerfile
├── docker-compose.yml        # 生产环境配置
├── docker-compose.test.yml   # 测试环境配置
├── docker-compose.yapi.yml   # YAPI环境配置
└── README.md                 # 说明文档
```

## 快速开始

### 1. 复制环境变量示例文件

```bash
cp .env.example .env
```

### 2. 按需修改环境变量

编辑`.env`文件，根据实际情况修改数据库密码和其他配置。

### 3. 使用部署脚本

部署脚本提供了多种命令来简化Docker环境的管理。

```bash
# 显示帮助信息
./deploy.sh help

# 部署服务
./deploy.sh deploy

# 查看服务状态
./deploy.sh status

# 查看服务日志
./deploy.sh logs ky-admin

# 停止服务
./deploy.sh stop

# 重启服务
./deploy.sh restart

# 备份数据库
./deploy.sh backup_db

# 恢复数据库
./deploy.sh restore_db backups/mysql_backup_20250322_123000.sql
```

## 环境说明

### 生产环境

使用`docker-compose.yml`配置文件，包含以下服务：

- **ky-admin**: 后端服务
- **mysql**: 数据库服务
- **redis**: 缓存服务
- **nginx**: Web服务器

### 测试环境

使用`docker-compose.test.yml`配置文件，包含：

- **mysql**: 测试数据库
- **redis**: 测试缓存服务

### YAPI环境

使用`docker-compose.yapi.yml`配置文件，用于启动YAPI接口管理平台：

- **mongodb**: YAPI数据库
- **yapi**: YAPI接口管理平台

## 访问说明

### 主应用

- HTTP: `http://localhost`
- HTTPS: `https://localhost` (需配置SSL证书)
- API: `http://localhost/api/`
- Swagger: `http://localhost/swagger/`

### YAPI

- URL: `http://localhost:3000`
- 账号: `admin@kysion.com`
- 密码: `ky-admin123`

## 自定义配置

### SSL配置

1. 将SSL证书文件放入`nginx/ssl/`目录
2. 复制并修改SSL配置示例：

   ```bash
   cp nginx/conf.d/default-ssl.conf.example nginx/conf.d/default-ssl.conf
   ```

3. 编辑`default-ssl.conf`文件，修改证书路径和域名
