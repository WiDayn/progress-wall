package main

import (
	"log"
	"progress-wall-backend/config"
	"progress-wall-backend/database"
	"progress-wall-backend/routes"
)

func main() {

	cfg := config.Load()

	// 打印数据库配置，检查用户名/密码是否正确
	fmt.Printf("DB_USER=%s, DB_PASSWORD=%s, DB_HOST=%s\n", cfg.DB.User, cfg.DB.Password, cfg.DB.Host)

	// 下面是你原来的 InitDB 调用
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 获取数据库实例
	db := database.GetDB()

	// 自动迁移数据库表
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化依赖注入层
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg)

	// 初始化路由
	deps := router.HandlerDependencies{
		UserService: userService,
	}
	r := router.NewRouter(deps)

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 执行数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 初始化基础数据
	if err := database.Seed(db); err != nil {
		log.Fatal("初始化基础数据失败:", err)
	}

	log.Println("数据库初始化完成")

	// 设置路由
	r := routes.SetupRoutes(db, cfg)

	// 启动HTTP服务器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("服务器启动在端口 %s\n", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
