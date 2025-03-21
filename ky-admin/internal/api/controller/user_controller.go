package controller

import (
	"github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/dto"
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

// GetUserInfo 获取当前用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		api.Error(ctx, errors.NewAuthError(errors.ErrInvalidToken))
		return
	}

	// 获取用户信息
	userInfo, err := c.userService.GetUserWithRoles(userID.(uint))
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, userInfo)
}

// UpdateUserInfo 更新当前用户信息
func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		api.Error(ctx, errors.NewAuthError(errors.ErrInvalidToken))
		return
	}

	// 解析请求参数
	var updateDTO dto.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 获取用户信息
	user, err := c.userService.GetByID(userID.(uint))
	if err != nil {
		api.Error(ctx, err)
		return
	}

	// 更新用户信息
	if updateDTO.Nickname != "" {
		user.Nickname = updateDTO.Nickname
	}
	if updateDTO.Email != "" {
		user.Email = updateDTO.Email
	}
	if updateDTO.Phone != "" {
		user.Phone = updateDTO.Phone
	}
	if updateDTO.Avatar != "" {
		user.Avatar = updateDTO.Avatar
	}

	// 保存更新
	if err := c.userService.UpdateUser(user); err != nil {
		api.Error(ctx, err)
		return
	}

	// 获取更新后的用户信息
	userInfo, err := c.userService.GetUserWithRoles(userID.(uint))
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, userInfo)
}

// ChangePassword 修改密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		api.Error(ctx, errors.NewAuthError(errors.ErrInvalidToken))
		return
	}

	// 解析请求参数
	var changeDTO dto.ChangePasswordDTO
	if err := ctx.ShouldBindJSON(&changeDTO); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 修改密码
	if err := c.userService.ChangePassword(userID.(uint), changeDTO.OldPassword, changeDTO.NewPassword); err != nil {
		api.Error(ctx, err)
		return
	}

	api.SuccessWithMessage(ctx, "密码修改成功", nil)
}

// GetUserList 获取用户列表
func (c *UserController) GetUserList(ctx *gin.Context) {
	// 解析分页参数
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

	// 获取用户列表
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

	api.SuccessWithPage(ctx, userDTOs, total, query.Page, query.PageSize)
}

// GetUserByID 根据ID获取用户信息
func (c *UserController) GetUserByID(ctx *gin.Context) {
	// 解析用户ID
	var req dto.IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 获取用户信息
	user, err := c.userService.GetUserWithRoles(req.ID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, user)
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	// 解析请求参数
	var createDTO dto.CreateUserDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 创建用户
	userID, err := c.userService.CreateUserWithRoles(createDTO)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, gin.H{"id": userID})
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// 解析用户ID
	var req dto.IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 解析请求参数
	var updateDTO dto.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 获取用户信息
	user, err := c.userService.GetByID(req.ID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	// 更新用户信息
	if updateDTO.Nickname != "" {
		user.Nickname = updateDTO.Nickname
	}
	if updateDTO.Email != "" {
		user.Email = updateDTO.Email
	}
	if updateDTO.Phone != "" {
		user.Phone = updateDTO.Phone
	}
	if updateDTO.Avatar != "" {
		user.Avatar = updateDTO.Avatar
	}

	// 保存更新
	if err := c.userService.UpdateUser(user); err != nil {
		api.Error(ctx, err)
		return
	}

	api.SuccessWithMessage(ctx, "更新成功", nil)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// 解析用户ID
	var req dto.IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 删除用户
	if err := c.userService.DeleteUser(req.ID); err != nil {
		api.Error(ctx, err)
		return
	}

	api.SuccessWithMessage(ctx, "删除成功", nil)
}

// AssignUserRoles 分配用户角色
func (c *UserController) AssignUserRoles(ctx *gin.Context) {
	// 解析请求参数
	var assignDTO dto.AssignRoleDTO
	if err := ctx.ShouldBindJSON(&assignDTO); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 分配用户角色
	if err := c.userService.AssignRoles(assignDTO.UserID, assignDTO.RoleIDs); err != nil {
		api.Error(ctx, err)
		return
	}

	api.SuccessWithMessage(ctx, "角色分配成功", nil)
}
