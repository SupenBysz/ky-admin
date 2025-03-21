# 命令行工具说明

本目录包含了项目相关的命令行工具，主要用于数据库迁移、初始化等操作。

## 工具列表

### 1. migrate

完整的数据库迁移工具，支持多种操作类型和配置文件。

特点：

- 支持迁移(migrate)、数据填充(seed)、重置(reset)、备份(backup)等操作
- 需要配置文件支持
- 适合完整的数据库部署和维护

用法：

```bash
go run cmd/migrate/main.go -config configs/config.yaml -action=migrate
```

参数说明：

- `-config`: 配置文件路径，默认为"configs/config.yaml"
- `-action`: 操作类型，可选值：
  - `migrate`: 仅迁移表结构
  - `seed`: 迁移表结构并初始化基础数据
  - `reset`: 重置数据库（删除所有表并重新创建）
  - `backup`: 备份数据库

### 2. simplemigrate

简化版数据库迁移工具，无需配置文件，通过环境变量即可配置。

特点：

- 使用简单，无需配置文件
- 自动执行表结构迁移和基础数据初始化
- 支持通过环境变量配置数据库连接信息
- 适合开发环境和快速部署

用法：

```bash
go run cmd/simplemigrate/main.go
```

环境变量配置：

- `DB_HOST`: 数据库主机地址，默认为"10.68.73.213"
- `DB_PORT`: 数据库端口，默认为"3306"
- `DB_USER`: 数据库用户名，默认为"kysion_dev"
- `DB_PASSWORD`: 数据库密码，默认为"TrY6eK8Bt28KwYii"
- `DB_NAME`: 数据库名称，默认为"kysiondb_dev"

示例：

```bash
# 使用默认配置
go run cmd/simplemigrate/main.go

# 指定数据库连接信息
DB_HOST=localhost DB_PORT=3306 DB_USER=root DB_PASSWORD=123456 DB_NAME=kysiondb_dev go run cmd/simplemigrate/main.go
```

## 迁移工具创建的表

迁移工具会创建以下表：

- `sys_user`: 用户表
- `sys_role`: 角色表
- `sys_permission`: 权限表
- `sys_user_role`: 用户-角色关联表
- `sys_role_permission`: 角色-权限关联表

## 初始化的基础数据

### 角色

工具会自动创建以下基础角色：

- `admin`: 系统管理员
- `manager`: 管理员
- `user`: 普通用户

### 用户

工具会自动创建一个超级管理员用户：

- 用户名: `admin`
- 密码: `admin123`
- 角色: `admin`

## 注意事项

1. 首次运行时会自动创建表结构和基础数据
2. 再次运行时会检查表结构和数据是否存在，存在则跳过初始化
3. simplemigrate工具默认会输出SQL执行日志，方便调试
4. 使用前请确保数据库用户具有创建表的权限

## 故障排除

### 连接数据库失败

检查环境变量或配置文件中的数据库连接信息是否正确，确保数据库服务已启动且网络连接正常。

### 创建表失败

检查数据库用户是否具有创建表的权限，以及数据库是否存在。

### 初始化数据失败

可能是由于表结构变更或数据已存在导致，检查错误信息以确定具体原因。
