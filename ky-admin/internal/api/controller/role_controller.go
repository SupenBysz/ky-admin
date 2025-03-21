package controller

import (
	"strconv"

	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/pkg/response"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	roleService service.PermissionService
}

// NewRoleController 创建角色控制器
func NewRoleController(roleService service.PermissionService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
}

// GetRoleList 获取角色列表
func (c *RoleController) GetRoleList(ctx *gin.Context) {
	// 获取所有角色
	roles, err := c.roleService.GetAllRoles()
	if err != nil {
		response.ServerError(ctx, "获取角色列表失败："+err.Error())
		return
	}

	// 转换为角色DTO
	var roleList []*dto.RoleInfo
	for i, roleName := range roles {
		roleList = append(roleList, &dto.RoleInfo{
			ID:   uint(i + 1), // 这里简化处理，实际应用需要从数据库获取真实ID
			Name: roleName,
			Code: roleName,
		})
	}

	response.Success(ctx, roleList)
}

// GetRoleByID 根据ID获取角色
func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	// 获取角色ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的角色ID")
		return
	}

	// 这里简化处理，实际应用需要从数据库获取角色详情
	response.Success(ctx, &dto.RoleDetailRes{
		ID:   uint(id),
		Name: "角色" + idStr,
		Code: "role_" + idStr,
		Desc: "角色" + idStr + "的描述",
	})
}

// CreateRole 创建角色
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req dto.CreateRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误："+err.Error())
		return
	}

	// 创建角色，添加默认权限
	err := c.roleService.AddPermissionForRole(req.Code, "*", "*")
	if err != nil {
		response.ServerError(ctx, "创建角色失败："+err.Error())
		return
	}

	response.Success(ctx, gin.H{"role": req.Code})
}

// UpdateRole 更新角色
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	// 获取角色ID
	idStr := ctx.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的角色ID")
		return
	}

	var req dto.UpdateRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误："+err.Error())
		return
	}

	// 更新角色
	// 在实际应用中，需要实现角色更新逻辑

	response.Success(ctx, gin.H{"updated": true})
}

// DeleteRole 删除角色
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	// 获取角色ID
	idStr := ctx.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的角色ID")
		return
	}

	// 删除角色
	// 在实际应用中，应检查角色是否存在，以及是否有用户关联

	response.Success(ctx, gin.H{"deleted": true})
}

// GetAllPermissions 获取所有权限
func (c *RoleController) GetAllPermissions(ctx *gin.Context) {
	// 在实际应用中，应从数据库或配置中获取所有可用权限
	var permissions []*dto.PermissionInfo = []*dto.PermissionInfo{
		{ID: 1, Name: "查看用户列表", Object: "/api/users", Action: "GET", Desc: "允许查看所有用户"},
		{ID: 2, Name: "查看用户详情", Object: "/api/users/:id", Action: "GET", Desc: "允许查看用户详细信息"},
		{ID: 3, Name: "创建用户", Object: "/api/users", Action: "POST", Desc: "允许创建新用户"},
		{ID: 4, Name: "更新用户", Object: "/api/users/:id", Action: "PUT", Desc: "允许更新用户信息"},
		{ID: 5, Name: "删除用户", Object: "/api/users/:id", Action: "DELETE", Desc: "允许删除用户"},
		{ID: 6, Name: "查看角色列表", Object: "/api/roles", Action: "GET", Desc: "允许查看所有角色"},
		{ID: 7, Name: "查看角色详情", Object: "/api/roles/:id", Action: "GET", Desc: "允许查看角色详细信息"},
		{ID: 8, Name: "创建角色", Object: "/api/roles", Action: "POST", Desc: "允许创建新角色"},
		{ID: 9, Name: "更新角色", Object: "/api/roles/:id", Action: "PUT", Desc: "允许更新角色信息"},
		{ID: 10, Name: "删除角色", Object: "/api/roles/:id", Action: "DELETE", Desc: "允许删除角色"},
	}

	response.Success(ctx, permissions)
}

// AssignPermissions 分配权限
func (c *RoleController) AssignPermissions(ctx *gin.Context) {
	var req dto.AssignPermissionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误："+err.Error())
		return
	}

	// 为角色添加权限
	// 在实际应用中，需要先获取权限详情，再为角色分配权限
	// 这里简化处理，假设req.RoleID能够映射到角色名称，权限ID能够映射到资源和操作

	// 模拟成功
	response.Success(ctx, gin.H{"assigned": true})
}
