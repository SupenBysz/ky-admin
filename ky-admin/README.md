# KY-Admin 后端服务

KY-Admin是一个基于Go语言开发的管理系统后端服务，提供高性能、可扩展的API接口和管理功能。

## 项目架构

本项目采用现代化的设计理念和架构模式，主要包括：

- RESTful API设计
- 分层架构：表示层、业务逻辑层、数据访问层
- 依赖注入设计模式
- 中间件机制
- 统一的错误处理和日志系统

## 目录结构

```
ky-admin/
├── cmd/               # 主要的应用程序入口点
├── configs/           # 配置文件
├── docs/              # 文档
├── internal/          # 内部包，不对外暴露
│   ├── api/           # API处理器
│   ├── middleware/    # HTTP中间件
│   ├── model/         # 数据模型
│   ├── repository/    # 数据存储层
│   └── service/       # 业务逻辑层
├── pkg/               # 可重用的外部包
├── scripts/           # 脚本文件
└── test/              # 测试代码
```

## 技术栈

- 语言：Go 1.24+
- Web框架：Gin
- ORM：GORM
- 配置管理：Viper
- 日志：Zap
- 数据库：MariaDB 11.6.2+
- API文档：YAPI

## 开始使用

### 环境要求

- Go 1.24+
- MariaDB 11.6.2+

### 安装与运行

1. 克隆仓库

```bash
git clone https://github.com/SupenBysz/ky-admin.git
cd ky-admin
```

2. 安装依赖

```bash
go mod download
```

3. 配置数据库

```bash
# 编辑configs/config.yaml文件，配置数据库连接信息
```

4. 构建与运行

```bash
make build
make run
```

## 开发指南

请参考[开发文档](docs/development.md)获取更多细节。

## 贡献代码

欢迎提交Pull Request或Issue。在提交代码前，请确保：

1. 添加必要的测试用例
2. 遵循项目的代码规范
3. 更新相关文档

## 许可证

Copyright © 2025 Kysion.com
