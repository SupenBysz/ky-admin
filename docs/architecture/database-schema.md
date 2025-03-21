# 数据库模式设计

本文档描述了KY-Admin系统的数据库设计，包括主要表结构、关系和索引设计。

## 数据库概述

KY-Admin使用MySQL 8.0作为主要关系型数据库，采用UTF8MB4字符集和InnoDB存储引擎。数据库设计遵循第三范式，同时针对查询性能进行了适当的反规范化设计。

## 核心表结构

### 用户认证与授权

#### 用户表 (users)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| username | VARCHAR(50) | 否 | | 用户名，唯一索引 |
| password | VARCHAR(100) | 否 | | 加密后的密码 |
| nickname | VARCHAR(50) | 是 | NULL | 用户昵称 |
| email | VARCHAR(100) | 是 | NULL | 电子邮箱，唯一索引 |
| phone | VARCHAR(20) | 是 | NULL | 手机号，唯一索引 |
| avatar | VARCHAR(255) | 是 | NULL | 头像URL |
| status | TINYINT | 否 | 1 | 状态：0-禁用，1-启用 |
| tenant_id | BIGINT | 否 | 0 | 租户ID |
| last_login | DATETIME | 是 | NULL | 最后登录时间 |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | DATETIME | 是 | NULL | 删除时间（软删除） |

索引：

- PRIMARY KEY (id)
- UNIQUE INDEX idx_username (username, deleted_at)
- UNIQUE INDEX idx_email (email, deleted_at)
- UNIQUE INDEX idx_phone (phone, deleted_at)
- INDEX idx_tenant (tenant_id)

#### 角色表 (roles)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| name | VARCHAR(50) | 否 | | 角色名称 |
| code | VARCHAR(50) | 否 | | 角色编码 |
| description | VARCHAR(200) | 是 | NULL | 角色描述 |
| status | TINYINT | 否 | 1 | 状态：0-禁用，1-启用 |
| tenant_id | BIGINT | 否 | 0 | 租户ID |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | DATETIME | 是 | NULL | 删除时间（软删除） |

索引：

- PRIMARY KEY (id)
- UNIQUE INDEX idx_name_tenant (name, tenant_id, deleted_at)
- UNIQUE INDEX idx_code_tenant (code, tenant_id, deleted_at)
- INDEX idx_tenant (tenant_id)

#### 权限表 (permissions)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| name | VARCHAR(50) | 否 | | 权限名称 |
| code | VARCHAR(50) | 否 | | 权限编码 |
| module | VARCHAR(50) | 否 | | 所属模块 |
| action | VARCHAR(50) | 否 | | 操作类型（例如：create, read, update, delete） |
| description | VARCHAR(200) | 是 | NULL | 权限描述 |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | DATETIME | 是 | NULL | 删除时间（软删除） |

索引：

- PRIMARY KEY (id)
- UNIQUE INDEX idx_code (code, deleted_at)
- INDEX idx_module (module)

#### 用户角色关联表 (user_roles)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| user_id | BIGINT | 否 | | 用户ID，外键 |
| role_id | BIGINT | 否 | | 角色ID，外键 |
| tenant_id | BIGINT | 否 | 0 | 租户ID |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |

索引：

- PRIMARY KEY (id)
- UNIQUE INDEX idx_user_role (user_id, role_id)
- INDEX idx_user (user_id)
- INDEX idx_role (role_id)
- INDEX idx_tenant (tenant_id)

#### 角色权限关联表 (role_permissions)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| role_id | BIGINT | 否 | | 角色ID，外键 |
| permission_id | BIGINT | 否 | | 权限ID，外键 |
| tenant_id | BIGINT | 否 | 0 | 租户ID |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |

索引：

- PRIMARY KEY (id)
- UNIQUE INDEX idx_role_permission (role_id, permission_id)
- INDEX idx_role (role_id)
- INDEX idx_permission (permission_id)
- INDEX idx_tenant (tenant_id)

### 系统管理

#### 租户表 (tenants)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| name | VARCHAR(100) | 否 | | 租户名称 |
| code | VARCHAR(50) | 否 | | 租户编码，唯一 |
| domain | VARCHAR(100) | 是 | NULL | 租户域名 |
| status | TINYINT | 否 | 1 | 状态：0-禁用，1-启用 |
| expiry_date | DATE | 是 | NULL | 过期日期 |
| contact_name | VARCHAR(50) | 是 | NULL | 联系人姓名 |
| contact_phone | VARCHAR(20) | 是 | NULL | 联系人电话 |
| contact_email | VARCHAR(100) | 是 | NULL | 联系人邮箱 |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | DATETIME | 是 | NULL | 删除时间（软删除） |

索引：

- PRIMARY KEY (id)
- UNIQUE INDEX idx_code (code, deleted_at)
- UNIQUE INDEX idx_domain (domain, deleted_at)

#### 部门表 (departments)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| name | VARCHAR(50) | 否 | | 部门名称 |
| parent_id | BIGINT | 是 | 0 | 父部门ID |
| order_num | INT | 否 | 0 | 显示顺序 |
| leader | VARCHAR(50) | 是 | NULL | 负责人 |
| phone | VARCHAR(20) | 是 | NULL | 联系电话 |
| email | VARCHAR(100) | 是 | NULL | 邮箱 |
| status | TINYINT | 否 | 1 | 状态：0-禁用，1-启用 |
| tenant_id | BIGINT | 否 | 0 | 租户ID |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | DATETIME | 是 | NULL | 删除时间（软删除） |

索引：

- PRIMARY KEY (id)
- INDEX idx_parent (parent_id)
- INDEX idx_tenant (tenant_id)

#### 菜单表 (menus)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| name | VARCHAR(50) | 否 | | 菜单名称 |
| parent_id | BIGINT | 是 | 0 | 父菜单ID |
| path | VARCHAR(200) | 是 | NULL | 路由地址 |
| component | VARCHAR(255) | 是 | NULL | 组件路径 |
| permission | VARCHAR(100) | 是 | NULL | 权限标识 |
| type | TINYINT | 否 | 0 | 类型：0-目录，1-菜单，2-按钮 |
| icon | VARCHAR(100) | 是 | NULL | 图标 |
| order_num | INT | 否 | 0 | 显示顺序 |
| visible | TINYINT | 否 | 1 | 是否可见：0-隐藏，1-显示 |
| status | TINYINT | 否 | 1 | 状态：0-禁用，1-启用 |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | DATETIME | 是 | NULL | 删除时间（软删除） |

索引：

- PRIMARY KEY (id)
- INDEX idx_parent (parent_id)

#### 角色菜单关联表 (role_menus)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| role_id | BIGINT | 否 | | 角色ID，外键 |
| menu_id | BIGINT | 否 | | 菜单ID，外键 |
| tenant_id | BIGINT | 否 | 0 | 租户ID |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | 否 | CURRENT_TIMESTAMP | 更新时间 |

索引：

- PRIMARY KEY (id)
- UNIQUE INDEX idx_role_menu (role_id, menu_id)
- INDEX idx_role (role_id)
- INDEX idx_menu (menu_id)
- INDEX idx_tenant (tenant_id)

### 系统监控

#### 操作日志表 (op_logs)

| 字段名 | 类型 | 允许空 | 默认值 | 描述 |
|-------|------|-------|-------|------|
| id | BIGINT | 否 | | 主键，自增ID |
| user_id | BIGINT | 是 | NULL | 操作用户ID |
| username | VARCHAR(50) | 是 | NULL | 操作用户名 |
| ip | VARCHAR(50) | 是 | NULL | 操作IP |
| module | VARCHAR(50) | 否 | | 操作模块 |
| action | VARCHAR(50) | 否 | | 操作类型 |
| description | VARCHAR(200) | 是 | NULL | 操作描述 |
| request_method | VARCHAR(10) | 是 | NULL | 请求方法 |
| request_url | VARCHAR(255) | 是 | NULL | 请求URL |
| request_params | TEXT | 是 | NULL | 请求参数 |
| response_code | INT | 是 | NULL | 响应状态码 |
| response_data | TEXT | 是 | NULL | 响应数据 |
| execution_time | INT | 是 | NULL | 执行时间(毫秒) |
| tenant_id | BIGINT | 否 | 0 | 租户ID |
| created_at | DATETIME | 否 | CURRENT_TIMESTAMP | 创建时间 |

索引：

- PRIMARY KEY (id)
- INDEX idx_user (user_id)
- INDEX idx_module (module)
- INDEX idx_action (action)
- INDEX idx_created (created_at)
- INDEX idx_tenant (tenant_id)

## 数据库关系图

下面是主要表之间的关系图：

```
users 1 ---- N user_roles N ---- 1 roles
      |                           |
      |                           |
      v                           v
departments                 role_permissions
                                  |
                                  |
                                  v
                             permissions
      
roles 1 ---- N role_menus N ---- 1 menus
```

## 数据库分表策略

- **操作日志表**: 按月分表，如`op_logs_202503`
- **业务日志表**: 按用户行为类型分表，如`user_logs`、`payment_logs`等

## 索引策略

- 所有外键关系都建立索引
- 对于频繁查询的字段建立索引
- 对于频繁排序的字段建立索引
- 组合索引遵循最左前缀原则
- 定期维护索引，删除不必要的索引

## 事务管理

系统采用以下事务管理策略：

1. 所有涉及多表操作的业务逻辑使用事务包装
2. 使用声明式事务管理，通过注解或装饰器标记事务边界
3. 事务隔离级别默认使用READ COMMITTED
4. 对于特殊业务场景，可以根据需要调整隔离级别

## 扩展与分区

随着数据量增长，系统计划采取以下扩展策略：

1. 对于大表使用分区表技术，按时间或ID范围分区
2. 使用主从复制实现读写分离
3. 使用分库分表中间件实现水平拆分
4. 对于不常用的历史数据，定期归档到历史表
