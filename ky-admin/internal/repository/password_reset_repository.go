package repository

import (
	"time"

	"github.com/SupenBysz/ky-admin/internal/model"
	"gorm.io/gorm"
)

// PasswordResetRepository 密码重置令牌仓库接口
type PasswordResetRepository interface {
	Create(userID uint, email, token string, expiration time.Time) (*model.PasswordReset, error)
	FindByToken(token string) (*model.PasswordReset, error)
	MarkTokenAsUsed(token string) error
	DeleteExpired() error
}

// passwordResetRepositoryImpl 密码重置令牌仓库实现
type passwordResetRepositoryImpl struct {
	db *gorm.DB
}

// NewPasswordResetRepository 创建密码重置令牌仓库实例
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepositoryImpl{db: db}
}

// Create 创建密码重置令牌
func (r *passwordResetRepositoryImpl) Create(userID uint, email, token string, expiration time.Time) (*model.PasswordReset, error) {
	// 删除该用户之前的所有重置令牌
	if err := r.db.Where("user_id = ?", userID).Delete(&model.PasswordReset{}).Error; err != nil {
		return nil, err
	}

	reset := &model.PasswordReset{
		UserID:    userID,
		Email:     email,
		Token:     token,
		ExpiresAt: expiration,
	}

	if err := r.db.Create(reset).Error; err != nil {
		return nil, err
	}

	return reset, nil
}

// FindByToken 通过令牌查找密码重置记录
func (r *passwordResetRepositoryImpl) FindByToken(token string) (*model.PasswordReset, error) {
	var reset model.PasswordReset
	if err := r.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&reset).Error; err != nil {
		return nil, err
	}
	return &reset, nil
}

// MarkTokenAsUsed 标记令牌为已使用
func (r *passwordResetRepositoryImpl) MarkTokenAsUsed(token string) error {
	return r.db.Model(&model.PasswordReset{}).Where("token = ?", token).Updates(map[string]interface{}{
		"used":       true,
		"updated_at": time.Now(),
	}).Error
}

// DeleteExpired 删除过期的令牌
func (r *passwordResetRepositoryImpl) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.PasswordReset{}).Error
}
