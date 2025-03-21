package repository

import (
	"errors"

	"github.com/SupenBysz/ky-admin/internal/model"
	"gorm.io/gorm"
)

// RoleRepository 角色仓库接口
type RoleRepository interface {
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id uint) error
	FindByID(id uint) (*model.Role, error)
	FindByCode(code string) (*model.Role, error)
	FindAll() ([]*model.Role, error)
	FindByPage(query *model.RolePageQuery) ([]*model.Role, int64, error)
	FindByIDs(ids []uint) ([]*model.Role, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
	FindByTenantID(tenantID uint) ([]*model.Role, error)
}

// roleRepository 角色仓库实现
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create 创建角色
func (r *roleRepository) Create(role *model.Role) error {
	// 检查角色是否已存在
	var count int64
	if err := r.db.Model(&model.Role{}).Where("code = ?", role.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleAlreadyExists
	}

	// 开启事务
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// 创建角色
	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// Update 更新角色
func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *roleRepository) Delete(id uint) error {
	// 开启事务
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// 删除角色权限关联
	if err := tx.Where("role_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除用户角色关联
	if err := tx.Where("role_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除角色
	if err := tx.Delete(&model.Role{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// FindByID 根据ID查询角色
func (r *roleRepository) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

// FindByCode 根据编码查询角色
func (r *roleRepository) FindByCode(code string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

// FindAll 查询所有角色
func (r *roleRepository) FindAll() ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindByPage 分页查询角色
func (r *roleRepository) FindByPage(query *model.RolePageQuery) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	db := r.db.Model(&model.Role{})

	// 构建查询条件
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Code != "" {
		db = db.Where("code LIKE ?", "%"+query.Code+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.TenantID != nil {
		db = db.Where("tenant_id = ?", *query.TenantID)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// FindByIDs 根据ID列表查询角色
func (r *roleRepository) FindByIDs(ids []uint) ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.db.Where("id IN ?", ids).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// AssignPermissions 分配权限
func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 开启事务
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// 删除原有角色权限关联
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新的角色权限关联
	if len(permissionIDs) > 0 {
		var rolePermissions []model.RolePermission
		for _, permID := range permissionIDs {
			rolePermissions = append(rolePermissions, model.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
				TenantID:     1, // 默认租户，后续实现多租户后需要修改
			})
		}
		if err := tx.Create(&rolePermissions).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// FindByTenantID 根据租户ID查询角色
func (r *roleRepository) FindByTenantID(tenantID uint) ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.db.Where("tenant_id = ?", tenantID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
