# 项目文档

本目录包含与项目相关的所有文档，包括设计文档、开发指南、测试说明等。

## 目录结构

```
docs/
├── milestones/           # 项目里程碑计划
│   ├── milestone1-foundation-plan.md         # 第一阶段：基础架构
│   ├── milestone2-core-features-plan.md      # 第二阶段：核心功能
│   ├── milestone3-business-modules-plan.md   # 第三阶段：业务模块
│   └── milestone4-optimization-plan.md       # 第四阶段：优化
├── CONTRIBUTING.md       # 贡献指南
├── development.md        # 开发指南
├── testing-framework.md  # 测试框架说明
├── permission-system.md  # 权限系统设计
├── multi-tenant.md       # 多租户设计
├── code-generator.md     # 代码生成器说明
├── yapi-integration.md   # YAPI集成说明
└── *.md                  # 其他文档
```

> 注意：Swagger API文档已迁移至 `ky-admin/swagger` 目录

## 文档更新

### Swagger文档更新

使用以下命令更新Swagger文档：

```bash
cd kysion.com/backend
./scripts/swagger-update.sh
```

### 文档贡献指南

1. 所有文档应使用Markdown格式编写
2. 图表和流程图可以使用Mermaid或PlantUML
3. 文档应保持更新，与代码同步
4. 添加新功能时应同步更新相关文档

## API文档访问

Swagger API文档可通过以下方式访问：

1. 启动服务后访问 <http://localhost:8080/swagger/index.html>
2. 使用 `./scripts/swagger-update.sh` 生成最新文档
3. 查看 `ky-admin/swagger` 目录下的JSON或YAML文件
