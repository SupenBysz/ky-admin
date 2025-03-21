# KY-Admin 后端测试系统

## 目录结构

```
backend/
├── ky-admin/            # 后台管理系统核心代码
│   ├── swagger/         # Swagger API文档
│   ├── cmd/             # 主程序入口
│   ├── pkg/             # 核心功能包
│   └── internal/        # 内部实现
├── tests/               # 测试目录
│   ├── api/             # API测试
│   ├── integration/     # 集成测试
│   ├── unit/            # 单元测试
│   ├── performance/     # 性能测试
│   ├── fixtures/        # 测试数据和资源
│   ├── mocks/           # 模拟对象
│   ├── reports/         # 测试报告输出目录
│   ├── config.yaml      # 测试配置文件
│   └── Makefile         # 测试专用Makefile
├── docs/                # 项目文档
│   ├── milestones/      # 项目里程碑计划
│   └── *.md             # 各类文档文件
├── docker/              # Docker配置和数据
│   ├── data/            # Docker数据目录
│   ├── docker-compose.test.yml  # 测试环境配置
│   └── docker-compose.yapi.yml  # YAPI环境配置
├── scripts/             # 工具脚本
│   ├── tools/           # 工具源码
│   ├── bin/             # 编译后工具
│   ├── test.sh          # 测试执行脚本
│   ├── run-test.sh      # 交互式测试执行脚本
│   ├── swagger-update.sh # Swagger文档更新脚本
│   └── sync-swagger-to-yapi.js # 将Swagger同步到YAPI的脚本
└── Makefile             # 项目管理
```

## 测试执行方法

### 使用交互式脚本

```bash
# 进入backend目录
cd kysion.com/backend

# 执行测试脚本
./scripts/run-test.sh
```

### 使用Makefile

```bash
# 进入backend目录
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

### 直接在tests目录使用

```bash
# 进入tests目录
cd kysion.com/backend/tests

# 运行API测试
make test-api

# 查看帮助
make help
```

## Docker环境

项目提供了两种Docker环境配置：

### 测试环境

```bash
cd kysion.com/backend/docker
docker-compose -f docker-compose.test.yml up -d
```

### YAPI接口管理平台

```bash
cd kysion.com/backend/docker
docker-compose -f docker-compose.yapi.yml up -d
```

详细信息请参阅 [docker/README.md](docker/README.md)。

## 测试配置

测试配置集中在 `tests/config.yaml` 文件中管理，包括：

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

测试报告保存在 `tests/reports` 目录，每次测试会创建一个时间戳目录。
可以通过以下方式快速访问最新报告：

1. 使用 `./scripts/run-test.sh` 选择 "仅查看报告" 选项
2. 或直接访问 `http://localhost:8089` (需启动报告服务器)
3. 或通过 `make serve-report` 启动报告服务器

## API文档管理

### 更新Swagger文档

使用以下命令更新Swagger文档：

```bash
cd kysion.com/backend
./scripts/swagger-update.sh
```

生成的Swagger文档位于 `ky-admin/swagger` 目录。

### 同步Swagger到YAPI

使用以下命令将Swagger文档同步到YAPI平台：

```bash
cd kysion.com/backend
node ./scripts/sync-swagger-to-yapi.js
```

## 项目文档

项目文档位于 `docs` 目录，包含以下内容：

- 项目里程碑计划
- 开发指南
- 贡献指南
- 功能模块设计文档

详细信息请参阅 [docs/README.md](docs/README.md)。
