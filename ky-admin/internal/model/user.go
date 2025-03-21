package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password    string         `json:"-" gorm:"size:255;not null;comment:密码"`
	Nickname    string         `json:"nickname" gorm:"size:50;comment:昵称"`
	Email       string         `json:"email" gorm:"size:100;comment:邮箱"`
	Phone       string         `json:"phone" gorm:"size:20;comment:电话"`
	Avatar      string         `json:"avatar" gorm:"size:255;comment:头像"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 0:禁用 1:启用"` // 使用int保持一致
	TenantID    uint           `json:"tenant_id" gorm:"default:1;comment:租户ID"`      // 添加租户ID字段
	DeptID      uint           `json:"dept_id" gorm:"comment:部门ID"`                  // 添加部门ID字段
	LastLoginAt *time.Time     `json:"last_login_at" gorm:"comment:最后登录时间"`          // 添加最后登录时间
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (User) TableName() string {
	return "sys_user"
}

// LoginDTO 登录请求数据传输对象
type LoginDTO struct {
	Username string `json:"username" binding:"required" example:"admin"`       // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
}

// RegisterDTO 注册请求数据传输对象
type RegisterDTO struct {
	Username string `json:"username" binding:"required" example:"admin"`       // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
	Nickname string `json:"nickname" example:"管理员"`                            // 昵称
	Email    string `json:"email" example:"admin@example.com"`                 // 邮箱
	Phone    string `json:"phone" example:"13800138000"`                       // 电话
	RoleIDs  []uint `json:"roleIds,omitempty"`                                 // 角色ID列表
}

// UpdateUserDTO 更新用户信息请求数据传输对象
type UpdateUserDTO struct {
	Nickname string `json:"nickname" example:"管理员"`            // 昵称
	Email    string `json:"email" example:"admin@example.com"` // 邮箱
	Phone    string `json:"phone" example:"13800138000"`       // 电话
	Avatar   string `json:"avatar"`                            // 头像
	Status   *int8  `json:"status" example:"1"`                // 状态
	RoleIDs  []uint `json:"roleIds,omitempty"`                 // 角色ID列表
}

// ChangePasswordDTO 修改密码请求数据传输对象
type ChangePasswordDTO struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"password123"` // 旧密码
	NewPassword string `json:"newPassword" binding:"required" example:"newpass123"`  // 新密码
}

// UserDTO 用户信息数据传输对象
type UserDTO struct {
	ID        uint      `json:"id"`              // 用户ID
	Username  string    `json:"username"`        // 用户名
	Nickname  string    `json:"nickname"`        // 昵称
	Email     string    `json:"email"`           // 邮箱
	Phone     string    `json:"phone"`           // 电话
	Avatar    string    `json:"avatar"`          // 头像
	Status    int8      `json:"status"`          // 状态
	TenantID  uint      `json:"tenantId"`        // 租户ID
	DeptID    uint      `json:"deptId"`          // 部门ID
	CreatedAt time.Time `json:"createdAt"`       // 创建时间
	UpdatedAt time.Time `json:"updatedAt"`       // 更新时间
	Roles     []RoleDTO `json:"roles,omitempty"` // 角色列表
}

// UserPageQuery 用户分页查询参数
type UserPageQuery struct {
	Username string `form:"username" json:"username"` // 用户名
	Nickname string `form:"nickname" json:"nickname"` // 昵称
	Status   *int8  `form:"status" json:"status"`     // 状态
	TenantID *uint  `form:"tenantId" json:"tenantId"` // 租户ID
	DeptID   *uint  `form:"deptId" json:"deptId"`     // 部门ID
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页条数
}

// AssignRoleDTO 分配角色请求数据传输对象
type AssignRoleDTO struct {
	UserID  uint   `json:"userId" binding:"required"`  // 用户ID
	RoleIDs []uint `json:"roleIds" binding:"required"` // 角色ID列表
}
