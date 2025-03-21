package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index;comment:用户ID"`
	RoleID    uint           `json:"role_id" gorm:"not null;index;comment:角色ID"`
	TenantID  uint           `json:"tenant_id" gorm:"not null;index;comment:租户ID;default:1"` // 租户ID
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (UserRole) TableName() string {
	return "sys_user_role"
}
