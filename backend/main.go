package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"progress-wall-backend/config"
	"progress-wall-backend/database"
	"progress-wall-backend/routes"
	"progress-wall-backend/services"

	"github.com/robfig/cron/v3"
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

	// 初始化并启动定时任务调度器（核心新增逻辑）
	var cronInstance *cron.Cron // 声明定时任务实例
	// 创建调度器实例（传入数据库连接和通知服务URL）(若是后端有了NotificationService，则将URL改成cfg.NotificationService.URL)
	schedulerIns := services.NewScheduler(db, "http://localhost:8080")
	// 启动定时任务，返回cron实例用于后续关闭
	cronInstance = schedulerIns.Start()
	defer cronInstance.Stop() // 程序退出时停止定时任务
	log.Println("定时任务调度器已启动")

	// 启动HTTP服务器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("服务器启动在端口 %s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

	// 等待中断信号（优雅退出逻辑）
	quit := make(chan os.Signal, 1)
	// 监听Ctrl+C和kill信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞等待信号
	log.Println("开始优雅关闭服务...")

	// 若定时任务已启动，等待其完全停止
	if cronInstance != nil {
		cronInstance.Stop()
		log.Println("定时任务已停止")
	}

	log.Println("服务已正常退出")
}
