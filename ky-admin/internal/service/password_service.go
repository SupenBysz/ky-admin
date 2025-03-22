package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/repository"
	"github.com/SupenBysz/ky-admin/pkg/utils"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrInvalidResetToken = errors.New("无效的重置令牌")
	ErrExpiredResetToken = errors.New("重置令牌已过期")
	ErrTokenAlreadyUsed  = errors.New("令牌已被使用")
)

// PasswordService 密码服务接口
type PasswordService interface {
	// ForgotPassword 忘记密码，生成重置令牌并发送邮件
	ForgotPassword(dto *dto.ForgotPasswordDTO) error
	// ResetPassword 重置密码
	ResetPassword(dto *dto.ResetPasswordDTO) error
	// ValidateResetToken 验证重置令牌
	ValidateResetToken(token string) (bool, error)
}

// passwordService 密码服务实现
type passwordService struct {
	userRepo          repository.UserRepository
	passwordResetRepo repository.PasswordResetRepository
	emailSvc          EmailService
	appURL            string
}

// NewPasswordService 创建密码服务
func NewPasswordService(
	userRepo repository.UserRepository,
	passwordResetRepo repository.PasswordResetRepository,
	emailSvc EmailService,
	appURL string,
) PasswordService {
	return &passwordService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		emailSvc:          emailSvc,
		appURL:            appURL,
	}
}

// generateResetToken 生成重置令牌
func (s *passwordService) generateResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ForgotPassword 忘记密码，生成重置令牌并发送邮件
func (s *passwordService) ForgotPassword(dto *dto.ForgotPasswordDTO) error {
	// 查找用户
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return ErrUserNotFound
	}

	// 生成重置令牌
	token, err := s.generateResetToken()
	if err != nil {
		return err
	}

	// 设置过期时间（24小时后）
	expiresAt := time.Now().Add(24 * time.Hour)

	// 保存令牌到数据库
	_, err = s.passwordResetRepo.Create(user.ID, user.Email, token, expiresAt)
	if err != nil {
		return err
	}

	// 生成重置链接
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.appURL, token)

	// 发送重置密码邮件
	return s.emailSvc.SendPasswordResetEmail(user.Email, resetLink, user.Username)
}

// ValidateResetToken 验证重置令牌
func (s *passwordService) ValidateResetToken(token string) (bool, error) {
	resetToken, err := s.passwordResetRepo.FindByToken(token)
	if err != nil {
		return false, ErrInvalidResetToken
	}

	if resetToken.Used {
		return false, ErrTokenAlreadyUsed
	}

	if resetToken.ExpiresAt.Before(time.Now()) {
		return false, ErrExpiredResetToken
	}

	return true, nil
}

// ResetPassword 重置密码
func (s *passwordService) ResetPassword(dto *dto.ResetPasswordDTO) error {
	// 验证令牌
	resetToken, err := s.passwordResetRepo.FindByToken(dto.Token)
	if err != nil {
		return ErrInvalidResetToken
	}

	if resetToken.Used {
		return ErrTokenAlreadyUsed
	}

	if resetToken.ExpiresAt.Before(time.Now()) {
		return ErrExpiredResetToken
	}

	// 查找用户
	user, err := s.userRepo.FindByID(resetToken.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	// 更新密码
	hashedPassword, err := utils.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 标记令牌为已使用
	return s.passwordResetRepo.MarkTokenAsUsed(dto.Token)
}
