# simplemigrate - 简化版数据库迁移工具

`simplemigrate`是一个简化版的数据库迁移工具，专为快速开发和测试环境设计。无需复杂的配置文件，通过环境变量即可方便地自定义数据库连接信息。

## 功能特点

- **简单易用**：无需复杂配置，一键执行迁移
- **自动化**：自动创建表结构和基础数据
- **可配置**：通过环境变量灵活配置数据库连接
- **增量式**：检测已有数据，避免重复迁移和初始化

## 使用方法

### 基本用法

```bash
# 使用默认配置运行
go run cmd/simplemigrate/main.go
```

### 自定义数据库连接

通过环境变量自定义数据库连接信息：

```bash
# 设置数据库连接信息并运行
DB_HOST=localhost \
DB_PORT=3306 \
DB_USER=root \
DB_PASSWORD=password \
DB_NAME=kysiondb \
go run cmd/simplemigrate/main.go
```

### 编译后运行

```bash
# 编译
go build -o simplemigrate cmd/simplemigrate/main.go

# 运行
./simplemigrate
```

## 环境变量配置

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| DB_HOST | 数据库主机地址 | 10.68.73.213 |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户名 | kysion_dev |
| DB_PASSWORD | 数据库密码 | TrY6eK8Bt28KwYii |
| DB_NAME | 数据库名称 | kysiondb_dev |

## 工作原理

1. 工具首先连接到指定的数据库
2. 检查并创建必要的表结构：
   - `sys_user`：用户表
   - `sys_role`：角色表
   - `sys_permission`：权限表
   - `sys_user_role`：用户-角色关联表
   - `sys_role_permission`：角色-权限关联表
3. 检查数据库中是否已有基础数据
4. 如果没有基础数据，则自动创建：
   - 三个基础角色：admin、manager、user
   - 一个超级管理员用户：admin (密码: admin123)

## 代码结构

```
cmd/simplemigrate/
├── main.go         # 入口文件
└── utils/
    └── migration.go # 迁移逻辑实现
```

## 开发与扩展

如需扩展或修改迁移工具的功能，可以编辑以下文件：

1. `utils/migration.go`：修改迁移逻辑、表结构或初始数据
2. `main.go`：修改程序入口

## 常见问题

### 1. 数据库连接失败

症状：出现"连接数据库失败"的错误信息。

解决方案：

- 检查环境变量中的数据库连接信息是否正确
- 确认数据库服务器是否运行并可访问
- 检查网络连接和防火墙设置

### 2. 表创建失败

症状：出现"自动迁移表结构失败"的错误信息。

解决方案：

- 确保数据库用户具有创建表的权限
- 检查数据库是否已存在且可访问

### 3. 数据初始化失败

症状：出现"创建角色失败"或"创建管理员用户失败"的错误信息。

解决方案：

- 检查表结构是否与模型定义一致
- 确认是否有唯一索引冲突（如角色代码重复）
