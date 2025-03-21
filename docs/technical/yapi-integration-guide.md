# YAPI集成使用指南

## 1. 简介

本指南详细说明了如何在KY-Admin项目中使用YAPI进行API文档管理。YAPI是一个强大的API管理平台，提供了接口管理、MockAPI服务、自动化测试等功能。

## 2. 部署YAPI

我们提供了两种部署YAPI的方式，推荐使用Docker方式。

### 2.1 Docker部署（推荐）

项目中已经准备了Docker Compose配置文件，可以轻松部署YAPI：

```bash
# 进入脚本目录
cd scripts

# 启动YAPI服务
docker-compose -f docker-compose.yapi.yml up -d

# 查看服务状态
docker-compose -f docker-compose.yapi.yml ps
```

YAPI服务将在 <http://localhost:3000> 上启动，初始管理员账号:

- 用户名: <admin@kysion.com>
- 密码: ky-admin123

### 2.2 手动部署

如果需要手动部署，请参考YAPI官方文档：<https://hellosean1025.github.io/yapi/>

## 3. Swagger集成

KY-Admin使用Swagger自动生成API文档，我们提供了一键更新脚本：

```bash
# 进入项目根目录
cd /path/to/ky-admin

# 运行更新脚本
./scripts/swagger-update.sh
```

更新后，可以通过 <http://localhost:8080/swagger/index.html> 访问Swagger UI。

## 4. 同步到YAPI

更新Swagger文档后，可以手动或自动同步到YAPI平台：

### 4.1 手动同步

1. 登录YAPI平台
2. 进入你的项目
3. 点击"数据管理" -> "导入数据"
4. 选择"Swagger"，上传`docs/swagger/swagger.json`文件
5. 点击导入

### 4.2 自动同步

我们提供了自动同步脚本，需要先配置YAPI项目的token：

1. 在YAPI平台的项目设置中获取token
2. 编辑`scripts/sync-swagger-to-yapi.js`文件，更新`yapiToken`和`yapiBaseUrl`
3. 执行以下命令同步：

```bash
# 安装依赖
npm install axios --no-save

# 运行同步脚本
node scripts/sync-swagger-to-yapi.js
```

也可以通过环境变量设置配置：

```bash
YAPI_URL=http://your-yapi-server:3000 YAPI_TOKEN=your-project-token node scripts/sync-swagger-to-yapi.js
```

### 4.3 CI/CD集成

可以将同步过程集成到CI/CD流程中：

```yaml
# 示例GitHub Actions配置
jobs:
  sync-api-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.24
      - name: Setup Node
        uses: actions/setup-node@v2
        with:
          node-version: '18'
      - name: Update Swagger
        run: ./scripts/swagger-update.sh
      - name: Sync to YAPI
        run: |
          npm install axios --no-save
          YAPI_URL=${{ secrets.YAPI_URL }} YAPI_TOKEN=${{ secrets.YAPI_TOKEN }} SYNC_TO_YAPI=true node scripts/sync-swagger-to-yapi.js
```

## 5. 使用YAPI Mock服务

YAPI提供了Mock服务，可以在前端开发中使用：

1. 在YAPI项目页面找到接口
2. 点击接口，查看详情
3. 在右上角找到"Mock地址"
4. 复制地址，前端可直接使用该地址调用Mock服务

## 6. API文档规范

为确保生成高质量的API文档，请遵循以下规范：

1. 使用Swagger注释:
   - 每个接口必须有`@Summary`和`@Description`
   - 指定`@Tags`分组API
   - 详细描述请求参数和响应

2. 结构体注释:
   - 每个字段添加注释说明用途
   - 使用`binding`标签说明验证规则
   - 使用`example`标签提供示例值

3. 示例:

```go
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body LoginRequest true "登录信息"
// @Success 200 {object} Response{data=LoginResponse} "登录成功"
// @Failure 400 {object} Response "参数错误"
// @Router /auth/login [post]
func Login(c *gin.Context) {
    // ...
}
```

## 7. 故障排除

### 7.1 Swagger生成失败

- 检查Swagger注解格式是否正确
- 确保导入了正确的包
- 尝试使用`swag fmt`命令格式化注解

### 7.2 YAPI同步失败

- 确认YAPI服务是否正常运行
- 检查token是否正确
- 检查网络连接
- 查看YAPI日志排查问题

### 7.3 接口文档不全

- 确保所有API都添加了完整的Swagger注解
- 检查结构体定义是否完整
- 重新生成Swagger文档

## 8. 注意事项

1. YAPI中的修改不会自动同步回代码，建议始终以代码中的Swagger注解为准
2. 每次API变更后，应及时更新Swagger文档并同步到YAPI
3. 在公共环境部署YAPI时，务必加强安全措施，避免数据泄露
4. 定期备份YAPI数据

## 9. 参考资源

- [YAPI官方文档](https://hellosean1025.github.io/yapi/)
- [Swagger官方文档](https://swagger.io/docs/)
- [swaggo/swag GitHub](https://github.com/swaggo/swag)
- [gin-swagger GitHub](https://github.com/swaggo/gin-swagger)
