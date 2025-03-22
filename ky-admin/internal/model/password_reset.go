package model

import (
	"time"

	"gorm.io/gorm"
)

// PasswordReset 密码重置模型
type PasswordReset struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index;not null;comment:用户ID"`
	Email     string         `json:"email" gorm:"size:100;index;not null;comment:邮箱"`
	Token     string         `json:"token" gorm:"size:255;not null;comment:重置令牌"`
	Used      bool           `json:"used" gorm:"default:false;comment:是否已使用"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null;comment:过期时间"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (PasswordReset) TableName() string {
	return "sys_password_reset"
}
