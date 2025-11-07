package database

import (
	"fmt"
	"log"
	"progress-wall-backend/config"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB   *gorm.DB
	once sync.Once
)

// InitDB initializes the database connection based on configuration
func InitDB(cfg *config.Config) error {
	var err error

	once.Do(func() {
		// 构造 MySQL DSN
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Name,
		)

		// 初始化 GORM 配置
		gormConfig := &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
		if cfg.Server.Mode == "release" {
			gormConfig.Logger = logger.Default.LogMode(logger.Silent)
		}

		// 连接数据库
		DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %v", err)
			return
		}

		// 配置连接池
		sqlDB, err := DB.DB()
		if err != nil {
			err = fmt.Errorf("failed to get underlying sql.DB: %v", err)
			return
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		log.Printf("Database connected successfully (Type: mysql)")
	})

	return err
}

// GetDB returns the initialized *gorm.DB
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitDB first.")
	}
	return DB
}

// CloseDB closes the underlying database connection
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
