# KY-Admin 测试文档

本目录包含KY-Admin项目的测试代码和相关资源。

## 目录结构

```
tests/
├── api/           # API集成测试
├── fixtures/      # 测试固定数据
├── integration/   # 集成测试
├── mocks/         # 模拟对象
├── performance/   # 性能测试
├── reports/       # 测试报告模板（实际报告输出到独立目录）
└── unit/          # 单元测试
```

## 测试报告

所有测试报告都会输出到项目的统一位置：

```
/backend/temp/test-reports/
```

每次测试运行都会创建一个带有时间戳的目录，便于追踪不同测试运行的结果。最新的测试报告可以通过符号链接访问：

```
/backend/temp/test-reports/latest
```

## 运行测试

可以使用Makefile中定义的命令来运行测试：

```bash
# 进入测试目录
cd ky-admin/tests

# 运行所有测试
make test-all

# 运行API测试
make test-api

# 运行并查看覆盖率报告
make test-coverage

# 查看测试帮助
make help
```

也可以直接使用测试脚本：

```bash
# 使用交互式菜单
./scripts/run-test.sh

# 指定测试类型
./scripts/run-test.sh --type=all
./scripts/run-test.sh --type=api
./scripts/run-test.sh --type=coverage

# 生成并查看报告
./scripts/run-test.sh --type=coverage --view
```

## 测试类型

1. **单元测试**: 测试独立组件的功能
2. **集成测试**: 测试组件间的交互
3. **API测试**: 测试API端点的功能
4. **性能测试**: 测试系统在负载下的性能

## 编写测试

编写测试时，请遵循以下原则：

1. 测试文件命名为 `xxx_test.go`
2. 测试函数命名为 `TestXxx`
3. 单元测试应该是独立的、快速的
4. 使用模拟对象替代外部依赖
5. 使用断言库简化测试代码

## CI/CD集成

测试会在GitHub Actions CI/CD工作流中自动运行。测试报告会被保存为构建产物，可以在GitHub Actions界面查看。

测试覆盖率数据会被用来生成徽章，显示在项目的README.md中。
