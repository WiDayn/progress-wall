package router

import (
	"progress-wall-backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HandlerDependencies struct {
	UserService services.UserService
}

func NewRouter(deps HandlerDependencies) *gin.Engine {
	router := gin.Default()

	// CORS 配置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 用户相关路由组
	userHandler := NewUserHandler(deps.UserService)
	userGroup := router.Group("/api/auth")
	{
		userGroup.POST("/register", userHandler.Register) // 注册
		userGroup.POST("/login", userHandler.Login)       // 登录
	}

	return router
}
