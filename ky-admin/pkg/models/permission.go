package models

import (
	"errors"

	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	StatusModel
	Name        string  `gorm:"type:varchar(50);not null;uniqueIndex:idx_permission_name;comment:权限名称"`
	Code        string  `gorm:"type:varchar(100);not null;uniqueIndex:idx_permission_code;comment:权限代码"`
	Description string  `gorm:"type:varchar(200);comment:权限描述"`
	Type        uint    `gorm:"default:2;comment:权限类型 1:系统 2:菜单 3:操作 4:数据"`
	ParentID    uint    `gorm:"default:0;comment:父权限ID"`
	Path        string  `gorm:"type:varchar(200);comment:路径"`
	Method      string  `gorm:"type:varchar(10);comment:HTTP方法,GET,POST,PUT,DELETE"`
	Sort        int     `gorm:"default:0;comment:排序"`
	Icon        string  `gorm:"type:varchar(100);comment:图标"`
	Component   string  `gorm:"type:varchar(100);comment:组件路径"`
	IsHidden    bool    `gorm:"default:false;comment:是否隐藏"`
	IsCache     bool    `gorm:"default:true;comment:是否缓存"`
	IsFrame     bool    `gorm:"default:false;comment:是否外链"`
	IsAffix     bool    `gorm:"default:false;comment:是否固定"`
	Roles       []*Role `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 表名
func (Permission) TableName() string {
	return "permissions"
}

// BeforeDelete 删除前钩子
func (p *Permission) BeforeDelete(tx *gorm.DB) error {
	// 系统权限不允许删除
	if p.Type == PermissionTypeSystem {
		return ErrSystemPermissionNotDeletable
	}

	// 检查是否有子权限
	var count int64
	if err := tx.Model(&Permission{}).Where("parent_id = ?", p.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrHasChildPermissions
	}

	return nil
}

// Enable 启用权限
func (p *Permission) Enable() error {
	if p.Status == uint8(StatusEnabled) {
		return nil
	}
	p.Status = uint8(StatusEnabled)
	return DB.Model(p).Update("status", uint8(StatusEnabled)).Error
}

// Disable 禁用权限
func (p *Permission) Disable() error {
	if p.Status == uint8(StatusDisabled) {
		return nil
	}
	p.Status = uint8(StatusDisabled)
	return DB.Model(p).Update("status", uint8(StatusDisabled)).Error
}

// IsEnabled 是否启用
func (p *Permission) IsEnabled() bool {
	return p.Status == uint8(StatusEnabled)
}

// IsSystemPermission 是否是系统权限
func (p *Permission) IsSystemPermission() bool {
	return p.Type == PermissionTypeSystem
}

// IsMenuPermission 是否是菜单权限
func (p *Permission) IsMenuPermission() bool {
	return p.Type == PermissionTypeMenu
}

// IsActionPermission 是否是操作权限
func (p *Permission) IsActionPermission() bool {
	return p.Type == PermissionTypeAction
}

// IsDataPermission 是否是数据权限
func (p *Permission) IsDataPermission() bool {
	return p.Type == PermissionTypeData
}

// GetPermissionByID 根据ID获取权限
func GetPermissionByID(id uint) (*Permission, error) {
	var permission Permission
	if err := DB.First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

// GetPermissionByCode 根据权限代码获取权限
func GetPermissionByCode(code string) (*Permission, error) {
	var permission Permission
	if err := DB.Where("code = ?", code).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

// GetPermissionsByType 根据类型获取权限列表
func GetPermissionsByType(permType uint) ([]*Permission, error) {
	var permissions []*Permission
	if err := DB.Where("type = ?", permType).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetMenuPermissions 获取菜单权限列表
func GetMenuPermissions() ([]*Permission, error) {
	return GetPermissionsByType(PermissionTypeMenu)
}

// GetAllPermissions 获取所有权限
func GetAllPermissions() ([]*Permission, error) {
	var permissions []*Permission
	if err := DB.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
