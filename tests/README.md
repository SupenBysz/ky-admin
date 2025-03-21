# 测试目录结构说明

本项目采用统一的测试目录结构，将所有测试相关内容集中管理。

## 目录结构

```
tests/
├── api/            # API测试
├── integration/    # 集成测试
├── unit/           # 单元测试
├── performance/    # 性能测试
├── fixtures/       # 测试数据和资源
├── mocks/          # 模拟对象
├── reports/        # 测试报告输出目录
│   └── latest/     # 指向最新测试报告的符号链接
├── config.yaml     # 测试配置文件
└── Makefile        # 测试专用Makefile
```

## 测试执行方法

### 使用交互式脚本

项目提供了交互式测试脚本，使用非常简单：

```bash
# 进入项目根目录
cd kysion.com/backend

# 执行测试脚本
./scripts/run-test.sh
```

### 使用项目Makefile

项目根目录提供了Makefile简化测试执行：

```bash
# 在项目根目录
cd kysion.com/backend

# 运行API测试
make test-api

# 运行覆盖率测试
make test-coverage

# 运行所有测试
make test-all

# 生成测试报告
make test-report

# 查看测试报告
make serve-report

# 清理测试报告
make clean

# 查看帮助信息
make help
```

### 直接使用tests目录的Makefile

tests目录也提供了独立的Makefile，可以直接在tests目录下执行：

```bash
# 进入tests目录
cd kysion.com/backend/tests

# 运行API测试
make test-api

# 运行覆盖率测试
make test-coverage

# 运行所有测试
make test-all

# 生成测试报告
make test-report

# 查看测试报告
make serve-report

# 清理测试报告
make clean

# 查看帮助信息
make help
```

## 测试配置

测试配置集中在 `config.yaml` 文件中管理，包括：

- 路径配置
- 服务器配置
- 覆盖率阈值
- 数据库配置
- 测试类型定义

## 添加新测试

1. 根据测试类型将测试文件放在对应目录
2. 测试文件命名规范：`*_test.go`
3. 如需添加测试数据，放在`fixtures`目录下
4. 如需添加测试辅助函数，放在对应测试目录下的 `helpers.go` 文件中

## 查看测试报告

测试报告保存在 `reports` 目录，每次测试会创建一个时间戳目录。
可以通过以下方式快速访问最新报告：

1. 使用 `../scripts/run-test.sh` 选择 "仅查看报告" 选项
2. 或直接访问 `http://localhost:8089` (需启动报告服务器)
3. 或通过 `make serve-report` 启动报告服务器
