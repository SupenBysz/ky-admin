package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// Encryptor 配置加密器
type Encryptor struct {
	key []byte
}

// NewEncryptor 创建配置加密器实例
func NewEncryptor(secretKey string) *Encryptor {
	// 使用SHA-256散列确保密钥长度
	hash := sha256.Sum256([]byte(secretKey))
	return &Encryptor{
		key: hash[:],
	}
}

// Encrypt 加密字符串
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// 创建GCM分组模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 创建随机数
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密字符串
func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 检查数据长度
	if len(data) < gcm.NonceSize() {
		return "", errors.New("密文长度不足")
	}

	// 提取nonce
	nonce, ciphertextBytes := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// LoadEncryptedConfig 加载加密的配置项
func (c *Config) LoadEncryptedConfig(key, encryptedValue string) (string, error) {
	// 获取加密密钥
	secretKey := c.GetString("security.encryption_key")
	if secretKey == "" {
		secretKey = "default-encryption-key" // 默认密钥，应在生产环境中修改
	}

	// 创建加密器
	encryptor := NewEncryptor(secretKey)

	// 解密值
	if encryptedValue == "" {
		// 尝试从配置中获取
		encryptedValue = c.GetString(key)
	}

	decrypted, err := encryptor.Decrypt(encryptedValue)
	if err != nil {
		return "", err
	}

	return decrypted, nil
}

// SaveEncryptedConfig 保存加密的配置项
func (c *Config) SaveEncryptedConfig(key, plainValue string) (string, error) {
	// 获取加密密钥
	secretKey := c.GetString("security.encryption_key")
	if secretKey == "" {
		secretKey = "default-encryption-key" // 默认密钥，应在生产环境中修改
	}

	// 创建加密器
	encryptor := NewEncryptor(secretKey)

	// 加密值
	encrypted, err := encryptor.Encrypt(plainValue)
	if err != nil {
		return "", err
	}

	// 保存到配置
	c.Set(key, encrypted)

	return encrypted, nil
}
