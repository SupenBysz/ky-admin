package service

import (
	"strconv"

	"github.com/casbin/casbin/v2"
)

// PermissionService 权限服务接口
type PermissionService interface {
	// 检查用户是否具有某个角色
	HasRole(userID uint, role string) (bool, error)

	// 检查用户是否具有访问资源的权限
	HasPermission(userID uint, obj string, act string) (bool, error)

	// 为用户分配角色
	AssignRole(userID uint, role string) error

	// 移除用户角色
	RemoveRole(userID uint, role string) error

	// 为角色添加权限
	AddPermissionForRole(role string, obj string, act string) error

	// 移除角色权限
	RemovePermissionForRole(role string, obj string, act string) error

	// 获取用户所有角色
	GetRolesForUser(userID uint) ([]string, error)

	// 获取用户所有权限
	GetPermissionsForUser(userID uint) ([][]string, error)

	// 获取所有角色
	GetAllRoles() ([]string, error)

	// 重新加载权限策略
	ReloadPolicy() error
}

// permissionService 权限服务实现
type permissionService struct {
	enforcer *casbin.Enforcer
}

// NewPermissionService 创建权限服务
func NewPermissionService(enforcer *casbin.Enforcer) PermissionService {
	return &permissionService{
		enforcer: enforcer,
	}
}

// HasRole 检查用户是否具有某个角色
func (s *permissionService) HasRole(userID uint, role string) (bool, error) {
	return s.enforcer.HasRoleForUser(strconv.FormatUint(uint64(userID), 10), role)
}

// HasPermission 检查用户是否具有访问资源的权限
func (s *permissionService) HasPermission(userID uint, obj string, act string) (bool, error) {
	return s.enforcer.Enforce(strconv.FormatUint(uint64(userID), 10), obj, act)
}

// AssignRole 为用户分配角色
func (s *permissionService) AssignRole(userID uint, role string) error {
	_, err := s.enforcer.AddRoleForUser(strconv.FormatUint(uint64(userID), 10), role)
	if err != nil {
		return err
	}
	return s.enforcer.SavePolicy()
}

// RemoveRole 移除用户角色
func (s *permissionService) RemoveRole(userID uint, role string) error {
	_, err := s.enforcer.DeleteRoleForUser(strconv.FormatUint(uint64(userID), 10), role)
	if err != nil {
		return err
	}
	return s.enforcer.SavePolicy()
}

// AddPermissionForRole 为角色添加权限
func (s *permissionService) AddPermissionForRole(role string, obj string, act string) error {
	_, err := s.enforcer.AddPolicy(role, obj, act)
	if err != nil {
		return err
	}
	return s.enforcer.SavePolicy()
}

// RemovePermissionForRole 移除角色权限
func (s *permissionService) RemovePermissionForRole(role string, obj string, act string) error {
	_, err := s.enforcer.RemovePolicy(role, obj, act)
	if err != nil {
		return err
	}
	return s.enforcer.SavePolicy()
}

// GetRolesForUser 获取用户所有角色
func (s *permissionService) GetRolesForUser(userID uint) ([]string, error) {
	return s.enforcer.GetRolesForUser(strconv.FormatUint(uint64(userID), 10))
}

// GetPermissionsForUser 获取用户所有权限
func (s *permissionService) GetPermissionsForUser(userID uint) ([][]string, error) {
	return s.enforcer.GetImplicitPermissionsForUser(strconv.FormatUint(uint64(userID), 10))
}

// GetAllRoles 获取所有角色
func (s *permissionService) GetAllRoles() ([]string, error) {
	return s.enforcer.GetAllRoles()
}

// ReloadPolicy 重新加载权限策略
func (s *permissionService) ReloadPolicy() error {
	return s.enforcer.LoadPolicy()
}
