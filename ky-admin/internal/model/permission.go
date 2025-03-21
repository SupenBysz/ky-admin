package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"type:varchar(50);not null;comment:权限名称"`
	Code        string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:权限编码"`
	Description string         `json:"description" gorm:"type:varchar(255);comment:权限描述"`
	Type        int8           `json:"type" gorm:"type:tinyint(1);default:2;comment:权限类型 1:菜单 2:操作"`
	Module      string         `json:"module" gorm:"type:varchar(50);comment:所属模块"`
	TenantID    uint           `json:"tenantId" gorm:"comment:租户ID"`
	ParentID    uint           `json:"parentId" gorm:"default:0;comment:父级权限ID"`
	Path        string         `json:"path" gorm:"type:varchar(100);comment:路径"`
	Order       int            `json:"order" gorm:"type:int;default:0;comment:排序"`
	Icon        string         `json:"icon" gorm:"type:varchar(100);comment:图标"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Roles       []*Role        `json:"roles,omitempty" gorm:"many2many:sys_role_permission;"`
}

// TableName 设置表名
func (Permission) TableName() string {
	return "sys_permission"
}

// PermissionDTO 权限数据传输对象
type PermissionDTO struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Code        string           `json:"code"`
	Description string           `json:"description"`
	Type        int8             `json:"type"`
	Module      string           `json:"module"`
	TenantID    uint             `json:"tenantId"`
	ParentID    uint             `json:"parentId"`
	Path        string           `json:"path"`
	Order       int              `json:"order"`
	Icon        string           `json:"icon"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	Children    []*PermissionDTO `json:"children,omitempty"`
}

// CreatePermissionDTO 创建权限请求数据传输对象
type CreatePermissionDTO struct {
	Name        string `json:"name" binding:"required" example:"用户查询"`       // 权限名称
	Code        string `json:"code" binding:"required" example:"user:query"` // 权限编码
	Description string `json:"description" example:"查询用户列表"`                 // 权限描述
	Type        int8   `json:"type" example:"2"`                             // 权限类型
	Module      string `json:"module" example:"user"`                        // 所属模块
	ParentID    uint   `json:"parentId" example:"0"`                         // 父级权限ID
	Path        string `json:"path" example:"/user/list"`                    // 路径
	Order       int    `json:"order" example:"0"`                            // 排序
	Icon        string `json:"icon" example:"user"`                          // 图标
}

// UpdatePermissionDTO 更新权限请求数据传输对象
type UpdatePermissionDTO struct {
	Name        string `json:"name" example:"用户查询"`          // 权限名称
	Description string `json:"description" example:"查询用户列表"` // 权限描述
	Type        *int8  `json:"type" example:"2"`             // 权限类型
	Module      string `json:"module" example:"user"`        // 所属模块
	ParentID    *uint  `json:"parentId" example:"0"`         // 父级权限ID
	Path        string `json:"path" example:"/user/list"`    // 路径
	Order       *int   `json:"order" example:"0"`            // 排序
	Icon        string `json:"icon" example:"user"`          // 图标
}

// PermissionPageQuery 权限分页查询参数
type PermissionPageQuery struct {
	Name     string `form:"name" json:"name"`         // 权限名称
	Code     string `form:"code" json:"code"`         // 权限编码
	Type     *int8  `form:"type" json:"type"`         // 权限类型
	Module   string `form:"module" json:"module"`     // 所属模块
	TenantID *uint  `form:"tenantId" json:"tenantId"` // 租户ID
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页条数
}
