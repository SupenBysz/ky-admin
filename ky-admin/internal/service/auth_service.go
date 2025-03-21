package service

import (
	"time"

	"github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/model"
	"github.com/SupenBysz/ky-admin/internal/repository"
	"github.com/SupenBysz/ky-admin/pkg/jwt"
	"github.com/SupenBysz/ky-admin/pkg/utils"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(loginDTO *dto.LoginDTO) (string, string, *dto.UserDTO, error)
	Register(registerDTO *dto.RegisterDTO) (*dto.UserDTO, error)
	RefreshToken(refreshToken string) (string, string, error)
	GetUserInfo(userID uint) (*dto.UserDTO, error)
}

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository
	jwtSvc   *jwt.JWTService
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository, jwtSvc *jwt.JWTService) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

// Login 用户登录
func (s *authService) Login(loginDTO *dto.LoginDTO) (string, string, *dto.UserDTO, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(loginDTO.Username)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return "", "", nil, errors.ErrInvalidCredentials
		}
		return "", "", nil, errors.NewSystemError(err)
	}

	// 检查用户状态
	if user.Status != 1 {
		return "", "", nil, errors.ErrUserDisabled
	}

	// 验证密码
	if !utils.CheckPassword(user.Password, loginDTO.Password) {
		return "", "", nil, errors.ErrInvalidCredentials
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(user); err != nil {
		return "", "", nil, errors.NewSystemError(err)
	}

	// 生成令牌
	token, err := s.jwtSvc.GenerateToken(user.ID, user.TenantID, user.Username)
	if err != nil {
		return "", "", nil, errors.NewSystemError(err)
	}

	// 生成刷新令牌
	refreshToken, err := s.jwtSvc.GenerateRefreshToken(user.ID, user.TenantID, user.Username)
	if err != nil {
		return "", "", nil, errors.NewSystemError(err)
	}

	// 查询用户角色
	roles, _ := s.userRepo.FindUserRoles(user.ID)

	// 转换角色为DTO
	roleDTOs := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, dto.RoleDTO{
			ID:   role.ID,
			Name: role.Name,
			Code: role.Code,
		})
	}

	// 转换为DTO
	userDTO := &dto.UserDTO{
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
		Roles:     roleDTOs,
	}

	return token, refreshToken, userDTO, nil
}

// Register 用户注册
func (s *authService) Register(registerDTO *dto.RegisterDTO) (*dto.UserDTO, error) {
	// 检查密码强度（简单实现，实际应用需要更复杂的密码策略）
	if len(registerDTO.Password) < 6 {
		return nil, errors.ErrWeakPassword
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(registerDTO.Password)
	if err != nil {
		return nil, errors.NewSystemError(err)
	}

	// 创建用户
	user := &model.User{
		Username: registerDTO.Username,
		Password: hashedPassword,
		Nickname: registerDTO.Nickname,
		Email:    registerDTO.Email,
		Phone:    registerDTO.Phone,
		Status:   1, // 默认启用
		TenantID: 1, // 默认租户，后续实现多租户后需要修改
	}

	// 保存用户
	if err := s.userRepo.Create(user); err != nil {
		if err == repository.ErrUserAlreadyExists {
			return nil, errors.ErrUserAlreadyExists
		}
		return nil, errors.NewSystemError(err)
	}

	// 分配角色（如果有）
	if len(registerDTO.RoleIDs) > 0 {
		if err := s.userRepo.AssignRoles(user.ID, registerDTO.RoleIDs); err != nil {
			return nil, errors.NewSystemError(err)
		}
	}

	// 查询用户角色
	roles, _ := s.userRepo.FindUserRoles(user.ID)

	// 转换角色为DTO
	roleDTOs := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, dto.RoleDTO{
			ID:   role.ID,
			Name: role.Name,
			Code: role.Code,
		})
	}

	// 转换为DTO
	userDTO := &dto.UserDTO{
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
		Roles:     roleDTOs,
	}

	return userDTO, nil
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	return s.jwtSvc.RefreshToken(refreshToken)
}

// GetUserInfo 获取用户信息
func (s *authService) GetUserInfo(userID uint) (*dto.UserDTO, error) {
	// 查找用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.NewSystemError(err)
	}

	// 查询用户角色
	roles, _ := s.userRepo.FindUserRoles(user.ID)

	// 转换角色为DTO
	roleDTOs := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, dto.RoleDTO{
			ID:   role.ID,
			Name: role.Name,
			Code: role.Code,
		})
	}

	// 转换为DTO
	userDTO := &dto.UserDTO{
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
		Roles:     roleDTOs,
	}

	return userDTO, nil
}
