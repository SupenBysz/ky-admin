package service

import (
	"github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/model"
	"github.com/SupenBysz/ky-admin/internal/repository"
	"github.com/SupenBysz/ky-admin/pkg/utils"
)

// UserService 用户服务接口
type UserService interface {
	GetByID(id uint) (*model.User, error)
	GetUsers(page, pageSize int) ([]model.User, int64, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
	Login(username, password string) (*model.User, error)
	GetUserRoles(userID uint) ([]model.Role, error)
	GetUserWithRoles(userID uint) (*dto.UserDTO, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
	CreateUserWithRoles(req dto.CreateUserDTO) (uint, error)
	AssignRoles(userID uint, roleIDs []uint) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUsers 获取用户列表
func (s *userService) GetUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.FindAll(page, pageSize)
}

// CreateUser 创建用户
func (s *userService) CreateUser(user *model.User) error {
	// 加密密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.NewSystemError(err)
	}
	user.Password = hashedPassword
	return s.userRepo.Create(user)
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(user *model.User) error {
	// 如果密码不为空，需要加密
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return errors.NewSystemError(err)
		}
		user.Password = hashedPassword
	}
	return s.userRepo.Update(user)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

// Login 用户登录
func (s *userService) Login(username, password string) (*model.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.NewAuthError(errors.ErrInvalidCredentials)
	}

	// 验证密码
	if !utils.CheckPassword(user.Password, password) {
		return nil, errors.NewAuthError(errors.ErrInvalidCredentials)
	}

	return user, nil
}

// GetUserRoles 获取用户角色
func (s *userService) GetUserRoles(userID uint) ([]model.Role, error) {
	// 查询用户角色
	return s.userRepo.FindUserRoles(userID)
}

// GetUserWithRoles 获取带有角色信息的用户
func (s *userService) GetUserWithRoles(userID uint) (*dto.UserDTO, error) {
	// 查询用户信息
	user, err := s.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 查询用户角色
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	rolesDTOs := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		rolesDTOs = append(rolesDTOs, dto.RoleDTO{
			ID:   role.ID,
			Name: role.Name,
			Code: role.Code,
		})
	}

	return &dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Status:    int8(user.Status),
		TenantID:  user.TenantID,
		DeptID:    user.DeptID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Roles:     rolesDTOs,
	}, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 查询用户
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(user.Password, oldPassword) {
		return errors.NewUserError(errors.ErrOldPasswordWrong)
	}

	// 验证新密码强度
	if len(newPassword) < 6 {
		return errors.NewUserError(errors.ErrWeakPassword)
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.NewSystemError(err)
	}

	// 更新密码
	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// CreateUserWithRoles 创建用户并分配角色
func (s *userService) CreateUserWithRoles(req dto.CreateUserDTO) (uint, error) {
	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
		TenantID: req.TenantID,
		DeptID:   req.DeptID,
	}

	// 创建用户
	if err := s.CreateUser(user); err != nil {
		return 0, err
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		if err := s.AssignRoles(user.ID, req.RoleIDs); err != nil {
			return 0, err
		}
	}

	return user.ID, nil
}

// AssignRoles 分配角色
func (s *userService) AssignRoles(userID uint, roleIDs []uint) error {
	return s.userRepo.AssignRoles(userID, roleIDs)
}
