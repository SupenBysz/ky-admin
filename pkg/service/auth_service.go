package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/SupenBysz/ky-admin/pkg/logger"
	"github.com/SupenBysz/ky-admin/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// AuthService 认证服务
type AuthService struct {
	dbService *DBService
}

// NewAuthService 创建认证服务
func NewAuthService(dbService *DBService) *AuthService {
	return &AuthService{
		dbService: dbService,
	}
}

// TokenClaims JWT令牌声明
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IP       string `json:"ip"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string       `json:"token"`
	TokenType string       `json:"token_type"`
	ExpiresIn int64        `json:"expires_in"`
	User      *models.User `json:"user"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 根据用户名查找用户
	user, err := s.dbService.GetUserByUsername(req.Username)
	if err != nil {
		logger.Warn("登录失败，用户不存在", zap.String("username", req.Username), zap.Error(err))
		return nil, models.ErrUserNotFound
	}

	// 检查用户状态
	if user.Status == uint8(models.StatusDisabled) {
		logger.Warn("登录失败，用户已禁用", zap.String("username", req.Username))
		return nil, models.ErrUserDisabled
	}

	if user.Status == uint8(models.StatusDeleted) {
		logger.Warn("登录失败，用户已锁定", zap.String("username", req.Username))
		return nil, models.ErrUserLocked
	}

	// 验证密码
	if !user.ValidatePassword(req.Password) {
		logger.Warn("登录失败，密码错误", zap.String("username", req.Username))
		return nil, models.ErrPasswordIncorrect
	}

	// 更新登录信息
	if req.IP != "" {
		if err := s.dbService.UpdateUserLoginInfo(user.ID, req.IP); err != nil {
			logger.Error("更新登录信息失败", zap.Error(err))
		}
	}

	// 生成JWT令牌
	token, expiresIn, err := s.GenerateToken(user)
	if err != nil {
		logger.Error("生成令牌失败", zap.Error(err))
		return nil, err
	}

	// 预加载用户的角色
	if err := s.dbService.db.Preload("Roles").First(user, user.ID).Error; err != nil {
		logger.Error("加载用户角色失败", zap.Error(err))
	}

	return &LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
		User:      user,
	}, nil
}

// GenerateToken 生成JWT令牌
func (s *AuthService) GenerateToken(user *models.User) (string, int64, error) {
	cfg := config.GetConfig()
	secretKey := cfg.GetString("jwt.secret")
	if secretKey == "" {
		secretKey = "ky-admin-secret-key" // 默认密钥
	}

	expiresIn := int64(cfg.GetInt("jwt.expires_in"))
	if expiresIn <= 0 {
		expiresIn = 3600 * 24 // 默认24小时
	}

	// 设置令牌过期时间
	expirationTime := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// 创建Token声明
	claims := &TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
			Issuer:    "ky-admin",
		},
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名Token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}

// VerifyToken 验证JWT令牌
func (s *AuthService) VerifyToken(tokenString string) (*TokenClaims, error) {
	cfg := config.GetConfig()
	secretKey := cfg.GetString("jwt.secret")
	if secretKey == "" {
		secretKey = "ky-admin-secret-key" // 默认密钥
	}

	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名方法: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// GetUserFromToken 从令牌中获取用户信息
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	// 验证令牌
	claims, err := s.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 根据用户ID获取用户信息
	user, err := s.dbService.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// HasPermission 检查用户是否拥有指定权限
func (s *AuthService) HasPermission(userID uint, permissionCode string) (bool, error) {
	var user models.User
	err := s.dbService.db.Preload("Roles.Permissions").First(&user, userID).Error
	if err != nil {
		return false, err
	}

	// 检查用户的所有角色
	for _, role := range user.Roles {
		// 检查角色状态
		if role.Status != uint8(models.StatusEnabled) {
			continue
		}

		// 检查角色的所有权限
		for _, permission := range role.Permissions {
			// 检查权限状态和编码
			if permission.Status == uint8(models.StatusEnabled) && permission.Code == permissionCode {
				return true, nil
			}
		}
	}

	return false, nil
}
