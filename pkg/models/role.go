package models

import (
	"errors"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	StatusModel
	Name        string        `gorm:"type:varchar(50);not null;uniqueIndex:idx_role_name;comment:角色名称"`
	Code        string        `gorm:"type:varchar(50);not null;uniqueIndex:idx_role_code;comment:角色编码"`
	Description string        `gorm:"type:varchar(200);comment:角色描述"`
	Sort        int           `gorm:"default:0;comment:排序"`
	IsSystem    bool          `gorm:"default:false;comment:是否系统角色"`
	Permissions []*Permission `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Users       []*User       `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 表名
func (Role) TableName() string {
	return "roles"
}

// BeforeDelete 删除前钩子
func (r *Role) BeforeDelete(tx *gorm.DB) error {
	// 系统角色不允许删除
	if r.IsSystem {
		return ErrSystemRoleNotDeletable
	}
	return nil
}

// Enable 启用角色
func (r *Role) Enable() error {
	if r.Status == uint8(StatusEnabled) {
		return nil
	}
	r.Status = uint8(StatusEnabled)
	return DB.Model(r).Update("status", uint8(StatusEnabled)).Error
}

// Disable 禁用角色
func (r *Role) Disable() error {
	if r.Status == uint8(StatusDisabled) {
		return nil
	}
	r.Status = uint8(StatusDisabled)
	return DB.Model(r).Update("status", uint8(StatusDisabled)).Error
}

// IsEnabled 是否启用
func (r *Role) IsEnabled() bool {
	return r.Status == uint8(StatusEnabled)
}

// AddPermission 为角色添加权限
func (r *Role) AddPermission(permissionID uint) error {
	if DB == nil {
		return ErrDBNotRegistered
	}
	return DB.Model(r).Association("Permissions").Append(&Permission{StatusModel: StatusModel{BaseModel: BaseModel{ID: permissionID}}})
}

// RemovePermission 移除角色权限
func (r *Role) RemovePermission(permissionID uint) error {
	if DB == nil {
		return ErrDBNotRegistered
	}
	return DB.Model(r).Association("Permissions").Delete(&Permission{StatusModel: StatusModel{BaseModel: BaseModel{ID: permissionID}}})
}

// ClearPermissions 清空角色权限
func (r *Role) ClearPermissions() error {
	if DB == nil {
		return ErrDBNotRegistered
	}
	return DB.Model(r).Association("Permissions").Clear()
}

// HasPermission 是否拥有权限
func (r *Role) HasPermission(permissionID uint) bool {
	for _, permission := range r.Permissions {
		if permission.ID == permissionID {
			return true
		}
	}
	return false
}

// HasPermissionByCode 是否拥有指定代码的权限
func (r *Role) HasPermissionByCode(code string) bool {
	for _, permission := range r.Permissions {
		if permission.Code == code {
			return true
		}
	}
	return false
}

// AssignPermissions 为角色分配权限
func (r *Role) AssignPermissions(tx *gorm.DB, permissionIDs []uint) error {
	// 清空原有权限
	if err := tx.Model(r).Association("Permissions").Clear(); err != nil {
		return err
	}

	// 没有新权限要分配，直接返回
	if len(permissionIDs) == 0 {
		return nil
	}

	// 创建权限切片
	permissions := make([]*Permission, 0, len(permissionIDs))
	for _, id := range permissionIDs {
		permissions = append(permissions, &Permission{
			StatusModel: StatusModel{
				BaseModel: BaseModel{ID: id},
			},
		})
	}

	// 分配新权限
	return tx.Model(r).Association("Permissions").Append(permissions)
}

// GetRoleByID 根据ID获取角色
func GetRoleByID(id uint) (*Role, error) {
	var role Role
	if err := DB.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

// GetRoleByCode 根据代码获取角色
func GetRoleByCode(code string) (*Role, error) {
	var role Role
	if err := DB.Preload("Permissions").Where("code = ?", code).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

// GetAllRoles 获取所有角色
func GetAllRoles() ([]*Role, error) {
	var roles []*Role
	if err := DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
