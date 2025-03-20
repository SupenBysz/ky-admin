# 开发指南

本文档提供KY-Admin项目的开发指南，帮助开发者快速上手和理解项目结构。

## 目录

- [项目架构](#项目架构)
- [开发环境](#开发环境)
- [项目结构](#项目结构)
- [核心模块](#核心模块)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [常见问题](#常见问题)

## 项目架构

KY-Admin采用了清晰的分层架构设计：

```
┌─────────────┐
│  API层(Gin) │
└──────┬──────┘
       │
┌──────▼──────┐
│   服务层    │
└──────┬──────┘
       │
┌──────▼──────┐
│   仓储层    │
└──────┬──────┘
       │
┌──────▼──────┐
│  数据库层   │
└─────────────┘
```

- **API层**：处理HTTP请求，参数验证，返回结果
- **服务层**：实现业务逻辑
- **仓储层**：处理数据存储和查询
- **数据库层**：与数据库交互

## 开发环境

### 环境要求

- Go 1.20+
- MariaDB 11.6.2+
- Git

### 开发工具推荐

- IDE: GoLand, VSCode
- API测试: Postman
- 数据库工具: DBeaver, MySQL Workbench
- Git客户端: SourceTree, GitKraken

### 环境设置

1. 克隆项目

```bash
git clone https://github.com/SupenBysz/ky-admin.git
cd ky-admin
```

2. 安装依赖

```bash
go mod download
```

3. 运行开发服务器

```bash
make run
```

## 项目结构

```
ky-admin/
├── cmd/                 # 应用程序入口点
├── configs/             # 配置文件
├── docs/                # 文档
├── internal/            # 内部代码，不对外暴露
│   ├── api/             # API处理器
│   │   ├── handler/     # 请求处理器
│   │   ├── middleware/  # HTTP中间件
│   │   ├── router/      # 路由定义
│   │   └── validator/   # 请求验证器
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   └── service/         # 业务逻辑层
├── pkg/                 # 可重用的外部包
│   ├── config/          # 配置工具
│   ├── logger/          # 日志工具
│   ├── database/        # 数据库工具
│   ├── utils/           # 通用工具
│   └── errors/          # 错误处理
├── scripts/             # 脚本文件
└── test/                # 测试代码
```

## 核心模块

### 配置管理

使用Viper库实现配置管理，支持：

- 多环境配置
- 配置热加载
- 环境变量覆盖

```go
// 示例: 读取配置
cfg := config.GetConfig()
dbHost := cfg.GetString("database.host")
```

### 日志系统

使用zap实现高性能日志：

- 结构化日志
- 多级别日志
- 文件轮转

```go
// 示例: 使用日志
logger := log.GetLogger()
logger.Info("操作成功", zap.String("user", "admin"))
```

### 数据库操作

使用GORM进行数据库操作：

- 模型定义
- 查询构建
- 关联处理

```go
// 示例: 数据库操作
user := &model.User{Username: "admin"}
result := db.Create(user)
```

### API处理

使用Gin框架处理HTTP请求：

- 路由定义
- 中间件
- 参数验证
- 统一响应

```go
// 示例: API处理器
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, err)
        return
    }
    
    user, err := h.userService.Create(req)
    if err != nil {
        response.Fail(c, err)
        return
    }
    
    response.Success(c, user)
}
```

## 开发流程

1. **功能开发**
   - 从dev分支创建功能分支
   - 实现新功能
   - 添加单元测试
   - 提交PR到dev分支

2. **Bug修复**
   - 从dev分支创建bugfix分支
   - 修复bug
   - 添加测试
   - 提交PR到dev分支

3. **发布流程**
   - 从dev分支创建release分支
   - 测试并修复issue
   - 生成版本文档
   - 合并到main和dev分支

## 代码规范

### 命名规范

- 包名：小写，简短，不使用下划线
- 函数/方法：驼峰命名，如`GetUserByID`
- 变量：驼峰命名，如`userID`
- 常量：全大写，使用下划线，如`MAX_CONNECTIONS`
- 接口：以er结尾，如`Reader`

### 注释规范

- 包注释：在package语句之前添加，描述包的功能
- 函数/方法注释：描述功能、参数、返回值和错误
- 代码块注释：解释复杂的逻辑
- 符合golint要求

### 错误处理

- 使用有意义的错误消息
- 对错误进行包装以提供上下文
- 避免使用panic，除非是不可恢复的错误

## 常见问题

### 数据库连接问题

**问题**：启动时报数据库连接错误

**解决方案**：

- 检查配置文件中的数据库连接信息
- 确保MariaDB服务已启动
- 检查网络连接和防火墙设置

### 热更新配置

**问题**：修改配置文件不生效

**解决方案**：

- 确保使用的是开发环境
- 检查配置文件路径是否正确
- 重启应用（某些配置需要重启）

### 测试失败

**问题**：单元测试失败

**解决方案**：

- 确保测试数据库配置正确
- 检查测试环境
- 确保所有依赖都已模拟或提供
