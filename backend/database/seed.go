package database

import (
	"fmt"
	"log"
	"progress-wall-backend/models"
	"progress-wall-backend/utils"

	"gorm.io/gorm"
)

// Seed 初始化基础数据
func Seed(db *gorm.DB) error {
	log.Println("开始初始化基础数据...")

	// 创建默认用户
	if err := createDefaultUsers(db); err != nil {
		return err
	}

	// 创建默认权限
	if err := createDefaultPermissions(db); err != nil {
		return err
	}

	log.Println("基础数据初始化完成")
	return nil
}

// createDefaultUsers 创建默认用户
// 该函数负责在系统初始化时创建默认的管理员用户
// 参数: db - 数据库连接实例
// 返回: error - 如果创建失败则返回错误，否则返回nil
func createDefaultUsers(db *gorm.DB) error {
	// 检查admin用户是否已存在
	var existingUser models.User
	result := db.Where("username = ?", "admin").First(&existingUser)
	if result.Error == nil {
		log.Println("admin用户已存在，跳过创建")
		return nil
	}

	// 加密管理员密码
	adminPasswordHash, err := utils.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("加密管理员密码失败: %w", err)
	}

	// 创建默认管理员用户
	adminUser := models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: adminPasswordHash,
		Nickname: "系统管理员",
		Status:   models.UserStatusEnabled,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("创建默认管理员用户失败: %w", err)
	}

	log.Println("默认用户创建成功: admin")
	return nil
}

// createDefaultPermissions 创建默认权限组
// 该函数负责在系统初始化时创建默认的权限组，包括管理员权限组和访客权限组
// 参数: db - 数据库连接实例
// 返回: error - 如果创建失败则返回错误，否则返回nil
func createDefaultPermissions(db *gorm.DB) error {
	adminPermissionGroup := models.UserPermissionGroup{
		GroupName:         "admin",
		PermissionStrings: "*",
	}
	db.Create(&adminPermissionGroup)
	guestPermissionGroup := models.UserPermissionGroup{
		GroupName:         "guest",
		PermissionStrings: "users.list;users.detail",
	}
	db.Create(&guestPermissionGroup)
	return nil
}
