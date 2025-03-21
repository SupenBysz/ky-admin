package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/SupenBysz/ky-admin/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// RunMigration 执行数据库迁移和初始化
func RunMigration() {
	// 从环境变量获取数据库连接信息，或使用默认值
	dbHost := getEnv("DB_HOST", "10.68.73.213")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "kysion_dev")
	dbPassword := getEnv("DB_PASSWORD", "TrY6eK8Bt28KwYii")
	dbName := getEnv("DB_NAME", "kysiondb_dev")

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("数据库连接成功，开始执行迁移...")

	// 执行迁移
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	); err != nil {
		log.Fatalf("自动迁移表结构失败: %v", err)
	}

	fmt.Println("表结构迁移完成，开始初始化基础数据...")

	// 检查是否需要初始化角色数据
	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)

	if roleCount == 0 {
		// 创建基础角色
		roles := []model.Role{
			{Name: "admin", Code: "admin", Description: "系统管理员", Sort: 1},
			{Name: "manager", Code: "manager", Description: "管理员", Sort: 2},
			{Name: "user", Code: "user", Description: "普通用户", Sort: 3},
		}

		if err := db.Create(&roles).Error; err != nil {
			log.Fatalf("创建角色失败: %v", err)
		}

		fmt.Println("基础角色创建成功")
	} else {
		fmt.Println("角色数据已存在，跳过初始化")
	}

	// 检查是否需要初始化用户数据
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)

	if userCount == 0 {
		// 创建超级管理员用户
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("密码加密失败: %v", err)
		}

		adminUser := model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "系统管理员",
			Email:    "admin@example.com",
			Status:   1,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatalf("创建管理员用户失败: %v", err)
		}

		// 查询admin角色
		var adminRole model.Role
		if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
			log.Fatalf("查询管理员角色失败: %v", err)
		}

		// 为管理员分配角色
		if err := db.Exec("INSERT INTO sys_user_role (user_id, role_id) VALUES (?, ?)", adminUser.ID, adminRole.ID).Error; err != nil {
			log.Fatalf("分配管理员角色失败: %v", err)
		}

		fmt.Println("超级管理员用户创建成功")
	} else {
		fmt.Println("用户数据已存在，跳过初始化")
	}

	fmt.Println("数据库初始化完成！")
}

// getEnv 获取环境变量值，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
