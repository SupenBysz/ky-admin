package model

import "time"

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	RoleID       uint      `json:"roleId" gorm:"index;not null;comment:角色ID"`
	PermissionID uint      `json:"permissionId" gorm:"index;not null;comment:权限ID"`
	TenantID     uint      `json:"tenantId" gorm:"comment:租户ID"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TableName 设置表名
func (RolePermission) TableName() string {
	return "sys_role_permission"
}
