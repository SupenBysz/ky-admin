package dto

import (
	"time"
)

// IDRequest ID请求参数
type IDRequest struct {
	ID uint `uri:"id" binding:"required" example:"1"` // ID参数
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
	RoleIDs  []uint `json:"role_ids,omitempty"`                                // 角色ID列表
}

// CreateUserDTO 用户创建请求数据传输对象
type CreateUserDTO struct {
	Username string `json:"username" binding:"required" example:"admin"`       // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
	Nickname string `json:"nickname" example:"管理员"`                            // 昵称
	Email    string `json:"email" example:"admin@example.com"`                 // 邮箱
	Phone    string `json:"phone" example:"13800138000"`                       // 电话
	Avatar   string `json:"avatar"`                                            // 头像
	Status   int    `json:"status" example:"1"`                                // 状态 (使用int保持一致)
	TenantID uint   `json:"tenant_id" example:"1"`                             // 租户ID
	DeptID   uint   `json:"dept_id" example:"1"`                               // 部门ID
	RoleIDs  []uint `json:"role_ids,omitempty"`                                // 角色ID列表
}

// UpdateUserDTO 更新用户信息请求数据传输对象
type UpdateUserDTO struct {
	Nickname string `json:"nickname" example:"管理员"`            // 昵称
	Email    string `json:"email" example:"admin@example.com"` // 邮箱
	Phone    string `json:"phone" example:"13800138000"`       // 电话
	Avatar   string `json:"avatar"`                            // 头像
	Status   *int   `json:"status" example:"1"`                // 状态 (使用int保持一致)
	TenantID *uint  `json:"tenant_id" example:"1"`             // 租户ID
	DeptID   *uint  `json:"dept_id" example:"1"`               // 部门ID
	RoleIDs  []uint `json:"role_ids,omitempty"`                // 角色ID列表
}

// ChangePasswordDTO 修改密码请求数据传输对象
type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required" example:"password123"` // 旧密码
	NewPassword string `json:"new_password" binding:"required" example:"newpass123"`  // 新密码
}

// UserDTO 用户信息响应数据传输对象
type UserDTO struct {
	ID          uint       `json:"id"`                    // 用户ID
	Username    string     `json:"username"`              // 用户名
	Nickname    string     `json:"nickname"`              // 昵称
	Email       string     `json:"email"`                 // 邮箱
	Phone       string     `json:"phone"`                 // 电话
	Avatar      string     `json:"avatar"`                // 头像
	Status      int8       `json:"status"`                // 状态
	TenantID    uint       `json:"tenant_id"`             // 租户ID
	DeptID      uint       `json:"dept_id"`               // 部门ID
	LastLoginAt *time.Time `json:"last_login_at"`         // 最后登录时间
	CreatedAt   time.Time  `json:"created_at"`            // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`            // 更新时间
	Roles       []RoleDTO  `json:"roles,omitempty"`       // 角色列表
	Permissions []string   `json:"permissions,omitempty"` // 权限列表
}

// UserPageQueryDTO 用户分页查询参数
type UserPageQueryDTO struct {
	Username string `form:"username" json:"username"`   // 用户名
	Nickname string `form:"nickname" json:"nickname"`   // 昵称
	Status   *int   `form:"status" json:"status"`       // 状态 (使用int保持一致)
	TenantID *uint  `form:"tenant_id" json:"tenant_id"` // 租户ID
	DeptID   *uint  `form:"dept_id" json:"dept_id"`     // 部门ID
	Page     int    `form:"page" json:"page"`           // 页码
	PageSize int    `form:"page_size" json:"page_size"` // 每页条数
}

// UserPageResponseDTO 用户分页查询响应
type UserPageResponseDTO struct {
	Total int64     `json:"total"` // 总记录数
	List  []UserDTO `json:"list"`  // 用户列表
	Page  int       `json:"page"`  // 当前页码
	Size  int       `json:"size"`  // 每页大小
}

// AssignRoleDTO 分配角色请求数据传输对象
type AssignRoleDTO struct {
	UserID  uint   `json:"user_id" binding:"required"`  // 用户ID
	RoleIDs []uint `json:"role_ids" binding:"required"` // 角色ID列表
}

// 临时定义RoleDTO，后续将完善
type RoleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ForgotPasswordDTO 忘记密码DTO
type ForgotPasswordDTO struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

// ResetPasswordDTO 重置密码DTO
type ResetPasswordDTO struct {
	Token       string `json:"token" form:"token" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required,min=6"`
}

// UserResponseDTO 用户响应数据传输对象
type UserResponseDTO struct {
	ID          uint       `json:"id"`                    // 用户ID
	Username    string     `json:"username"`              // 用户名
	Nickname    string     `json:"nickname"`              // 昵称
	Email       string     `json:"email"`                 // 邮箱
	Phone       string     `json:"phone"`                 // 电话
	Avatar      string     `json:"avatar"`                // 头像
	Status      int8       `json:"status"`                // 状态
	TenantID    uint       `json:"tenant_id"`             // 租户ID
	DeptID      uint       `json:"dept_id"`               // 部门ID
	LastLoginAt *time.Time `json:"last_login_at"`         // 最后登录时间
	CreatedAt   time.Time  `json:"created_at"`            // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`            // 更新时间
	Roles       []RoleDTO  `json:"roles,omitempty"`       // 角色列表
	Permissions []string   `json:"permissions,omitempty"` // 权限列表
}
