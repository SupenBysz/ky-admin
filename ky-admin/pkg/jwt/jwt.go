package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义错误类型
var (
	ErrTokenExpired     = errors.New("令牌已过期")
	ErrTokenNotValidYet = errors.New("令牌尚未生效")
	ErrTokenMalformed   = errors.New("令牌格式错误")
	ErrTokenInvalid     = errors.New("无效的令牌")
)

// CustomClaims 自定义JWT声明结构
type CustomClaims struct {
	UserID   uint   `json:"userId"`
	TenantID uint   `json:"tenantId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTService JWT服务
type JWTService struct {
	secretKey     []byte
	tokenExpire   time.Duration
	refreshExpire time.Duration
	issuer        string
}

// New 创建JWTService实例
func New(secretKey string, tokenExpire, refreshExpire time.Duration, issuer string) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		tokenExpire:   tokenExpire,
		refreshExpire: refreshExpire,
		issuer:        issuer,
	}
}

// GenerateToken 生成令牌
func (j *JWTService) GenerateToken(userID uint, tenantID uint, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		TenantID: tenantID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateRefreshToken 生成刷新令牌
func (j *JWTService) GenerateRefreshToken(userID uint, tenantID uint, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		TenantID: tenantID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ParseToken 解析令牌
func (j *JWTService) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		// 根据错误类型返回特定错误
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		} else {
			return nil, ErrTokenInvalid
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新令牌
func (j *JWTService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := j.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// 生成新的访问令牌
	token, err := j.GenerateToken(claims.UserID, claims.TenantID, claims.Username)
	if err != nil {
		return "", "", err
	}

	// 生成新的刷新令牌
	newRefreshToken, err := j.GenerateRefreshToken(claims.UserID, claims.TenantID, claims.Username)
	if err != nil {
		return "", "", err
	}

	return token, newRefreshToken, nil
}
