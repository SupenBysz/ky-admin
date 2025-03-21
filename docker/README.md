# Docker 配置文件

本目录包含项目所需的Docker相关配置文件和数据。

## 目录结构

```
docker/
├── data/                     # 数据目录
│   └── mongo/                # MongoDB数据
├── docker-compose.test.yml   # 测试环境配置
├── docker-compose.yapi.yml   # YAPI环境配置
└── README.md                 # 说明文档
```

## 使用方法

### 启动测试环境

```bash
cd kysion.com/backend/docker
docker-compose -f docker-compose.test.yml up -d
```

测试环境包含：

- MySQL 数据库（端口：3306）
- Redis 缓存（端口：6379）

### 启动YAPI环境

```bash
cd kysion.com/backend/docker
docker-compose -f docker-compose.yapi.yml up -d
```

YAPI环境包含：

- MongoDB 数据库
- YAPI 接口管理平台（端口：3000）

## 访问YAPI

启动YAPI后，可通过以下方式访问：

- URL: <http://localhost:3000>
- 账号: <admin@kysion.com>
- 密码: ky-admin123
