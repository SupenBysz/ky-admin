package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/SupenBysz/ky-admin/pkg/database"
	"github.com/SupenBysz/ky-admin/pkg/logger"
	"github.com/SupenBysz/ky-admin/pkg/models"
	"gorm.io/gorm"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
	action := flag.String("action", "migrate", "操作类型：migrate(迁移)、seed(初始数据)、reset(重置数据库)、backup(备份数据库)")
	flag.Parse()

	// 初始化配置
	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	err = logger.InitLoggerWithConfig(conf)
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 创建数据库管理器并连接数据库
	dbManager := database.NewDBManager(conf)
	dbManager.SetLogger(logger.GetLogger())

	// 连接数据库
	if err := dbManager.Connect(); err != nil {
		logger.Fatal("连接数据库失败", logger.Err(err))
	}
	defer dbManager.Close()

	// 注册全局数据库管理器
	database.RegisterGlobalDBManager(dbManager)

	// 创建迁移管理器
	migrationManager, err := database.GetMigrationManager()
	if err != nil {
		logger.Fatal("创建迁移管理器失败", logger.Err(err))
	}

	// 根据操作类型执行不同的操作
	switch *action {
	case "migrate":
		if err := runMigration(migrationManager); err != nil {
			logger.Fatal("执行数据库迁移失败", logger.Err(err))
		}
		logger.Info("数据库迁移完成")

	case "seed":
		if err := runMigration(migrationManager); err != nil {
			logger.Fatal("执行数据库迁移失败", logger.Err(err))
		}

		if err := runSeed(migrationManager); err != nil {
			logger.Fatal("执行数据初始化失败", logger.Err(err))
		}
		logger.Info("数据初始化完成")

	case "reset":
		// 重置数据库（删除所有表并重新创建）
		if err := resetDatabase(migrationManager); err != nil {
			logger.Fatal("重置数据库失败", logger.Err(err))
		}

		if err := runMigration(migrationManager); err != nil {
			logger.Fatal("执行数据库迁移失败", logger.Err(err))
		}

		if err := runSeed(migrationManager); err != nil {
			logger.Fatal("执行数据初始化失败", logger.Err(err))
		}
		logger.Info("数据库重置完成")

	case "backup":
		backupPath, err := database.DatabaseBackup()
		if err != nil {
			logger.Fatal("备份数据库失败", logger.Err(err))
		}
		logger.Info(fmt.Sprintf("数据库备份完成: %s", backupPath))

	default:
		fmt.Println("无效的操作类型，可用选项: migrate, seed, reset, backup")
		os.Exit(1)
	}
}

// runMigration 执行数据库迁移
func runMigration(migrationManager *database.DBMigrationManager) error {
	logger.Info("开始执行数据库迁移")

	// 待迁移的模型
	models := []interface{}{
		&models.User{},
		&models.Role{},
		&models.Permission{},
	}

	return migrationManager.RunMigration(models...)
}

// runSeed 执行数据初始化
func runSeed(migrationManager *database.DBMigrationManager) error {
	logger.Info("开始初始化基础数据")

	return migrationManager.RunSeed(func(tx *gorm.DB) error {
		// 检查是否已有数据
		var count int64
		if err := tx.Model(&models.User{}).Count(&count).Error; err != nil {
			return err
		}

		// 已有数据，跳过初始化
		if count > 0 {
			logger.Info("已有数据，跳过初始化")
			return nil
		}

		// 创建超级管理员角色
		adminRole := &models.Role{
			Name:        "超级管理员",
			Code:        "admin",
			Description: "系统超级管理员",
			StatusModel: models.StatusModel{
				Status: uint8(models.StatusEnabled),
			},
		}

		if err := tx.Create(adminRole).Error; err != nil {
			return err
		}

		// 创建基础权限
		permissions := []models.Permission{
			{
				Name:        "系统管理",
				Code:        "system:manage",
				Type:        models.PermissionTypeSystem,
				Description: "系统管理权限",
				StatusModel: models.StatusModel{
					Status: uint8(models.StatusEnabled),
				},
			},
			{
				Name:        "用户管理",
				Code:        "user:manage",
				Type:        models.PermissionTypeMenu,
				Description: "用户管理权限",
				StatusModel: models.StatusModel{
					Status: uint8(models.StatusEnabled),
				},
			},
			{
				Name:        "角色管理",
				Code:        "role:manage",
				Type:        models.PermissionTypeMenu,
				Description: "角色管理权限",
				StatusModel: models.StatusModel{
					Status: uint8(models.StatusEnabled),
				},
			},
			{
				Name:        "权限管理",
				Code:        "permission:manage",
				Type:        models.PermissionTypeMenu,
				Description: "权限管理权限",
				StatusModel: models.StatusModel{
					Status: uint8(models.StatusEnabled),
				},
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 为超级管理员角色分配所有权限
		var allPermissionIDs []uint
		for _, p := range permissions {
			allPermissionIDs = append(allPermissionIDs, p.ID)
		}

		if err := adminRole.AssignPermissions(tx, allPermissionIDs); err != nil {
			return err
		}

		// 创建初始管理员用户
		adminUser := &models.User{
			Username: "admin",
			Email:    "admin@kysion.com",
			Nickname: "系统管理员",
			StatusModel: models.StatusModel{
				Status: uint8(models.StatusEnabled),
			},
			IsAdmin: true,
		}

		// 设置默认密码 admin123
		if err := adminUser.SetPassword("admin123"); err != nil {
			return err
		}

		if err := tx.Create(adminUser).Error; err != nil {
			return err
		}

		// 为用户分配角色
		if err := adminUser.AssignRole(tx, adminRole.ID); err != nil {
			return err
		}

		return nil
	})
}

// resetDatabase 重置数据库（删除所有表并重新创建）
func resetDatabase(migrationManager *database.DBMigrationManager) error {
	logger.Info("开始重置数据库")

	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("获取数据库连接失败")
	}

	// 获取所有表名
	var tables []string
	if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return err
	}

	// 禁用外键约束
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	// 删除所有表
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)).Error; err != nil {
			return err
		}
	}

	// 启用外键约束
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return err
	}

	logger.Info("数据库表已清空")
	return nil
}
