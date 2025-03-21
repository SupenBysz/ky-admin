package router

import (
	"github.com/SupenBysz/ky-admin/internal/api/controller"
	"github.com/SupenBysz/ky-admin/internal/api/middleware"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(
	r *gin.Engine,
	userService service.UserService,
	jwtService service.JWTService,
	permissionService service.PermissionService,
	enforcer *casbin.Enforcer,
) {
	// 创建控制器
	authController := controller.NewAuthController(userService, jwtService)
	userController := controller.NewUserController(userService)
	roleController := controller.NewRoleController(permissionService)

	// JWT中间件
	authMiddleware := middleware.JWTAuthMiddleware(jwtService)

	// 权限中间件（暂未实现）
	// authorizeMiddleware := middleware.Authorize(enforcer, userService)

	// API组
	api := r.Group("/api")
	{
		// 认证路由 - 不需要认证
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
			auth.POST("/refresh", authController.RefreshToken)
		}

		// 用户路由 - 需要认证
		user := api.Group("/user").Use(authMiddleware)
		{
			user.GET("", userController.GetUserInfo)             // 获取当前用户信息
			user.PUT("", userController.UpdateUserInfo)          // 更新当前用户信息
			user.PUT("/password", userController.ChangePassword) // 修改密码
		}

		// 用户管理路由 - 需要认证
		users := api.Group("/users").Use(authMiddleware)
		{
			users.GET("", userController.GetUserList)           // 获取用户列表
			users.GET("/:id", userController.GetUserByID)       // 获取指定用户信息
			users.POST("", userController.CreateUser)           // 创建用户
			users.PUT("/:id", userController.UpdateUser)        // 更新用户
			users.DELETE("/:id", userController.DeleteUser)     // 删除用户
			users.POST("/role", userController.AssignUserRoles) // 分配用户角色
		}

		// 角色管理路由 - 需要认证
		roles := api.Group("/roles").Use(authMiddleware)
		{
			roles.GET("", roleController.GetRoleList)                    // 获取角色列表
			roles.GET("/:id", roleController.GetRoleByID)                // 获取角色详情
			roles.POST("", roleController.CreateRole)                    // 创建角色
			roles.PUT("/:id", roleController.UpdateRole)                 // 更新角色
			roles.DELETE("/:id", roleController.DeleteRole)              // 删除角色
			roles.GET("/permissions", roleController.GetAllPermissions)  // 获取所有权限
			roles.POST("/permissions", roleController.AssignPermissions) // 分配权限
		}
	}
}
