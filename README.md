# KY-Admin 后端项目

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/yourusername/ky-admin)
[![Test Coverage](https://img.shields.io/badge/coverage-85%25-green)](https://github.com/yourusername/ky-admin)

KY-Admin 是一个基于Go+Gin+GORM的企业级后台管理系统后端，提供高效、安全、可扩展的管理平台API支持。

## 📋 开发进度

- ✅ 阶段1（基础框架搭建）已于2025-03-21完成
- 🔄 阶段2（核心功能开发）进度65%，预计2025-04-15完成
- 🔄 阶段3（业务模块开发）进度10%，已开始设计与规划
- 🔄 阶段4（性能优化与系统完善）进度25%，已完成项目结构和CI/CD优化

**最新更新（2025-03-22）：**

- 优化了项目结构，将构建产物和日志移出源码目录
- 改进了CI/CD工作流，统一了流程并修复了版本问题
- 用户认证与权限系统基本完成，RBAC模型已实现

## 🚀 快速开始

### 开发环境启动

```bash
# 克隆项目
git clone https://github.com/yourusername/ky-admin.git

# 进入项目目录
cd ky-admin

# 安装依赖
go mod download

# 创建并编辑配置
cp configs/config.example.yaml configs/config.yaml

# 启动服务
go run main.go
```

### 使用Docker启动

```bash
# 使用部署脚本
cd docker
chmod +x deploy.sh
./deploy.sh deploy
```

## 🏗️ 项目结构

```
.
├── api             # API层
├── configs         # 配置文件
├── docs            # 文档
├── docker          # Docker相关配置
├── internal        # 内部应用代码
│   ├── app         # 应用实例
│   ├── config      # 配置管理
│   ├── controller  # 控制器
│   ├── middleware  # 中间件
│   ├── model       # 数据模型
│   ├── repository  # 数据操作层
│   ├── service     # 业务逻辑层
│   └── utils       # 工具函数
├── pkg             # 可重用的包
├── scripts         # 脚本文件
└── tests           # 测试代码
```

## 🔧 技术栈

- **Web框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM框架**: [GORM](https://gorm.io/)
- **数据库**: MySQL 8.0, Redis 7.0
- **认证授权**: JWT, RBAC
- **文档**: Swagger, YAPI
- **测试框架**: Go标准库测试框架, Testify
- **日志**: Zap
- **容器化**: Docker, Docker Compose
- **CI/CD**: GitHub Actions

## 📚 API文档

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **YAPI**: `http://localhost:3000`

## 💻 开发工具

- **编辑器**: VSCode, GoLand
- **API测试**: Postman, Swagger UI
- **数据库工具**: MySQL Workbench, Redis Desktop Manager

## 🛠️ 容器化与部署

项目提供了完整的Docker容器化解决方案：

- **多环境支持**: 开发、测试、生产环境配置
- **一键部署**: 使用 `./deploy.sh` 脚本管理容器
- **CI/CD集成**: 通过GitHub Actions自动测试、构建和部署
- **数据持久化**: 自动管理数据卷和备份
- **健康检查**: 容器健康状态监控

详情请参考 [容器化与CI/CD文档](docs/containerization-cicd.md)

## 🧪 测试

```bash
# 运行所有测试
cd tests
make test-all

# 生成测试覆盖率报告
make test-coverage
```

## 📄 协议

本项目采用MIT协议，详情请见[LICENSE](LICENSE)文件。
