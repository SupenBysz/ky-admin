package config

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey     string `mapstructure:"secret_key"`     // JWT密钥
	TokenExpire   int64  `mapstructure:"token_expire"`   // 访问令牌过期时间（小时）
	RefreshExpire int64  `mapstructure:"refresh_expire"` // 刷新令牌过期时间（小时）
	Issuer        string `mapstructure:"issuer"`         // 签发者
}
