# migrate - 高级数据库迁移工具

`migrate`是一个功能完整的数据库迁移工具，支持多种操作模式和配置选项，适合生产环境部署和完整的数据库维护。

## 功能特点

- **多种操作模式**：支持迁移、数据填充、重置和备份等多种操作
- **配置灵活**：通过配置文件进行详细配置
- **事务支持**：所有操作支持事务，确保数据一致性
- **完整日志**：提供详细的操作日志，方便追踪和调试

## 使用方法

### 基本用法

```bash
# 基本迁移操作
go run cmd/migrate/main.go -config configs/config.yaml -action=migrate
```

### 可用操作

```bash
# 仅迁移表结构
go run cmd/migrate/main.go -config configs/config.yaml -action=migrate

# 迁移表结构并填充初始数据
go run cmd/migrate/main.go -config configs/config.yaml -action=seed

# 重置数据库（删除所有表并重新创建）
go run cmd/migrate/main.go -config configs/config.yaml -action=reset

# 备份数据库
go run cmd/migrate/main.go -config configs/config.yaml -action=backup
```

## 配置文件

工具需要一个有效的配置文件，默认路径为`configs/config.yaml`。配置文件应包含以下内容：

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  name: kysiondb
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

您可以根据不同环境创建多个配置文件：

- `config.development.yaml`：开发环境配置
- `config.production.yaml`：生产环境配置
- `config.test.yaml`：测试环境配置

## 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-config` | 配置文件路径 | `configs/config.yaml` |
| `-action` | 执行的操作类型 | `migrate` |

### 操作类型说明

| 操作类型 | 说明 |
|---------|------|
| `migrate` | 执行数据库表结构迁移，不影响已有数据 |
| `seed` | 执行表结构迁移并填充初始数据 |
| `reset` | 删除所有表并重新创建，谨慎使用 |
| `backup` | 备份当前数据库状态 |

## 工作原理

1. 工具首先读取配置文件，获取数据库连接信息和其他配置
2. 根据指定的操作类型执行相应的操作：
   - `migrate`：创建或更新表结构
   - `seed`：创建表结构并初始化基础数据
   - `reset`：删除所有表并重新创建
   - `backup`：创建数据库备份文件

## 代码结构

`main.go`文件包含了所有迁移逻辑，主要分为以下部分：

- 命令行参数解析
- 配置文件加载
- 数据库连接管理
- 迁移操作执行
- 数据初始化
- 错误处理和日志

## 注意事项

1. **数据备份**：在执行`reset`操作前，建议先使用`backup`操作备份数据
2. **生产环境**：在生产环境中使用时，请谨慎选择操作类型，避免意外删除数据
3. **权限要求**：执行该工具需要数据库用户具有创建、修改和删除表的权限
4. **配置文件**：确保配置文件格式正确，包含所有必要的配置项

## 故障排除

### 1. 配置文件未找到

症状：出现"配置文件未找到"的错误信息。

解决方案：

- 检查配置文件路径是否正确
- 确保配置文件存在并可读取

### 2. 数据库连接失败

症状：出现"连接数据库失败"的错误信息。

解决方案：

- 检查配置文件中的数据库连接信息
- 确认数据库服务器是否运行并可访问

### 3. 迁移操作失败

症状：出现"执行数据库迁移失败"的错误信息。

解决方案：

- 检查日志以确定具体错误原因
- 确保数据库用户具有足够的权限
- 验证表结构是否与模型定义一致
