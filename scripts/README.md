# 脚本工具目录

此目录包含用于项目构建、测试和部署的各种脚本工具。

## 目录结构

```
scripts/
├── tools/           # 工具源码
├── bin/             # 编译后的工具二进制文件
├── test.sh          # 测试执行脚本
├── run-test.sh      # 交互式测试执行脚本
├── swagger-update.sh # Swagger文档更新脚本
└── sync-swagger-to-yapi.js # 将Swagger同步到YAPI的脚本
```

## 脚本说明

### run-test.sh

交互式测试执行脚本，提供友好的菜单界面，简化测试执行。

```bash
# 执行交互式测试
./run-test.sh

# 命令行参数
./run-test.sh --type=api     # 执行API测试
./run-test.sh --type=coverage # 执行覆盖率测试
./run-test.sh --type=all     # 执行所有测试
./run-test.sh --view         # 执行后查看报告
```

### test.sh

底层测试执行脚本，通常由run-test.sh或Makefile调用。

```bash
# 执行API测试
./test.sh api

# 执行覆盖率测试
./test.sh coverage

# 启动报告服务器
./test.sh coverage serve
```

### swagger-update.sh

更新Swagger API文档。

```bash
# 更新Swagger文档
./swagger-update.sh
```

### sync-swagger-to-yapi.js

将Swagger文档同步到YAPI平台。

```bash
# 同步Swagger到YAPI
node sync-swagger-to-yapi.js
```

## 工具开发

如需添加新工具，请遵循以下步骤：

1. 将源码放在`tools`目录
2. 编译生成的二进制文件放在`bin`目录
3. 更新本README文件，添加新工具的使用说明
