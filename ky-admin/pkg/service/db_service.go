package service

import (
	"context"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/models"
	"gorm.io/gorm"
)

// DBService 数据库服务
type DBService struct {
	db *gorm.DB
}

// NewDBService 创建数据库服务
func NewDBService(db *gorm.DB) *DBService {
	return &DBService{db: db}
}

// User 用户相关服务

// GetUserByID 通过ID获取用户
func (s *DBService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 通过用户名获取用户
func (s *DBService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func (s *DBService) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

// UpdateUser 更新用户
func (s *DBService) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}

// DeleteUser 删除用户
func (s *DBService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

// EnableUser 启用用户
func (s *DBService) EnableUser(id uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("status", uint8(models.StatusEnabled)).Error
}

// DisableUser 禁用用户
func (s *DBService) DisableUser(id uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("status", uint8(models.StatusDisabled)).Error
}

// UpdateUserLoginInfo 更新用户登录信息
func (s *DBService) UpdateUserLoginInfo(id uint, ip string) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": time.Now(),
		"last_login_ip": ip,
	}).Error
}

// ChangeUserPassword 修改用户密码
func (s *DBService) ChangeUserPassword(id uint, newPassword string) error {
	user := &models.User{
		StatusModel: models.StatusModel{
			BaseModel: models.BaseModel{
				ID: id,
			},
		},
		Password: newPassword,
	}
	return s.db.Model(user).Update("password", newPassword).Error
}

// AddUserRole 为用户添加角色
func (s *DBService) AddUserRole(userID, roleID uint) error {
	return s.db.Model(&models.User{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: userID}}}).
		Association("Roles").
		Append(&models.Role{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: roleID}}})
}

// RemoveUserRole 移除用户角色
func (s *DBService) RemoveUserRole(userID, roleID uint) error {
	return s.db.Model(&models.User{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: userID}}}).
		Association("Roles").
		Delete(&models.Role{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: roleID}}})
}

// ClearUserRoles 清空用户角色
func (s *DBService) ClearUserRoles(userID uint) error {
	return s.db.Model(&models.User{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: userID}}}).
		Association("Roles").
		Clear()
}

// Role 角色相关服务

// GetRoleByID 通过ID获取角色
func (s *DBService) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateRole 创建角色
func (s *DBService) CreateRole(role *models.Role) error {
	return s.db.Create(role).Error
}

// UpdateRole 更新角色
func (s *DBService) UpdateRole(role *models.Role) error {
	return s.db.Save(role).Error
}

// DeleteRole 删除角色
func (s *DBService) DeleteRole(id uint) error {
	return s.db.Delete(&models.Role{}, id).Error
}

// EnableRole 启用角色
func (s *DBService) EnableRole(id uint) error {
	return s.db.Model(&models.Role{}).Where("id = ?", id).Update("status", uint8(models.StatusEnabled)).Error
}

// DisableRole 禁用角色
func (s *DBService) DisableRole(id uint) error {
	return s.db.Model(&models.Role{}).Where("id = ?", id).Update("status", uint8(models.StatusDisabled)).Error
}

// AddRolePermission 为角色添加权限
func (s *DBService) AddRolePermission(roleID, permissionID uint) error {
	return s.db.Model(&models.Role{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: roleID}}}).
		Association("Permissions").
		Append(&models.Permission{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: permissionID}}})
}

// RemoveRolePermission 移除角色权限
func (s *DBService) RemoveRolePermission(roleID, permissionID uint) error {
	return s.db.Model(&models.Role{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: roleID}}}).
		Association("Permissions").
		Delete(&models.Permission{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: permissionID}}})
}

// ClearRolePermissions 清空角色权限
func (s *DBService) ClearRolePermissions(roleID uint) error {
	return s.db.Model(&models.Role{StatusModel: models.StatusModel{BaseModel: models.BaseModel{ID: roleID}}}).
		Association("Permissions").
		Clear()
}

// Permission 权限相关服务

// GetPermissionByID 通过ID获取权限
func (s *DBService) GetPermissionByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	if err := s.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// CreatePermission 创建权限
func (s *DBService) CreatePermission(permission *models.Permission) error {
	return s.db.Create(permission).Error
}

// UpdatePermission 更新权限
func (s *DBService) UpdatePermission(permission *models.Permission) error {
	return s.db.Save(permission).Error
}

// DeletePermission 删除权限
func (s *DBService) DeletePermission(id uint) error {
	return s.db.Delete(&models.Permission{}, id).Error
}

// EnablePermission 启用权限
func (s *DBService) EnablePermission(id uint) error {
	return s.db.Model(&models.Permission{}).Where("id = ?", id).Update("status", uint8(models.StatusEnabled)).Error
}

// DisablePermission 禁用权限
func (s *DBService) DisablePermission(id uint) error {
	return s.db.Model(&models.Permission{}).Where("id = ?", id).Update("status", uint8(models.StatusDisabled)).Error
}

// Transaction 数据库事务

// Transaction 事务处理
func (s *DBService) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return s.db.WithContext(ctx).Transaction(fn)
}
