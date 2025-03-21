package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/SupenBysz/ky-admin/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("令牌无效")
	ErrExpiredToken = errors.New("令牌已过期")
)

// JWTClaims 自定义JWT Claims
type JWTClaims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}

// JWTService JWT服务接口
type JWTService interface {
	GenerateToken(userID uint) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	ValidateRefreshToken(refreshToken string) (*JWTClaims, error)
}

// jwtService JWT服务实现
type jwtService struct {
	config config.JWTConfig
}

// NewJWTService 创建JWT服务
func NewJWTService(config config.JWTConfig) JWTService {
	return &jwtService{
		config: config,
	}
}

// GenerateToken 生成访问令牌
func (s *jwtService) GenerateToken(userID uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.TokenExpire))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// GenerateRefreshToken 生成刷新令牌
func (s *jwtService) GenerateRefreshToken(userID uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.RefreshExpire))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// ValidateToken 验证访问令牌
func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateRefreshToken 验证刷新令牌
func (s *jwtService) ValidateRefreshToken(refreshToken string) (*JWTClaims, error) {
	return s.ValidateToken(refreshToken)
}
