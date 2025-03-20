# 多租户系统设计与实现指南

## 1. 概述

多租户系统是KY-Admin框架的重要特性之一，支持在同一套系统上服务多个独立的客户组织（租户），同时保证数据隔离和安全性。本文档详细介绍了KY-Admin多租户系统的设计理念、实现方式和使用方法。

## 2. 设计理念

KY-Admin多租户系统遵循以下设计理念：

- **数据隔离**：确保不同租户之间的数据严格隔离，防止数据泄露
- **资源共享**：多个租户共享同一套系统资源，降低运维成本
- **扩展性**：支持灵活添加新租户，且不影响现有租户正常使用
- **定制化**：允许针对不同租户进行一定程度的个性化定制
- **统一管理**：提供统一的租户管理界面，简化管理员操作

## 3. 多租户架构模式

KY-Admin支持以下三种多租户架构模式：

### 3.1 独立数据库模式

每个租户使用独立的数据库实例，实现完全的数据隔离。

**优点**：

- 最高级别的数据隔离
- 容易备份和恢复特定租户的数据
- 租户间性能互不影响

**缺点**：

- 资源消耗较大
- 维护和升级成本高
- 跨租户操作复杂

### 3.2 共享数据库，独立Schema模式

所有租户共享同一个数据库实例，但每个租户使用独立的数据库Schema。

**优点**：

- 良好的数据隔离性
- 维护成本适中
- 支持部分数据库原生功能

**缺点**：

- 需要数据库支持Schema
- Schema数量可能受限
- 备份恢复相对复杂

### 3.3 共享数据库，共享Schema模式（默认模式）

所有租户共享同一个数据库实例和Schema，通过租户ID字段区分不同租户的数据。

**优点**：

- 资源利用率最高
- 维护成本最低
- 可支持最多租户数量

**缺点**：

- 数据隔离依赖于应用层
- 可能存在性能瓶颈
- 需要额外的安全保障措施

## 4. 数据模型设计

### 4.1 租户表

```sql
CREATE TABLE `sys_tenant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '租户ID',
  `name` varchar(50) NOT NULL COMMENT '租户名称',
  `domain` varchar(50) DEFAULT NULL COMMENT '租户域名',
  `code` varchar(20) NOT NULL COMMENT '租户编码',
  `contact_user` varchar(50) DEFAULT NULL COMMENT '联系人',
  `contact_phone` varchar(20) DEFAULT NULL COMMENT '联系电话',
  `expire_time` datetime DEFAULT NULL COMMENT '过期时间',
  `license_key` varchar(100) DEFAULT NULL COMMENT '授权码',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户信息表';
```

### 4.2 租户配置表

```sql
CREATE TABLE `sys_tenant_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `tenant_id` bigint(20) NOT NULL COMMENT '租户ID',
  `config_key` varchar(50) NOT NULL COMMENT '配置键',
  `config_value` text COMMENT '配置值',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_key` (`tenant_id`,`config_key`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户配置表';
```

### 4.3 租户数据库信息表（用于独立数据库模式）

```sql
CREATE TABLE `sys_tenant_database` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `tenant_id` bigint(20) NOT NULL COMMENT '租户ID',
  `db_type` varchar(20) NOT NULL COMMENT '数据库类型',
  `db_host` varchar(100) NOT NULL COMMENT '主机地址',
  `db_port` int(11) NOT NULL COMMENT '端口',
  `db_name` varchar(50) NOT NULL COMMENT '数据库名',
  `db_username` varchar(50) NOT NULL COMMENT '用户名',
  `db_password` varchar(100) NOT NULL COMMENT '密码',
  `db_params` varchar(255) DEFAULT NULL COMMENT '连接参数',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户数据库信息表';
```

### 4.4 业务表设计

在共享数据库模式下，所有业务表需要增加租户ID字段：

```sql
-- 示例：用户表
CREATE TABLE `sys_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `tenant_id` bigint(20) NOT NULL COMMENT '租户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  -- 其他字段
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username_tenant` (`username`,`tenant_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';
```

## 5. 核心功能实现

### 5.1 租户识别

KY-Admin支持多种租户识别方式：

1. **URL路径识别**：`https://example.com/{tenant_code}/api/xxx`
2. **子域名识别**：`https://{tenant_code}.example.com/api/xxx`
3. **请求参数识别**：`https://example.com/api/xxx?tenant=tenant_code`
4. **请求头识别**：`X-Tenant-Code: tenant_code`
5. **JWT Token识别**：从Token中解析租户信息

实现示例：

```go
// 租户识别中间件
func TenantIdentifyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        var tenantCode string
        
        // 尝试从请求头获取
        tenantCode = c.GetHeader("X-Tenant-Code")
        
        // 尝试从子域名获取
        if tenantCode == "" {
            host := c.Request.Host
            parts := strings.Split(host, ".")
            if len(parts) >= 3 {
                tenantCode = parts[0]
            }
        }
        
        // 尝试从URL路径获取
        if tenantCode == "" {
            path := c.Request.URL.Path
            pathParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
            if len(pathParts) > 0 {
                // 检查第一段是否为租户代码
                possibleTenantCode := pathParts[0]
                // 验证是否为有效租户代码
                if isTenantCodeValid(possibleTenantCode) {
                    tenantCode = possibleTenantCode
                    // 重写路径，去除租户代码部分
                    c.Request.URL.Path = "/" + strings.Join(pathParts[1:], "/")
                }
            }
        }
        
        // 尝试从请求参数获取
        if tenantCode == "" {
            tenantCode = c.Query("tenant")
        }
        
        // 尝试从JWT Token获取
        if tenantCode == "" {
            token := c.GetHeader("Authorization")
            if token != "" {
                // 解析JWT Token
                claims, err := parseJWTToken(token)
                if err == nil {
                    tenantCode = claims.TenantCode
                }
            }
        }
        
        // 如果找到租户代码，查询租户信息并设置到上下文
        if tenantCode != "" {
            tenant, err := getTenantByCode(tenantCode)
            if err == nil && tenant.Status == 1 {
                // 将租户信息设置到上下文
                c.Set("tenant", tenant)
                c.Set("tenantId", tenant.ID)
                c.Set("tenantCode", tenant.Code)
            } else {
                // 租户不存在或已禁用
                c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                    "code": 403,
                    "message": "Invalid tenant or tenant disabled",
                })
                return
            }
        } else {
            // 没有找到租户信息，可能是访问公共API
            // 可以设置默认租户或者允许继续访问
        }
        
        c.Next()
    }
}
```

### 5.2 数据过滤

在共享数据库模式下，需要对所有查询自动添加租户过滤条件，确保数据隔离：

```go
// GORM中间件，自动添加租户过滤
func TenantFilterMiddleware(db *gorm.DB) {
    db.Callback().Query().Before("gorm:query").Register("tenant_filter", func(db *gorm.DB) {
        // 获取当前上下文
        ctx := db.Statement.Context
        if ctx == nil {
            return
        }
        
        // 从上下文中获取租户ID
        tenantID, exists := ctx.Value("tenantId").(uint64)
        if !exists {
            return
        }
        
        // 获取当前模型
        if db.Statement.Schema == nil {
            return
        }
        
        // 检查模型是否有租户ID字段
        if field, ok := db.Statement.Schema.FieldsByName["TenantID"]; ok {
            // 添加租户过滤条件
            tableName := db.Statement.Table
            db.Where(fmt.Sprintf("%s.tenant_id = ?", tableName), tenantID)
        }
    })
}
```

### 5.3 租户管理

提供完整的租户管理功能，包括：

1. **租户创建**：初始化租户数据和配置
2. **租户配置**：管理租户信息和参数设置
3. **租户启停**：控制租户的启用和禁用状态
4. **租户删除**：删除租户数据（通常是逻辑删除）

```go
// 租户服务接口
type TenantService interface {
    // 创建租户
    Create(ctx context.Context, req dto.CreateTenantRequest) (*model.Tenant, error)
    
    // 获取租户信息
    Get(ctx context.Context, id uint64) (*model.Tenant, error)
    
    // 更新租户信息
    Update(ctx context.Context, id uint64, req dto.UpdateTenantRequest) (*model.Tenant, error)
    
    // 删除租户
    Delete(ctx context.Context, id uint64) error
    
    // 租户列表查询
    List(ctx context.Context, query dto.TenantQuery) ([]*model.Tenant, int64, error)
    
    // 创建租户初始数据
    InitTenantData(ctx context.Context, tenantId uint64) error
    
    // 启用租户
    Enable(ctx context.Context, id uint64) error
    
    // 禁用租户
    Disable(ctx context.Context, id uint64) error
    
    // 租户配置管理
    GetConfig(ctx context.Context, tenantId uint64, key string) (string, error)
    SetConfig(ctx context.Context, tenantId uint64, key string, value string) error
}
```

### 5.4 租户数据初始化

新创建租户时需要初始化基础数据：

```go
// 租户数据初始化
func (s *TenantServiceImpl) InitTenantData(ctx context.Context, tenantId uint64) error {
    // 1. 创建租户管理员账号
    admin := &model.User{
        TenantID: tenantId,
        Username: "admin",
        Password: crypto.HashPassword("123456"), // 初始密码
        RealName: "管理员",
        Status:   1,
    }
    if err := s.db.Create(admin).Error; err != nil {
        return fmt.Errorf("create tenant admin failed: %w", err)
    }
    
    // 2. 创建基础角色
    adminRole := &model.Role{
        TenantID:    tenantId,
        Name:        "管理员",
        Code:        "admin",
        Description: "租户管理员角色",
        Status:      1,
    }
    if err := s.db.Create(adminRole).Error; err != nil {
        return fmt.Errorf("create admin role failed: %w", err)
    }
    
    userRole := &model.Role{
        TenantID:    tenantId,
        Name:        "普通用户",
        Code:        "user",
        Description: "普通用户角色",
        Status:      1,
    }
    if err := s.db.Create(userRole).Error; err != nil {
        return fmt.Errorf("create user role failed: %w", err)
    }
    
    // 3. 分配角色给管理员
    userRole := &model.UserRole{
        UserID:   admin.ID,
        RoleID:   adminRole.ID,
        TenantID: tenantId,
    }
    if err := s.db.Create(userRole).Error; err != nil {
        return fmt.Errorf("assign role to admin failed: %w", err)
    }
    
    // 4. 初始化基础权限数据
    // ...
    
    // 5. 初始化系统配置
    // ...
    
    return nil
}
```

### 5.5 租户数据源切换（独立数据库模式）

在独立数据库模式下，需要根据租户动态切换数据源：

```go
// 数据源管理
type DataSourceManager struct {
    defaultDB *gorm.DB
    tenantDBs map[uint64]*gorm.DB
    mutex     sync.RWMutex
}

// 获取租户数据库连接
func (m *DataSourceManager) GetDB(ctx context.Context) *gorm.DB {
    // 从上下文获取租户ID
    tenantID, exists := ctx.Value("tenantId").(uint64)
    if !exists {
        // 未找到租户ID，使用默认数据库
        return m.defaultDB
    }
    
    // 尝试获取已存在的连接
    m.mutex.RLock()
    db, exists := m.tenantDBs[tenantID]
    m.mutex.RUnlock()
    
    if exists {
        return db
    }
    
    // 创建新连接
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    // 再次检查，防止并发创建
    db, exists = m.tenantDBs[tenantID]
    if exists {
        return db
    }
    
    // 获取租户数据库信息
    var tenantDB model.TenantDatabase
    if err := m.defaultDB.Where("tenant_id = ?", tenantID).First(&tenantDB).Error; err != nil {
        // 找不到租户数据库信息，使用默认数据库
        log.Printf("Tenant database not found for tenant ID %d, using default DB", tenantID)
        return m.defaultDB
    }
    
    // 构建DSN
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
        tenantDB.DBUsername,
        tenantDB.DBPassword,
        tenantDB.DBHost,
        tenantDB.DBPort,
        tenantDB.DBName,
        tenantDB.DBParams,
    )
    
    // 创建连接
    newDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("Failed to connect to tenant database: %v", err)
        return m.defaultDB
    }
    
    // 存储连接
    m.tenantDBs[tenantID] = newDB
    return newDB
}
```

## 6. 使用指南

### 6.1 租户中间件配置

在应用入口处配置租户中间件：

```go
func main() {
    // 初始化Gin
    r := gin.Default()
    
    // 注册中间件
    r.Use(middleware.TenantIdentifyMiddleware())
    
    // 配置路由
    setupRoutes(r)
    
    // 启动服务
    r.Run(":8080")
}
```

### 6.2 数据访问层使用

在服务层代码中传递上下文：

```go
func (s *UserServiceImpl) List(ctx context.Context, query dto.UserQuery) ([]*model.User, int64, error) {
    // 获取带有租户信息的数据库连接
    db := s.db.WithContext(ctx)
    
    // 构建查询
    // 无需手动添加租户过滤，中间件会自动处理
    db = db.Model(&model.User{})
    
    // 应用其他查询条件
    if query.Username != "" {
        db = db.Where("username LIKE ?", "%"+query.Username+"%")
    }
    
    // 查询结果
    var count int64
    db.Count(&count)
    
    var users []*model.User
    db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).Find(&users)
    
    return users, count, nil
}
```

### 6.3 租户管理控制台

实现租户管理界面，提供以下功能：

1. 创建新租户
2. 管理租户信息和状态
3. 配置租户参数
4. 查看租户使用统计

### 6.4 租户隔离测试

编写测试用例验证租户隔离性：

```go
func TestTenantDataIsolation(t *testing.T) {
    // 创建测试数据库和两个测试租户
    db := setupTestDB()
    tenant1 := createTestTenant(db, "tenant1")
    tenant2 := createTestTenant(db, "tenant2")
    
    // 在租户1下创建数据
    ctx1 := context.WithValue(context.Background(), "tenantId", tenant1.ID)
    user1 := createTestUser(db.WithContext(ctx1), "user1", tenant1.ID)
    
    // 在租户2下创建数据
    ctx2 := context.WithValue(context.Background(), "tenantId", tenant2.ID)
    user2 := createTestUser(db.WithContext(ctx2), "user2", tenant2.ID)
    
    // 测试租户1只能查看自己的数据
    var users1 []model.User
    db.WithContext(ctx1).Find(&users1)
    assert.Equal(t, 1, len(users1))
    assert.Equal(t, user1.ID, users1[0].ID)
    
    // 测试租户2只能查看自己的数据
    var users2 []model.User
    db.WithContext(ctx2).Find(&users2)
    assert.Equal(t, 1, len(users2))
    assert.Equal(t, user2.ID, users2[0].ID)
}
```

## 7. 最佳实践

1. **选择合适的租户模式**：
   - 小型应用优先考虑共享数据库模式
   - 对安全要求高的应用考虑独立数据库模式
   - 需要平衡性能和隔离性的应用考虑共享数据库、独立Schema模式

2. **设计时考虑租户**：
   - 所有业务表设计时都要考虑租户字段
   - 所有查询操作需考虑租户过滤
   - 跨租户操作需特别谨慎

3. **租户标识存储**：
   - 将租户标识存储在上下文中
   - 确保租户标识贯穿整个请求生命周期
   - 不要硬编码租户ID

4. **超级管理员模式**：
   - 设计超级管理员角色可跨租户操作
   - 明确超级管理员的权限边界
   - 记录超级管理员的操作日志

5. **性能优化**：
   - 租户字段添加索引
   - 设计合理的分表分库策略
   - 考虑租户级别的缓存隔离

## 8. 常见问题

1. **问题**：如何处理公共数据？  
   **解决**：可以设置特殊租户ID（如0）表示公共数据，或使用专门的公共数据表

2. **问题**：租户数量增长导致性能下降  
   **解决**：实施分库分表策略；对大租户使用独立数据库；优化查询和索引

3. **问题**：跨租户数据查询需求  
   **解决**：为超级管理员提供特殊API；使用数据仓库汇总分析；实现受控的跨租户查询机制

4. **问题**：租户配置如何影响系统行为  
   **解决**：设计租户配置缓存；实现动态特性开关；使用策略模式适应不同配置

5. **问题**：如何迁移租户数据  
   **解决**：设计数据导入导出工具；实现租户数据备份机制；支持跨租户数据复制

## 9. 扩展功能

1. **租户资源限额**：
   - 限制租户用户数量
   - 限制存储空间使用
   - 限制API调用频率

2. **租户数据统计**：
   - 数据量统计
   - 用户活跃度统计
   - 资源使用率监控

3. **租户生命周期管理**：
   - 试用期设置
   - 到期提醒和续费
   - 数据备份和归档

4. **租户个性化**：
   - 自定义主题和品牌
   - 自定义字段和表单
   - 功能模块开关

5. **多级租户**：
   - 支持租户下设子租户
   - 实现租户组织结构
   - 跨子租户数据共享

## 10. 参考资源

- [多租户SaaS应用架构最佳实践](https://docs.microsoft.com/en-us/azure/architecture/guide/multitenant/overview)
- [GORM文档](https://gorm.io/docs/)
- [Gin Web框架](https://github.com/gin-gonic/gin)
- [多租户数据库模式](https://docs.microsoft.com/en-us/azure/sql-database/saas-tenancy-app-design-patterns)
- [云原生多租户应用开发指南](https://www.nginx.com/blog/microservices-reference-architecture-nginx-multi-tenancy-patterns/)
