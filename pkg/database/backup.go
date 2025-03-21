package database

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/logger"
	"go.uber.org/zap"
)

// DatabaseBackup 执行数据库备份
func DatabaseBackup() (string, error) {
	manager, err := GetGlobalDBManager()
	if err != nil {
		return "", err
	}

	// 获取数据库配置
	dbConfig := manager.GetBackupDBInfo()

	// 获取备份目录
	backupDir := "./backups"

	// 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		logger.Error("创建备份目录失败", zap.Error(err))
		return "", err
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("%s_%s.sql", dbConfig.Database, timestamp))

	// 执行备份命令
	logger.Info("开始执行数据库备份",
		zap.String("database", dbConfig.Database),
		zap.String("file", backupFile),
	)

	cmd := exec.Command("mysqldump",
		"-h", dbConfig.Host,
		"-P", strconv.Itoa(dbConfig.Port),
		"-u", dbConfig.Username,
		fmt.Sprintf("-p%s", dbConfig.Password),
		dbConfig.Database,
		"-r", backupFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("数据库备份失败", zap.Error(err), zap.ByteString("output", output))
		return "", err
	}

	logger.Info("数据库备份完成", zap.String("file", backupFile))
	return backupFile, nil
}

// ListBackups 列出所有备份文件
func ListBackups() ([]string, error) {
	backupDir := "./backups"

	// 确保备份目录存在
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	files, err := os.ReadDir(backupDir)
	if err != nil {
		logger.Error("读取备份目录失败", zap.Error(err))
		return nil, err
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			backups = append(backups, filepath.Join(backupDir, file.Name()))
		}
	}

	return backups, nil
}

// RestoreFromBackup 从备份文件恢复数据库
func RestoreFromBackup(backupFile string) error {
	manager, err := GetGlobalDBManager()
	if err != nil {
		return err
	}

	// 检查备份文件是否存在
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupFile)
	}

	// 获取数据库配置
	dbConfig := manager.GetBackupDBInfo()

	logger.Info("开始从备份文件恢复数据库",
		zap.String("database", dbConfig.Database),
		zap.String("file", backupFile),
	)

	// 执行恢复命令
	cmd := exec.Command("mysql",
		"-h", dbConfig.Host,
		"-P", strconv.Itoa(dbConfig.Port),
		"-u", dbConfig.Username,
		fmt.Sprintf("-p%s", dbConfig.Password),
		dbConfig.Database,
		"-e", fmt.Sprintf("source %s", backupFile))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("数据库恢复失败", zap.Error(err), zap.ByteString("output", output))
		return err
	}

	logger.Info("数据库恢复完成")
	return nil
}
