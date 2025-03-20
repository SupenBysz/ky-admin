# YAPI集成方案

## 1. 概述

本文档描述了KY-Admin项目与YAPI(一个高效、易用、功能强大的API管理平台)的集成方案。YAPI能够帮助开发、产品、测试人员提供更优雅的接口管理服务，通过集成YAPI，我们可以实现API文档的自动生成、管理和测试。

## 2. YAPI简介

YAPI是由YMFE团队开发的一款开源API管理工具，可以本地部署，具有以下特性：

- 基于JSON5和Mockjs定义接口返回数据结构和文档
- 扁平化权限设计，兼顾了多人协作和管理
- 类似Postman的接口调试功能
- 自动化测试，支持对Response断言
- MockServer功能，支持随机Mock和Mock期望
- 支持Postman、HAR、Swagger数据导入
- 免费开源，支持内网部署

官方地址：[https://github.com/YMFE/yapi](https://github.com/YMFE/yapi)

## 3. 集成方案

### 3.1 YAPI部署

YAPI可以通过两种方式部署：

#### 方式一：官方yapi-cli工具部署

```bash
# 安装yapi-cli
npm install -g yapi-cli --registry https://registry.npm.taobao.org

# 启动可视化部署程序
yapi server

# 按照指引完成配置和部署
# 部署完成后，按提示启动服务
node {安装目录}/vendors/server/app.js
```

#### 方式二：Docker部署（推荐）

```bash
# 使用docker-compose一键部署
git clone https://github.com/fjc0k/docker-YApi.git
cd docker-YApi
cp .env.example .env
# 编辑.env文件，配置管理员账号等信息
docker-compose up -d
```

### 3.2 Swagger集成YAPI

KY-Admin使用Swagger自动生成API文档，并通过自动化工具将Swagger文档同步到YAPI，实现流程如下：

1. 在KY-Admin后端代码中使用Swagger注解定义API
2. 使用swag工具生成Swagger文档
3. 使用自动化工具将Swagger文档导入YAPI

#### 3.2.1 KY-Admin中的Swagger配置

```go
// 在main.go中添加Swagger文档
import (
    _ "github.com/SupenBysz/ky-admin/docs"  // 导入生成的docs
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// @title KY-Admin API
// @version 1.0
// @description KY-Admin后台管理系统API
// @host localhost:8080
// @BasePath /api/v1

func main() {
    // ...
    
    // 注册Swagger路由
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // ...
}
```

#### 3.2.2 API控制器中的Swagger注解示例

```go
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body LoginRequest true "登录信息"
// @Success 200 {object} Response{data=LoginResponse} "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "未授权"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
    // 处理登录逻辑
}
```

#### 3.2.3 自动同步Swagger到YAPI

我们选择使用自动化脚本定期将Swagger文档同步到YAPI，实现步骤如下：

1. 创建同步脚本（sync-swagger-to-yapi.js）：

```javascript
// 安装依赖: npm install axios --save
const axios = require('axios');
const fs = require('fs');
const path = require('path');

// 配置信息
const config = {
  yapiBaseUrl: 'http://yapi.example.com',  // YAPI服务器地址
  yapiToken: 'your-token-here',           // YAPI项目的token
  swaggerFile: path.join(__dirname, '../docs/swagger.json') // Swagger文件路径
};

async function syncSwaggerToYapi() {
  try {
    // 读取Swagger文件
    const swaggerContent = fs.readFileSync(config.swaggerFile, 'utf8');
    const swaggerJson = JSON.parse(swaggerContent);
    
    // 调用YAPI的导入接口
    const response = await axios.post(
      `${config.yapiBaseUrl}/api/open/import_swagger`,
      {
        type: 'swagger',
        token: config.yapiToken,
        json: swaggerJson
      }
    );
    
    console.log('同步结果:', response.data);
  } catch (error) {
    console.error('同步失败:', error);
  }
}

syncSwaggerToYapi();
```

2. 添加到CI/CD流程中，确保每次构建后自动同步文档

### 3.3 前端调用YAPI Mock服务

YAPI提供的Mock服务可以帮助前端开发人员在后端API尚未完成时进行开发。前端可以通过以下方式调用Mock服务：

1. 基础URL配置：

```javascript
// 配置axios基础URL
const apiBaseUrl = process.env.NODE_ENV === 'development' 
  ? 'http://yapi.example.com/mock/project_id'  // 开发环境使用Mock
  : '/api';  // 生产环境使用真实API
```

2. 调用示例：

```javascript
// 调用用户登录API
axios.post(`${apiBaseUrl}/auth/login`, {
  username: 'admin',
  password: 'password'
}).then(response => {
  console.log('登录成功', response.data);
}).catch(error => {
  console.error('登录失败', error);
});
```

## 4. YAPI使用指南

### 4.1 添加项目

1. 登录YAPI平台
2. 点击"添加项目"
3. 填写项目信息（名称、描述等）
4. 设置项目权限和成员

### 4.2 导入API文档

1. 进入项目页面
2. 点击"数据管理"→"导入数据"
3. 选择"Swagger"，可以直接上传文件或粘贴Swagger JSON内容
4. 点击"导入"完成

### 4.3 API文档查看与管理

1. 接口列表：项目主页可查看所有接口
2. 接口详情：点击接口名称查看详细信息
3. 接口编辑：点击"编辑"按钮可修改接口信息
4. 接口测试：点击"运行"按钮可进行接口测试

### 4.4 Mock数据使用

1. Mock地址：`http://yapi-server/mock/{project_id}/{interface_path}`
2. 查看接口详情页面可获取完整Mock地址
3. 支持Mock期望设置，根据请求参数返回不同数据

### 4.5 团队协作

1. 项目成员管理：在项目设置中添加/删除成员
2. 权限控制：可设置只读、读写、管理员等权限
3. 接口评论：支持在接口页面添加评论，方便沟通

## 5. 最佳实践

1. **规范API设计**：遵循RESTful API设计规范，便于管理和使用
2. **完善接口文档**：使用详细的Swagger注解描述接口功能、参数和响应
3. **自动化同步**：将Swagger到YAPI的同步集成到CI/CD流程中
4. **前后端分离开发**：前端可以基于YAPI Mock数据进行开发，不必等待后端完成
5. **测试驱动开发**：利用YAPI的自动化测试功能验证API行为

## 6. 注意事项

1. YAPI部署建议使用Docker方式，简化维护
2. 确保YAPI服务器有足够的资源，特别是内存
3. 定期备份YAPI数据库
4. 避免在公网暴露YAPI服务，最好在内网使用
5. 妥善保管YAPI的token，避免泄露

## 7. 相关资源

- [YAPI GitHub仓库](https://github.com/YMFE/yapi)
- [YAPI使用文档](https://hellosean1025.github.io/yapi/)
- [Swagger官方文档](https://swagger.io/docs/)
- [Gin Swagger文档](https://github.com/swaggo/gin-swagger)
- [Docker-YAPI部署方案](https://github.com/fjc0k/docker-YApi)
