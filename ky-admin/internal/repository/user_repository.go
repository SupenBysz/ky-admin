package repository

import (
	stderrors "errors"

	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAll(page, pageSize int) ([]model.User, int64, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	FindByEmail(email string) (*model.User, error)
	FindByPage(query *dto.UserPageQueryDTO) ([]model.User, int64, error)
	AssignRoles(userID uint, roleIDs []uint) error
	FindUserRoles(userID uint) ([]model.Role, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByID 通过ID查找用户
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername 通过用户名查找用户
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindAll 查找所有用户
func (r *userRepository) FindAll(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var count int64

	offset := (page - 1) * pageSize
	if err := r.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	// 检查用户名是否已存在
	var count int64
	if err := r.db.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrUserAlreadyExists
	}

	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	// 检查用户是否存在
	var exists bool
	if err := r.db.Model(&model.User{}).Select("1").Where("id = ?", user.ID).Limit(1).Find(&exists).Error; err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// FindByEmail 根据邮箱查询用户
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByPage 分页查询用户
func (r *userRepository) FindByPage(query *dto.UserPageQueryDTO) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := r.db.Model(&model.User{})

	// 构建查询条件
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+query.Nickname+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.TenantID != nil {
		db = db.Where("tenant_id = ?", *query.TenantID)
	}
	if query.DeptID != nil {
		db = db.Where("dept_id = ?", *query.DeptID)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// AssignRoles 分配用户角色
func (r *userRepository) AssignRoles(userID uint, roleIDs []uint) error {
	// 开启事务
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// 查找用户是否存在
	var user model.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 删除原有的用户角色关联
	if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新的用户角色关联
	if len(roleIDs) > 0 {
		var userRoles []model.UserRole
		for _, roleID := range roleIDs {
			userRoles = append(userRoles, model.UserRole{
				UserID: userID,
				RoleID: roleID,
			})
		}
		if err := tx.Create(&userRoles).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// FindUserRoles 查询用户的角色
func (r *userRepository) FindUserRoles(userID uint) ([]model.Role, error) {
	var roles []model.Role

	// 通过用户角色关联表查询角色
	err := r.db.Table("sys_role").
		Joins("JOIN sys_user_role ON sys_role.id = sys_user_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.deleted_at IS NULL", userID).
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}
