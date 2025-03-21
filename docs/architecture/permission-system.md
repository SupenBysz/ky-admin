# 权限系统设计文档

## 1. 概述

KY-Admin框架的权限系统是一个灵活、精细、高性能的RBAC（基于角色的访问控制）权限管理模块，支持多租户场景，可满足各类企业级应用的权限管理需求。本文档详细阐述了该权限系统的设计理念、架构模型、核心功能以及使用方法。

## 2. 设计理念

KY-Admin权限系统遵循以下设计理念：

1. **多层次权限控制**：支持用户-角色-权限-资源的多层次权限控制
2. **领域驱动设计**：基于DDD(领域驱动设计)思想构建权限模型
3. **高性能设计**：通过缓存优化、预加载等技术确保权限验证高性能
4. **可扩展性**：权限规则支持自定义扩展，满足复杂业务场景
5. **开发友好**：提供简洁API和声明式注解，降低开发者使用门槛
6. **多租户支持**：内置多租户隔离能力，满足SaaS应用需求

## 3. 核心模型

### 3.1 基础概念

- **用户（User）**：系统的操作者，可以是实际用户、第三方系统或API客户端
- **角色（Role）**：权限的集合，用户通过被分配到不同角色获得相应权限
- **权限（Permission）**：对资源执行特定操作的能力描述
- **资源（Resource）**：系统中可被访问的对象，如API、菜单、按钮、数据等
- **租户（Tenant）**：系统的组织单位，可以是公司、部门或团队
- **策略（Policy）**：定义特定条件下的权限规则

### 3.2 数据模型

#### 用户表（sys_user）

```sql
CREATE TABLE `sys_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '电话',
  `tenant_id` bigint(20) NOT NULL COMMENT '租户ID',
  `dept_id` bigint(20) DEFAULT NULL COMMENT '部门ID',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username_tenant` (`username`,`tenant_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_dept_id` (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';
```

#### 角色表（sys_role）

```sql
CREATE TABLE `sys_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `tenant_id` bigint(20) NOT NULL COMMENT '租户ID',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code_tenant` (`code`,`tenant_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
```

#### 用户角色关联表（sys_user_role）

```sql
CREATE TABLE `sys_user_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `tenant_id` bigint(20) NOT NULL COMMENT '租户ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';
```

#### 权限表（sys_permission）

```sql
CREATE TABLE `sys_permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `name` varchar(50) NOT NULL COMMENT '权限名称',
  `code` varchar(100) NOT NULL COMMENT '权限编码',
  `type` tinyint(4) NOT NULL COMMENT '类型：1-菜单，2-按钮，3-API，4-数据',
  `parent_id` bigint(20) DEFAULT '0' COMMENT '父权限ID',
  `path` varchar(255) DEFAULT NULL COMMENT '路径',
  `method` varchar(10) DEFAULT NULL COMMENT 'HTTP方法',
  `icon` varchar(100) DEFAULT NULL COMMENT '图标',
  `component` varchar(255) DEFAULT NULL COMMENT '组件',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';
```

#### 角色权限关联表（sys_role_permission）

```sql
CREATE TABLE `sys_role_permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `permission_id` bigint(20) NOT NULL COMMENT '权限ID',
  `tenant_id` bigint(20) NOT NULL COMMENT '租户ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`),
  KEY `idx_permission_id` (`permission_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';
```

#### 租户表（sys_tenant）

```sql
CREATE TABLE `sys_tenant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '租户ID',
  `name` varchar(50) NOT NULL COMMENT '租户名称',
  `code` varchar(50) NOT NULL COMMENT '租户编码',
  `contact_name` varchar(50) DEFAULT NULL COMMENT '联系人',
  `contact_phone` varchar(20) DEFAULT NULL COMMENT '联系电话',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `expire_time` timestamp NULL DEFAULT NULL COMMENT '过期时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';
```

### 3.3 对象模型

```go
// 用户模型
type User struct {
    ID        uint64    `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`
    RealName  string    `json:"realName"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    TenantID  uint64    `json:"tenantId"`
    DeptID    uint64    `json:"deptId"`
    Status    int8      `json:"status"`
    Roles     []Role    `json:"roles" gorm:"many2many:sys_user_role;"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// 角色模型
type Role struct {
    ID          uint64       `json:"id" gorm:"primaryKey"`
    Name        string       `json:"name"`
    Code        string       `json:"code"`
    TenantID    uint64       `json:"tenantId"`
    Description string       `json:"description"`
    Sort        int          `json:"sort"`
    Status      int8         `json:"status"`
    Permissions []Permission `json:"permissions" gorm:"many2many:sys_role_permission;"`
    CreatedAt   time.Time    `json:"createdAt"`
    UpdatedAt   time.Time    `json:"updatedAt"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// 权限模型
type Permission struct {
    ID        uint64    `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    Code      string    `json:"code"`
    Type      int8      `json:"type"`
    ParentID  uint64    `json:"parentId"`
    Path      string    `json:"path"`
    Method    string    `json:"method"`
    Icon      string    `json:"icon"`
    Component string    `json:"component"`
    Sort      int       `json:"sort"`
    Status    int8      `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// 租户模型
type Tenant struct {
    ID           uint64    `json:"id" gorm:"primaryKey"`
    Name         string    `json:"name"`
    Code         string    `json:"code"`
    ContactName  string    `json:"contactName"`
    ContactPhone string    `json:"contactPhone"`
    Status       int8      `json:"status"`
    ExpireTime   time.Time `json:"expireTime"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
```

## 4. 核心功能

### 4.1 用户认证

系统支持多种认证方式：

1. **基于JWT的认证**：
   - 用户登录后获取JWT令牌
   - 令牌包含用户ID、租户ID、角色等信息
   - 支持令牌刷新和自动续期

2. **集成OAuth2.0**：
   - 支持第三方账号登录
   - 支持OAuth2.0的授权码、隐式、客户端凭证等模式

3. **双因素认证**：
   - 可选支持双因素认证
   - 支持TOTP、短信等验证方式

### 4.2 访问控制

1. **基于角色的访问控制（RBAC）**：
   - 用户关联角色，角色关联权限
   - 支持角色继承

2. **基于属性的访问控制（ABAC）**：
   - 支持基于用户属性、资源属性、环境属性的控制
   - 支持复杂的条件表达式

3. **数据权限控制**：
   - 支持行级数据权限
   - 支持多种数据范围：全部、本部门、本部门及子部门、自定义部门、仅本人

### 4.3 权限验证

1. **API权限验证**：
   - 通过中间件自动验证API访问权限
   - 支持通配符路径和HTTP方法匹配

2. **界面元素权限**：
   - 菜单权限控制
   - 按钮权限控制

3. **数据权限过滤**：
   - 自动添加数据权限过滤条件
   - 支持跨表关联查询的权限控制

### 4.4 多租户支持

1. **租户隔离模式**：
   - 支持独立数据库、共享数据库独立schema、共享数据库共享表等模式
   - 默认采用共享数据库、共享表、行级租户隔离模式

2. **租户上下文**：
   - 租户信息贯穿整个请求生命周期
   - 支持租户识别、切换、验证

3. **超级管理员模式**：
   - 支持跨租户管理功能
   - 超级管理员不受租户隔离限制

## 5. 使用指南

### 5.1 权限控制注解

```go
// 1. 控制器方法权限注解
// @RequiresPermissions("user:create")
// @RequiresRoles("admin")
func (c *UserController) Create(ctx *gin.Context) {
    // ...
}

// 2. 多个权限验证，支持AND和OR逻辑
// @RequiresPermissions(value={"user:create", "user:edit"}, logical=Logical.OR)
func (c *UserController) Update(ctx *gin.Context) {
    // ...
}
```

### 5.2 中间件配置

```go
// 全局权限中间件
router.Use(middleware.Permission())

// 路由组权限中间件
userGroup := router.Group("/users")
userGroup.Use(middleware.RequiresPermissions("user:view"))
{
    userGroup.GET("", userController.List)
    userGroup.GET("/:id", userController.Get)
    // ...
}
```

### 5.3 数据权限配置

```go
// 1. 服务层数据权限配置
// @DataPermission(tableAlias="u")
func (s *UserServiceImpl) List(ctx context.Context, query dto.UserQuery) ([]model.User, error) {
    db := s.db.WithContext(ctx)
    db = permission.AddDataPermission(db, "u")
    // ...
}

// 2. 自定义数据权限规则
func (s *OrderServiceImpl) GetOrdersByCustomer(ctx context.Context, customerId uint64) ([]model.Order, error) {
    db := s.db.WithContext(ctx)
    db = permission.AddCustomDataPermission(db, func(db *gorm.DB) *gorm.DB {
        return db.Where("customer_id = ? OR created_by = ?", customerId, permission.GetCurrentUserId(ctx))
    })
    // ...
}
```

### 5.4 权限验证工具

```go
// 1. 检查当前用户是否有指定权限
if permission.HasPermission(ctx, "user:create") {
    // 有权限执行的代码
}

// 2. 检查当前用户是否具有指定角色
if permission.HasRole(ctx, "admin") {
    // 有角色执行的代码
}

// 3. 复杂权限表达式验证
if permission.CheckExpression(ctx, "(role:'admin' OR permission:'user:*') AND tenant:1") {
    // 满足表达式执行的代码
}
```

### 5.5 菜单和按钮权限

```jsx
// 前端权限指令示例（Vue）
<template>
  <div>
    <button v-permission="'user:create'">创建用户</button>
    <button v-permission="'user:edit'">编辑用户</button>
    <button v-role="'admin'">管理员操作</button>
  </div>
</template>
```

## 6. 缓存设计

为了提高权限验证性能，系统采用多级缓存设计：

1. **一级缓存**：请求上下文缓存，单次请求内有效
2. **二级缓存**：本地缓存，单个服务实例内有效
3. **三级缓存**：分布式缓存，整个系统内有效

缓存内容包括：

- 用户-角色映射
- 角色-权限映射
- 权限规则
- 租户信息

缓存更新策略：

- 权限变更时主动失效
- 设置合理的过期时间
- 采用版本号机制避免缓存一致性问题

## 7. 权限初始化

系统提供权限数据初始化功能：

1. **内置角色**：
   - 超级管理员：拥有所有权限，不受租户限制
   - 租户管理员：拥有租户内所有权限
   - 普通用户：基础功能权限

2. **内置权限**：
   - 系统管理权限
   - 用户管理权限
   - 角色管理权限
   - 权限管理权限
   - 租户管理权限

3. **初始化方式**：
   - 通过SQL脚本初始化
   - 通过代码自动检测并初始化
   - 支持自定义初始化数据

```go
// 初始化示例
func InitPermissions() {
    // 注册系统权限
    for _, api := range router.GetAPIList() {
        permission.RegisterAPIPermission(api.Path, api.Method, api.Name)
    }
    
    // 创建基础角色
    permission.CreateDefaultRoles()
    
    // 分配基础权限
    permission.AssignDefaultPermissions()
}
```

## 8. 权限审计

系统提供完善的权限操作审计功能：

1. **记录内容**：
   - 权限变更记录
   - 权限验证失败记录
   - 敏感操作记录

2. **审计方式**：
   - 数据库记录
   - 日志文件
   - 外部审计系统对接

3. **审计查询**：
   - 支持多维度查询
   - 支持导出审计记录
   - 支持审计报告生成

## 9. 最佳实践

1. **权限粒度控制**：
   - 根据实际需求确定权限粒度
   - 避免过细粒度带来的管理复杂性

2. **角色设计**：
   - 基于职责设计角色
   - 避免角色爆炸
   - 善用角色继承关系

3. **权限命名规范**：
   - 采用`资源:操作`的命名方式
   - 例如：`user:create`, `role:view`
   - 支持通配符：`user:*`, `*:view`

4. **租户隔离**：
   - 确保数据不会跨租户泄露
   - 公共数据考虑独立标记

5. **性能优化**：
   - 合理使用缓存
   - 批量预加载权限数据
   - 避免频繁权限检查

## 10. 常见问题

1. **问题**：权限验证影响系统性能  
   **解决**：使用多级缓存，减少数据库查询；批量预加载权限数据

2. **问题**：权限规则过于复杂难以管理  
   **解决**：遵循最小权限原则；使用权限组简化管理；提供可视化配置工具

3. **问题**：跨租户数据访问需求  
   **解决**：设计特殊权限标识；超级管理员模式；数据共享机制

4. **问题**：如何处理动态权限需求  
   **解决**：支持运行时权限规则注册；提供规则引擎；自定义权限处理器

5. **问题**：权限数据一致性问题  
   **解决**：分布式缓存；缓存版本控制；变更事件通知机制

## 11. 扩展功能

1. **权限委派**：
   - 支持临时权限委派
   - 基于时间或操作次数的限制

2. **权限组**：
   - 将相关权限组合管理
   - 简化权限分配

3. **权限模板**：
   - 预定义权限模板
   - 快速应用到新租户

4. **动态权限规则**：
   - 支持运行时规则定义
   - 基于表达式的权限计算

5. **权限分析工具**：
   - 权限使用分析
   - 权限冗余检测
   - 权限风险评估

## 12. 参考资源

- [RBAC模型规范](https://csrc.nist.gov/CSRC/media/Publications/conference-paper/1992/10/13/proceedings-15th-national-computer-security-conference-1992/documents/1992-15th-NCSC-proceedings-vol-2.pdf)
- [OAuth 2.0规范](https://oauth.net/2/)
- [JWT规范](https://jwt.io/)
- [ABAC模型介绍](https://nvlpubs.nist.gov/nistpubs/specialpublications/NIST.sp.800-162.pdf)
- [Casbin权限库](https://github.com/casbin/casbin)
