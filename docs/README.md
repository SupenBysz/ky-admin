# KY-Admin 项目文档

欢迎来到KY-Admin项目文档。本文档提供了项目的概述、开发指南、技术架构和使用说明。

## 文档目录结构

```
docs/
├── README.md                   # 文档主页
│
├── guide/                      # 使用指南
│   ├── development.md          # 开发指南
│   ├── contributing.md         # 贡献指南
│   └── deployment.md           # 部署指南
│
├── architecture/               # 架构文档
│   ├── overview.md             # 系统架构概述
│   ├── permission-system.md    # 权限系统设计
│   ├── multi-tenant.md         # 多租户设计
│   └── database-schema.md      # 数据库模式
│
├── technical/                  # 技术文档
│   ├── testing-framework.md    # 测试框架
│   ├── code-generator.md       # 代码生成器
│   ├── containerization-cicd.md # 容器化与CI/CD
│   ├── yapi-integration.md     # YAPI集成
│   └── yapi-integration-guide.md # YAPI集成指南
│
├── project/                    # 项目管理
│   ├── milestone-plan.md       # 里程碑计划
│   ├── status.md               # 项目状态
│   ├── backend-progress.md     # 后端总体进度
│   └── subtask-progress.md     # 子任务进度
│
└── releases/                   # 版本发布
    ├── stage1-progress-tracker.md # 阶段1进度跟踪
    ├── v1.0-release-notes.md   # v1.0发布说明
    └── changelog.md            # 变更日志
```

## 文档说明

### 使用指南

- **开发指南**：提供本地开发环境设置、代码规范和工作流程的指导
- **贡献指南**：描述如何为项目做出贡献，包括分支策略、提交规范和PR流程
- **部署指南**：提供在不同环境中部署项目的详细步骤

### 架构文档

- **系统架构概述**：描述系统的整体架构、模块组成和交互关系
- **权限系统设计**：详细说明RBAC权限系统的实现方式
- **多租户设计**：多租户系统的架构和数据隔离方案
- **数据库模式**：数据库表结构和关系

### 技术文档

- **测试框架**：描述测试策略、工具和实践
- **代码生成器**：代码生成工具的使用说明
- **容器化与CI/CD**：Docker、Kubernetes配置和CI/CD流程
- **YAPI集成**：API文档集成的实现细节和使用说明

### 项目管理

- **里程碑计划**：项目总体规划和里程碑定义
- **项目状态**：当前项目进展和健康状况
- **后端总体进度**：后端开发的整体进度跟踪
- **子任务进度**：详细的子任务分解和进度

### 版本发布

- **阶段进度跟踪**：各开发阶段的进度记录
- **发布说明**：版本特性、改进和修复的说明
- **变更日志**：详细的代码变更记录
