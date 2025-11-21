package main

import (
	"fmt"
	"log"

	"progress-wall-backend/config"
	"progress-wall-backend/database"
	"progress-wall-backend/routes"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 打印配置信息（不打印敏感信息）
	log.Printf("服务器配置: 端口=%s, 模式=%s, 数据库类型=%s",
		cfg.Server.Port, cfg.Server.Mode, cfg.DB.Type)

	// 初始化数据库连接
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 执行数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化基础数据
	if err := database.Seed(db); err != nil {
		log.Fatalf("初始化基础数据失败: %v", err)
	}

	log.Println("数据库初始化完成")

	// 设置路由
	r := routes.SetupRoutes(db, cfg)

	// 启动HTTP服务器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("服务器启动在端口 %s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
