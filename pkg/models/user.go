package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	StatusModel
	Username     string    `gorm:"size:50;not null;uniqueIndex;comment:用户名" json:"username"`       // 用户名
	Password     string    `gorm:"size:100;not null;comment:密码" json:"-"`                          // 密码（加密存储）
	Nickname     string    `gorm:"size:50;comment:昵称" json:"nickname"`                             // 昵称
	Email        string    `gorm:"size:100;uniqueIndex;comment:邮箱" json:"email"`                   // 邮箱
	Mobile       string    `gorm:"size:20;uniqueIndex;comment:手机号" json:"mobile"`                  // 手机号
	Avatar       string    `gorm:"size:255;comment:头像" json:"avatar"`                              // 头像URL
	Gender       uint8     `gorm:"default:0;comment:性别 0:未知 1:男 2:女" json:"gender"`                // 性别
	Remark       string    `gorm:"size:200;comment:备注" json:"remark"`                              // 备注
	LastLoginAt  time.Time `gorm:"comment:最后登录时间" json:"last_login_at"`                            // 最后登录时间
	LastLoginIP  string    `gorm:"size:50;comment:最后登录IP" json:"last_login_ip"`                    // 最后登录IP
	PasswordLock bool      `gorm:"default:false;comment:密码锁定" json:"password_lock"`                // 密码锁定（禁止修改密码）
	Roles        []*Role   `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;" json:"roles"` // 角色列表
	IsAdmin      bool      `gorm:"default:false;comment:是否管理员" json:"is_admin"`                    // 是否是管理员
}

// TableName 指定表名
func (User) TableName() string {
	return "ky_user"
}

// BeforeSave 保存前加密密码
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 只有当密码有变更时才加密
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// ValidatePassword 验证密码
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ChangePassword 修改密码
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	// 检查密码锁定状态
	if u.PasswordLock {
		return errors.New("密码已锁定，无法修改")
	}

	// 验证旧密码
	if !u.ValidatePassword(oldPassword) {
		return errors.New("旧密码验证失败")
	}

	// 设置新密码
	u.Password = newPassword
	return DB.Save(u).Error
}

// UpdateLoginInfo 更新登录信息
func (u *User) UpdateLoginInfo(ip string) error {
	u.LastLoginAt = time.Now()
	u.LastLoginIP = ip
	return DB.Model(u).Updates(map[string]interface{}{
		"last_login_at": u.LastLoginAt,
		"last_login_ip": u.LastLoginIP,
	}).Error
}

// Enable 启用用户
func (u *User) Enable() error {
	if u.Status == uint8(StatusEnabled) {
		return nil
	}
	u.Status = uint8(StatusEnabled)
	return DB.Model(u).Update("status", uint8(StatusEnabled)).Error
}

// Disable 禁用用户
func (u *User) Disable() error {
	if u.Status == uint8(StatusDisabled) {
		return nil
	}
	u.Status = uint8(StatusDisabled)
	return DB.Model(u).Update("status", uint8(StatusDisabled)).Error
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleID uint) bool {
	// 需要预加载角色数据才能使用此方法
	for _, role := range u.Roles {
		if role.ID == roleID {
			return true
		}
	}
	return false
}

// AddRole 为用户添加角色
func (u *User) AddRole(roleID uint) error {
	return DB.Model(u).Association("Roles").Append(&Role{StatusModel: StatusModel{BaseModel: BaseModel{ID: roleID}}})
}

// AssignRole 为用户分配角色（事务中使用）
func (u *User) AssignRole(tx *gorm.DB, roleID uint) error {
	// 检查角色是否存在
	var count int64
	if err := tx.Model(&Role{}).Where("id = ?", roleID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return errors.New("角色不存在")
	}

	// 创建角色关联
	role := &Role{
		StatusModel: StatusModel{
			BaseModel: BaseModel{ID: roleID},
		},
	}

	// 添加角色关联
	return tx.Model(u).Association("Roles").Append(role)
}

// RemoveRole 移除用户角色
func (u *User) RemoveRole(roleID uint) error {
	return DB.Model(u).Association("Roles").Delete(&Role{StatusModel: StatusModel{BaseModel: BaseModel{ID: roleID}}})
}

// ClearRoles 清空用户角色
func (u *User) ClearRoles() error {
	return DB.Model(u).Association("Roles").Clear()
}

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
