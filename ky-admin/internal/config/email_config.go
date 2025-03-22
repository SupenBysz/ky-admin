package config

// EmailConfig 邮件配置
type EmailConfig struct {
	// SMTPHost SMTP服务器地址
	SMTPHost string `yaml:"smtp_host" json:"smtp_host"`
	// SMTPPort SMTP服务器端口
	SMTPPort int `yaml:"smtp_port" json:"smtp_port"`
	// SMTPUsername SMTP用户名
	SMTPUsername string `yaml:"smtp_username" json:"smtp_username"`
	// SMTPPassword SMTP密码
	SMTPPassword string `yaml:"smtp_password" json:"smtp_password"`
	// SenderEmail 发件人邮箱
	SenderEmail string `yaml:"sender_email" json:"sender_email"`
	// SenderName 发件人名称
	SenderName string `yaml:"sender_name" json:"sender_name"`
	// TemplateDir 邮件模板目录
	TemplateDir string `yaml:"template_dir" json:"template_dir"`
	// EnableTLS 是否启用TLS
	EnableTLS bool `yaml:"enable_tls" json:"enable_tls"`
}
