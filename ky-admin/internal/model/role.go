package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"type:varchar(50);not null;comment:角色名称"`
	Code        string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:角色编码"`
	Description string         `json:"description" gorm:"type:varchar(255);comment:角色描述"`
	Sort        int            `json:"sort" gorm:"type:int;default:0;comment:排序"`
	Status      int8           `json:"status" gorm:"type:tinyint(1);default:1;comment:状态 0:禁用 1:启用"`
	TenantID    uint           `json:"tenantId" gorm:"comment:租户ID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Permissions []*Permission  `json:"permissions,omitempty" gorm:"many2many:sys_role_permission;"`
}

// TableName 设置表名
func (Role) TableName() string {
	return "sys_role"
}

// RoleDTO 角色数据传输对象
type RoleDTO struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Description   string    `json:"description"`
	Sort          int       `json:"sort"`
	Status        int8      `json:"status"`
	TenantID      uint      `json:"tenantId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PermissionIDs []uint    `json:"permissionIds,omitempty"`
}

// CreateRoleDTO 创建角色请求数据传输对象
type CreateRoleDTO struct {
	Name          string `json:"name" binding:"required" example:"管理员"`   // 角色名称
	Code          string `json:"code" binding:"required" example:"admin"` // 角色编码
	Description   string `json:"description" example:"系统管理员角色"`           // 角色描述
	Sort          int    `json:"sort" example:"0"`                        // 排序
	Status        int8   `json:"status" example:"1"`                      // 状态
	PermissionIDs []uint `json:"permissionIds,omitempty"`                 // 权限ID列表
}

// UpdateRoleDTO 更新角色请求数据传输对象
type UpdateRoleDTO struct {
	Name          string `json:"name" example:"管理员"`            // 角色名称
	Description   string `json:"description" example:"系统管理员角色"` // 角色描述
	Sort          *int   `json:"sort" example:"0"`              // 排序
	Status        *int8  `json:"status" example:"1"`            // 状态
	PermissionIDs []uint `json:"permissionIds,omitempty"`       // 权限ID列表
}

// RolePageQuery 角色分页查询参数
type RolePageQuery struct {
	Name     string `form:"name" json:"name"`         // 角色名称
	Code     string `form:"code" json:"code"`         // 角色编码
	Status   *int8  `form:"status" json:"status"`     // 状态
	TenantID *uint  `form:"tenantId" json:"tenantId"` // 租户ID
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页条数
}
