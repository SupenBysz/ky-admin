package dto

import (
	"time"
)

// RoleInfo 角色信息
type RoleInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// RoleDetailRes 角色详情响应
type RoleDetailRes struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Desc      string    `json:"desc"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RoleListReq 角色列表查询请求
type RoleListReq struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100"`
	Name     string `form:"name"`
	Status   *uint  `form:"status"`
}

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Desc   string `json:"desc"`
	Status uint   `json:"status" binding:"required"`
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Status *uint  `json:"status,omitempty"`
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Object string `json:"object"`
	Action string `json:"action"`
	Desc   string `json:"desc"`
}

// AssignPermissionReq 分配权限请求
type AssignPermissionReq struct {
	RoleID        uint   `json:"roleId" binding:"required"`
	PermissionIDs []uint `json:"permissionIds" binding:"required"`
}

// RoleCreateDTO 角色创建请求数据传输对象
type RoleCreateDTO struct {
	Name        string `json:"name" binding:"required" example:"管理员"`   // 角色名称
	Code        string `json:"code" binding:"required" example:"admin"` // 角色编码
	Description string `json:"description" example:"系统管理员角色"`           // 角色描述
	Status      int    `json:"status" example:"1"`                      // 状态 (使用int保持一致)
	Sort        int    `json:"sort" example:"1"`                        // 排序
	PermIDs     []uint `json:"perm_ids,omitempty"`                      // 权限ID列表
}

// RoleUpdateDTO 角色更新请求数据传输对象
type RoleUpdateDTO struct {
	Name        string `json:"name" example:"管理员"`            // 角色名称
	Description string `json:"description" example:"系统管理员角色"` // 角色描述
	Status      *int   `json:"status" example:"1"`            // 状态 (使用int保持一致)
	Sort        *int   `json:"sort" example:"1"`              // 排序
	PermIDs     []uint `json:"perm_ids,omitempty"`            // 权限ID列表
}

// RoleResponseDTO 角色响应数据传输对象
type RoleResponseDTO struct {
	ID          uint                    `json:"id"`                    // 角色ID
	Name        string                  `json:"name"`                  // 角色名称
	Code        string                  `json:"code"`                  // 角色编码
	Description string                  `json:"description"`           // 角色描述
	Status      int                     `json:"status"`                // 状态 (使用int保持一致)
	Sort        int                     `json:"sort"`                  // 排序
	CreatedAt   time.Time               `json:"created_at"`            // 创建时间
	UpdatedAt   time.Time               `json:"updated_at"`            // 更新时间
	Permissions []PermissionResponseDTO `json:"permissions,omitempty"` // 权限列表
}

// RolePageQueryDTO 角色分页查询参数
type RolePageQueryDTO struct {
	Name   string `form:"name" json:"name"`           // 角色名称
	Code   string `form:"code" json:"code"`           // 角色编码
	Status *int   `form:"status" json:"status"`       // 状态 (使用int保持一致)
	Page   int    `form:"page" json:"page"`           // 页码
	Size   int    `form:"page_size" json:"page_size"` // 每页条数
}

// RolePageResponseDTO 角色分页查询响应
type RolePageResponseDTO struct {
	Total int64             `json:"total"` // 总记录数
	List  []RoleResponseDTO `json:"list"`  // 角色列表
	Page  int               `json:"page"`  // 当前页码
	Size  int               `json:"size"`  // 每页大小
}

// AssignPermissionDTO 分配权限请求数据传输对象
type AssignPermissionDTO struct {
	RoleID  uint   `json:"role_id" binding:"required"`  // 角色ID
	PermIDs []uint `json:"perm_ids" binding:"required"` // 权限ID列表
}
