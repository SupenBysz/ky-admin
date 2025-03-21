# 代码生成器设计与使用指南

## 1. 概述

代码生成器是KY-Admin框架的核心功能之一，旨在通过自动生成代码架构和基础功能，大幅提高开发效率，减少重复劳动，并确保代码质量和规范统一。本文档详细介绍了代码生成器的设计理念、实现方式和使用方法。

## 2. 设计理念

代码生成器遵循以下设计理念：

- **模板驱动**：基于Go模板引擎，实现高度可定制的代码生成
- **约定优于配置**：遵循固定的项目结构和命名规范，减少配置复杂度
- **完整性**：生成完整的CRUD操作及相关文件，包括模型、控制器、路由、服务、仓库等
- **可扩展性**：支持自定义模板，满足不同业务场景需求
- **代码质量**：生成符合最佳实践的高质量代码，包含注释和文档

## 3. 功能特性

代码生成器具备以下核心功能：

1. **数据库表结构解析**：自动读取数据库表结构，分析字段类型、约束、关系等
2. **多层架构代码生成**：生成符合项目架构的各层代码（模型、存储库、服务、控制器等）
3. **自动生成API接口**：包括CRUD基础操作及定制化接口
4. **自动生成YAPI文档**：为API接口生成标准YAPI文档
5. **自动生成单元测试**：为关键功能生成测试代码
6. **支持多种数据类型**：包括基础类型、数组、对象等
7. **关联关系处理**：支持一对一、一对多、多对多等数据关系
8. **定制化生成**：支持通过配置文件调整生成内容

## 4. 实现方式

### 4.1 技术栈

- Go语言标准库的text/template包
- Golang的reflect包进行类型反射
- GORM进行数据库操作和表结构解析
- Cobra命令行工具构建CLI界面

### 4.2 核心组件

1. **模板引擎**：基于Go标准库的template包，支持条件判断、循环、函数调用等
2. **数据库连接器**：连接各类数据库并读取表结构信息
3. **模型解析器**：解析数据库表结构或用户输入，构建模型信息
4. **代码生成器**：根据模板和模型信息生成代码文件
5. **命令行工具**：提供友好的命令行界面，接收用户参数

### 4.3 目录结构

```
pkg/generator/
├── cmd/              # 命令行工具实现
├── config/           # 配置相关代码
├── database/         # 数据库连接与解析
├── generator/        # 核心生成逻辑
├── model/            # 模型定义和处理
├── template/         # 模板文件目录
│   ├── controller/   # 控制器模板
│   ├── model/        # 模型模板
│   ├── repository/   # 存储库模板
│   ├── router/       # 路由模板
│   ├── service/      # 服务模板
│   └── test/         # 测试代码模板
└── util/             # 工具函数
```

## 5. 使用方法

### 5.1 命令行参数

代码生成器提供了丰富的命令行选项：

```bash
# 基础用法
go run cmd/generator/main.go generate --table user_info

# 完整参数示例
go run cmd/generator/main.go generate \
  --db-type mysql \
  --db-dsn "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local" \
  --table user_info \
  --module github.com/SupenBysz/ky-admin \
  --output ./internal \
  --template-dir ./pkg/generator/template \
  --force
```

参数说明：

- `--db-type`：数据库类型（支持mysql、postgres、sqlite）
- `--db-dsn`：数据库连接字符串
- `--table`：要生成代码的表名，支持多个表（逗号分隔）
- `--module`：Go模块名称
- `--output`：代码输出目录
- `--template-dir`：自定义模板目录（可选）
- `--force`：强制覆盖已存在的文件
- `--exclude`：排除生成的代码类型（model,service,controller等，逗号分隔）
- `--prefix`：表名前缀，生成代码时会移除此前缀

### 5.2 配置文件

除了命令行参数，还支持通过配置文件进行更详细的定制：

```yaml
# generator.yaml
database:
  type: mysql
  dsn: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

generator:
  module: github.com/SupenBysz/ky-admin
  output: ./internal
  templateDir: ./pkg/generator/template
  force: false
  
tables:
  - name: user_info
    modelName: User  # 自定义模型名称
    excludes: []     # 排除生成的部分
    fields:          # 自定义字段映射
      id:
        jsonTag: id
        validTag: "required"
      created_at:
        ignore: true  # 忽略某些字段
  
  - name: role_info
    modelName: Role
    relations:       # 定义关联关系
      - type: many2many
        table: user_role
        reference: User
```

### 5.3 模板自定义

用户可以自定义模板以满足特定需求：

1. 复制默认模板目录到自定义位置
2. 修改模板文件
3. 使用 `--template-dir` 参数指定自定义模板目录

模板文件使用Go模板语法，例如：

```go
// 模型模板示例
package model

import (
 "time"
 
 "gorm.io/gorm"
)

// {{.ModelName}} 表示{{.Comment}}
type {{.ModelName}} struct {
 {{range .Fields}}
 {{.Name}} {{.Type}} {{.Tags}} {{if .Comment}}// {{.Comment}}{{end}}
 {{end}}
}

// TableName 设置表名
func ({{.ModelName}}) TableName() string {
 return "{{.TableName}}"
}
```

## 6. 代码生成示例

### 6.1 输入：数据库表结构

假设有以下数据库表：

```sql
CREATE TABLE `user_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '电话',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';
```

### 6.2 输出：生成的代码文件

运行生成器后，将创建以下文件：

```
internal/
├── model/
│   └── user.go                 # 数据模型
├── repository/
│   └── user_repository.go      # 数据访问层
├── service/
│   └── user_service.go         # 业务逻辑层
├── controller/
│   └── user_controller.go      # 控制器层
├── router/
│   └── user_router.go          # 路由注册
└── dto/
    └── user_dto.go             # 数据传输对象
```

### 6.3 模型代码示例

生成的模型代码示例：

```go
// model/user.go
package model

import (
 "time"
 
 "gorm.io/gorm"
)

// User 表示用户信息
type User struct {
 ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:用户ID"`
 Username  string         `json:"username" gorm:"type:varchar(50);not null;uniqueIndex:idx_username;comment:用户名"`
 Password  string         `json:"-" gorm:"type:varchar(100);not null;comment:密码"`
 RealName  string         `json:"realName" gorm:"type:varchar(50);comment:真实姓名"`
 Email     string         `json:"email" gorm:"type:varchar(100);comment:邮箱"`
 Phone     string         `json:"phone" gorm:"type:varchar(20);comment:电话"`
 Avatar    string         `json:"avatar" gorm:"type:varchar(255);comment:头像"`
 Status    int8           `json:"status" gorm:"type:tinyint(4);not null;default:1;comment:状态：0-禁用，1-启用"`
 CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
 UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime;comment:更新时间"`
 DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (User) TableName() string {
 return "user_info"
}
```

### 6.4 控制器代码示例

生成的控制器代码示例：

```go
// controller/user_controller.go
package controller

import (
 "github.com/gin-gonic/gin"
 "github.com/SupenBysz/ky-admin/internal/dto"
 "github.com/SupenBysz/ky-admin/internal/service"
 "github.com/SupenBysz/ky-admin/pkg/response"
)

// UserController 用户控制器
type UserController struct {
 userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
 return &UserController{
  userService: userService,
 }
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.CreateUserRequest true "用户信息"
// @Success 200 {object} response.Response{data=dto.UserResponse} "成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /users [post]
func (c *UserController) Create(ctx *gin.Context) {
 var req dto.CreateUserRequest
 if err := ctx.ShouldBindJSON(&req); err != nil {
  response.Failed(ctx, err)
  return
 }
 
 user, err := c.userService.Create(ctx, req)
 if err != nil {
  response.Failed(ctx, err)
  return
 }
 
 response.Success(ctx, user)
}

// 其他方法（Get, List, Update, Delete等）
// ...
```

## 7. 注意事项与最佳实践

1. **命名规范**：
   - 表名使用下划线命名法（snake_case）
   - 生成的Go代码使用驼峰命名法（modelName使用大驼峰，字段使用小驼峰）

2. **表结构设计**：
   - 每个表建议包含id, created_at, updated_at, deleted_at字段
   - 为字段添加有意义的注释，会自动生成到代码文档中
   - 合理设置字段类型和约束

3. **性能考虑**：
   - 生成大量代码可能需要一些时间，特别是处理复杂表关系时
   - 对于大型项目，建议按模块分批生成

4. **代码审查**：
   - 生成的代码应经过审查后再投入使用
   - 可能需要对生成的代码进行一些微调

5. **版本控制**：
   - 建议将自定义模板纳入版本控制
   - 记录使用的生成器版本，便于后续维护

## 8. 常见问题

1. **问题**：生成的代码与数据库表结构不一致  
   **解决**：确保数据库连接正常，且有足够的权限查询表结构

2. **问题**：字段类型映射不正确  
   **解决**：检查数据库字段类型，必要时通过配置文件手动映射

3. **问题**：关联关系处理不正确  
   **解决**：在配置文件中明确指定表间关系

4. **问题**：生成的API不符合REST规范  
   **解决**：自定义API模板或调整路由配置

5. **问题**：生成的代码编译失败  
   **解决**：检查Go模块设置，确保依赖项正确配置

## 9. 扩展与进阶用法

1. **批量生成**：

   ```bash
   # 生成多个表的代码
   go run cmd/generator/main.go generate --table user_info,role_info,permission_info
   ```

2. **关联表生成**：

   ```bash
   # 使用配置文件指定关联关系
   go run cmd/generator/main.go generate --config generator.yaml
   ```

3. **生成前端代码**：

   ```bash
   # 同时生成后端和前端代码
   go run cmd/generator/main.go generate --table user_info --with-frontend
   ```

4. **插件扩展**：
   编写自定义插件扩展功能，如生成特定业务逻辑或集成第三方服务

## 10. 参考资源

- [Go Template文档](https://golang.org/pkg/text/template/)
- [GORM官方文档](https://gorm.io/docs/)
- [Cobra命令行库](https://github.com/spf13/cobra)
- [代码生成器样板项目](https://github.com/SupenBysz/ky-generator)
- [项目结构最佳实践](https://github.com/golang-standards/project-layout)
- [YAPI官方文档](https://hellosean1025.github.io/yapi/)
