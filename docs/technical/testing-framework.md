# KY-Admin 测试框架设计

## 1. 简介

本文档描述了KY-Admin项目的测试框架设计，包括测试目标、测试类型、技术选择、目录结构和实现策略。测试框架旨在确保代码质量、减少缺陷、提高系统稳定性和可维护性。

## 2. 测试目标

- 确保各个组件按照预期工作
- 验证API接口的功能正确性
- 检测性能瓶颈和资源使用情况
- 验证系统在各种条件下的稳定性
- 提高代码覆盖率，减少潜在缺陷
- 为代码重构提供安全保障
- 形成可自动化执行的测试套件

## 3. 测试类型

### 3.1 单元测试

- **目标**：测试独立的代码单元（函数、方法、类）
- **范围**：所有核心业务逻辑、工具函数、服务组件
- **工具**：Go标准测试框架 + testify
- **关注点**：功能正确性、边界条件、错误处理

### 3.2 集成测试

- **目标**：测试多个组件之间的交互
- **范围**：数据库交互、服务间调用、中间件链
- **工具**：Docker容器化测试环境
- **关注点**：组件间协作、数据流、事务处理

### 3.3 API测试

- **目标**：验证API接口的功能和行为
- **范围**：所有公开API端点
- **工具**：httptest包 + testify
- **关注点**：请求处理、响应格式、状态码、权限控制

### 3.4 性能测试

- **目标**：评估系统在不同负载下的性能
- **范围**：关键API、数据库操作、资源密集型处理
- **工具**：Go标准基准测试 + 自定义工具
- **关注点**：响应时间、吞吐量、资源消耗

### 3.5 端到端测试

- **目标**：测试整个系统的流程
- **范围**：核心业务流程、用户场景
- **工具**：专用测试脚本
- **关注点**：用户体验、业务流程完整性

## 4. 技术选择

### 4.1 测试框架

- **Go标准测试包**：使用Go内置的`testing`包作为基础
- **Testify**：提供扩展的断言、模拟和套件功能
  - `assert`：提供丰富的断言方法
  - `require`：断言失败时立即终止测试
  - `mock`：创建和使用模拟对象
  - `suite`：组织测试套件

### 4.2 辅助工具

- **httptest**：用于HTTP请求测试
- **sqlmock**：用于数据库操作模拟
- **gomock**：用于接口模拟
- **Docker Compose**：创建隔离的测试环境
- **GitHub Actions**：自动化测试执行

## 5. 目录结构

```
ky-admin/
├── test/                   # 测试相关代码
│   ├── helpers/            # 测试辅助函数
│   ├── mocks/              # 模拟对象
│   ├── integration/        # 集成测试
│   ├── api/                # API测试
│   ├── performance/        # 性能测试
│   └── e2e/                # 端到端测试
├── testdata/               # 测试数据
├── scripts/
│   ├── test.sh             # 测试执行脚本
│   └── docker-compose.test.yml # 测试环境配置
└── .github/workflows/
    └── test.yml            # CI测试工作流
```

## 6. 单元测试实现

### 6.1 测试文件命名

- 测试文件应命名为 `xxx_test.go`
- 测试函数应命名为 `TestXxx`

### 6.2 测试函数结构

```go
func TestFunction(t *testing.T) {
    // 准备测试数据
    ...
    
    // 执行被测试的函数
    result := FunctionUnderTest(args)
    
    // 验证结果
    assert.Equal(t, expected, result)
}
```

### 6.3 表驱动测试

```go
func TestCalculator(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"正数相加", 2, 3, 5},
        {"负数相加", -2, -3, -5},
        {"零值处理", 0, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 6.4 子测试

```go
func TestComplexFunction(t *testing.T) {
    t.Run("正常情况", func(t *testing.T) {
        // 测试正常情况
    })
    
    t.Run("错误处理", func(t *testing.T) {
        // 测试错误处理
    })
    
    t.Run("边界条件", func(t *testing.T) {
        // 测试边界条件
    })
}
```

### 6.5 Mock对象使用

```go
func TestServiceWithMock(t *testing.T) {
    // 创建Mock对象
    mockRepo := new(mocks.Repository)
    
    // 设置预期行为
    mockRepo.On("FindByID", 1).Return(entity.User{ID: 1, Name: "测试用户"}, nil)
    
    // 创建使用Mock的服务
    service := NewService(mockRepo)
    
    // 执行测试
    user, err := service.GetUser(1)
    
    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, "测试用户", user.Name)
    
    // 验证Mock被正确调用
    mockRepo.AssertExpectations(t)
}
```

### 6.6 测试辅助函数

使用`test/helpers`包中的辅助函数简化测试代码：

```go
func TestDatabase(t *testing.T) {
    // 初始化测试环境
    _, db := helpers.SetupTestEnvironment(t)
    defer helpers.CleanupTestEnvironment()
    
    // 执行数据库相关测试
    ...
}
```

## 7. 集成测试策略

### 7.1 测试环境

使用Docker Compose创建隔离的测试环境，包含：

- MySQL数据库
- Redis缓存
- 其他必要的外部依赖

### 7.2 测试数据

- 使用专用的测试数据库
- 每次测试前重置数据库状态
- 使用事务确保测试隔离

### 7.3 集成测试示例

```go
func TestUserRepository(t *testing.T) {
    // 初始化测试环境
    cfg, db := helpers.SetupTestEnvironment(t)
    defer helpers.CleanupTestEnvironment()
    
    // 准备测试数据
    helpers.ExecuteInTransaction(t, db, func(tx *gorm.DB) {
        // 创建测试数据
        tx.Create(&models.User{Name: "集成测试用户"})
        
        // 初始化仓库
        repo := repository.NewUserRepository(tx)
        
        // 执行测试
        users, err := repo.FindByName("集成测试用户")
        
        // 验证结果
        assert.NoError(t, err)
        assert.Len(t, users, 1)
        assert.Equal(t, "集成测试用户", users[0].Name)
    })
}
```

## 8. API测试

### 8.1 测试HTTP处理器

使用`httptest`包测试HTTP处理器：

```go
func TestUserHandler(t *testing.T) {
    // 设置Gin测试模式
    gin.SetMode(gin.TestMode)
    
    // 创建测试路由
    router := gin.New()
    
    // 注册处理器
    userHandler := handler.NewUserHandler(mockService)
    router.GET("/users/:id", userHandler.GetUser)
    
    // 创建测试请求
    req, _ := http.NewRequest("GET", "/users/1", nil)
    w := httptest.NewRecorder()
    
    // 执行请求
    router.ServeHTTP(w, req)
    
    // 验证响应
    assert.Equal(t, http.StatusOK, w.Code)
    
    // 解析响应
    var resp response.Response
    json.Unmarshal(w.Body.Bytes(), &resp)
    
    // 验证响应内容
    assert.Equal(t, 0, resp.Code)
    assert.Equal(t, "success", resp.Message)
}
```

### 8.2 API测试场景

- 正常请求处理
- 参数验证和错误处理
- 权限控制和认证
- 边界条件和特殊输入
- 响应格式和状态码

## 9. 性能测试

### 9.1 基准测试

使用Go的基准测试功能：

```go
func BenchmarkDatabaseQuery(b *testing.B) {
    // 初始化测试环境
    cfg, db := helpers.SetupBenchmarkEnvironment()
    defer helpers.CleanupBenchmarkEnvironment()
    
    // 准备测试数据
    // ...
    
    // 重置计时器
    b.ResetTimer()
    
    // 执行基准测试
    for i := 0; i < b.N; i++ {
        // 执行被测试的操作
        db.Where("status = ?", "active").Find(&users)
    }
}
```

### 9.2 负载测试

创建专用的负载测试脚本，测试系统在不同并发级别下的性能。

## 10. 测试覆盖率目标

- **单元测试**：核心模块覆盖率 ≥ 80%
- **集成测试**：关键业务流程覆盖率 ≥ 70%
- **API测试**：所有公开API端点覆盖率 100%
- **总体覆盖率**：≥ 75%

## 11. CI/CD集成

### 11.1 GitHub Actions配置

使用GitHub Actions自动化测试流程：

- 每次提交到主分支执行完整测试套件
- 每次提交PR执行相关测试
- 收集并报告测试覆盖率
- 测试失败时阻止合并

### 11.2 本地测试执行

提供便捷的测试执行脚本：

```bash
# 运行所有测试
./scripts/test.sh

# 运行特定类型的测试
./scripts/test.sh [basic|api|integration]
```

## 12. 最佳实践

1. **测试独立性**：每个测试应该独立且不依赖其他测试的状态
2. **可重复性**：测试应该在任何环境中产生相同的结果
3. **测试可读性**：测试应该清晰表达意图和预期行为
4. **测试维护**：随代码变化及时更新测试
5. **测试速度**：单元测试应该快速执行，集成测试可以适当慢一些
6. **测试数据分离**：使用专用的测试数据，避免污染生产数据
7. **测试自动化**：确保测试可以自动执行，便于CI/CD集成

## 13. 实施步骤

1. 配置基础测试环境
2. 建立测试目录结构和辅助工具
3. 为核心模块编写单元测试
4. 实现模拟对象和测试数据生成器
5. 配置集成测试环境
6. 为API端点编写测试
7. 设置性能测试和基准测试
8. 配置CI/CD集成
9. 建立测试覆盖率监控

## 14. 总结

KY-Admin的测试框架旨在通过全面的测试策略确保代码质量和系统稳定性。通过结合多种测试类型和自动化工具，我们可以在开发过程中尽早发现并解决问题，提高项目的整体质量。
