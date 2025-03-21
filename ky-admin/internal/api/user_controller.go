package api

import (
	"github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/middleware"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	user, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	// 获取用户信息并返回
	userDTO, err := c.userService.GetUserWithRoles(user.ID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, userDTO)
}

// GetUserInfo 获取当前用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID := middleware.CurrentUser(ctx)
	if userID == 0 {
		api.Error(ctx, errors.NewAuthError(errors.ErrInvalidToken))
		return
	}

	// 获取用户信息
	userDTO, err := c.userService.GetUserWithRoles(userID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, userDTO)
}

// GetUserByID 通过ID获取用户
func (c *UserController) GetUserByID(ctx *gin.Context) {
	var req dto.IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 获取用户信息
	userDTO, err := c.userService.GetUserWithRoles(req.ID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, userDTO)
}

// GetUsers 获取用户列表
func (c *UserController) GetUsers(ctx *gin.Context) {
	var query dto.UserPageQueryDTO
	if err := ctx.ShouldBindQuery(&query); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的查询参数", err))
		return
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	// 查询用户列表
	users, total, err := c.userService.GetUsers(query.Page, query.PageSize)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	// 转换为DTO
	userDTOs := make([]dto.UserDTO, 0, len(users))
	for _, user := range users {
		// 查询用户角色
		roles, _ := c.userService.GetUserRoles(user.ID)

		roleDTOs := make([]dto.RoleDTO, 0, len(roles))
		for _, role := range roles {
			roleDTOs = append(roleDTOs, dto.RoleDTO{
				ID:   role.ID,
				Name: role.Name,
				Code: role.Code,
			})
		}

		userDTOs = append(userDTOs, dto.UserDTO{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			Status:    int8(user.Status),
			TenantID:  user.TenantID,
			DeptID:    user.DeptID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Roles:     roleDTOs,
		})
	}

	// 调整参数顺序以适应SuccessWithPage
	api.SuccessWithPage(ctx, userDTOs, total, query.Page, query.PageSize)
}

// UpdateUserInfo 更新用户信息
func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	var req dto.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 获取当前用户ID
	userID := middleware.CurrentUser(ctx)
	if userID == 0 {
		api.Error(ctx, errors.NewAuthError(errors.ErrInvalidToken))
		return
	}

	// 获取用户信息
	user, err := c.userService.GetByID(userID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// 保存用户信息
	if err := c.userService.UpdateUser(user); err != nil {
		api.Error(ctx, err)
		return
	}

	// 返回更新后的用户信息
	userDTO, err := c.userService.GetUserWithRoles(userID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, userDTO)
}

// ChangePassword 修改密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req dto.ChangePasswordDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 获取当前用户ID
	userID := middleware.CurrentUser(ctx)
	if userID == 0 {
		api.Error(ctx, errors.NewAuthError(errors.ErrInvalidToken))
		return
	}

	// 调用服务修改密码
	if err := c.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		api.Error(ctx, err)
		return
	}

	api.SuccessWithMessage(ctx, "密码修改成功", nil)
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 调用服务创建用户
	userID, err := c.userService.CreateUserWithRoles(req)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, gin.H{"id": userID})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	var req dto.IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 调用服务删除用户
	if err := c.userService.DeleteUser(req.ID); err != nil {
		api.Error(ctx, err)
		return
	}

	api.SuccessWithMessage(ctx, "删除成功", nil)
}

// AssignRoles 分配角色
func (c *UserController) AssignRoles(ctx *gin.Context) {
	var req dto.AssignRoleDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 调用服务分配角色
	if err := c.userService.AssignRoles(req.UserID, req.RoleIDs); err != nil {
		api.Error(ctx, err)
		return
	}

	api.SuccessWithMessage(ctx, "角色分配成功", nil)
}
