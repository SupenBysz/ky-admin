package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"

	"github.com/SupenBysz/ky-admin/internal/config"
	"github.com/jordan-wright/email"
)

// EmailService 邮件服务接口
type EmailService interface {
	// SendEmail 发送邮件
	SendEmail(to []string, subject, content string) error
	// SendHTMLEmail 发送HTML邮件
	SendHTMLEmail(to []string, subject, htmlContent string) error
	// SendTemplateEmail 发送模板邮件
	SendTemplateEmail(to []string, subject, templateName string, data interface{}) error
	// SendPasswordResetEmail 发送密码重置邮件
	SendPasswordResetEmail(to string, resetLink string, username string) error
}

// emailService 邮件服务实现
type emailService struct {
	config config.EmailConfig
}

// NewEmailService 创建邮件服务
func NewEmailService(config config.EmailConfig) EmailService {
	return &emailService{
		config: config,
	}
}

// SendEmail 发送邮件
func (s *emailService) SendEmail(to []string, subject, content string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.config.SenderName, s.config.SenderEmail)
	e.To = to
	e.Subject = subject
	e.Text = []byte(content)

	return e.Send(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost),
	)
}

// SendHTMLEmail 发送HTML邮件
func (s *emailService) SendHTMLEmail(to []string, subject, htmlContent string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.config.SenderName, s.config.SenderEmail)
	e.To = to
	e.Subject = subject
	e.HTML = []byte(htmlContent)

	return e.Send(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost),
	)
}

// SendTemplateEmail 发送模板邮件
func (s *emailService) SendTemplateEmail(to []string, subject, templateName string, data interface{}) error {
	templatePath := filepath.Join(s.config.TemplateDir, templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("解析邮件模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行邮件模板失败: %w", err)
	}

	return s.SendHTMLEmail(to, subject, buf.String())
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *emailService) SendPasswordResetEmail(to string, resetLink string, username string) error {
	subject := "密码重置"
	data := map[string]interface{}{
		"Username":  username,
		"ResetLink": resetLink,
		"ExpiresIn": "24小时",
	}

	// 使用密码重置邮件模板
	return s.SendTemplateEmail([]string{to}, subject, "password_reset.html", data)
}
