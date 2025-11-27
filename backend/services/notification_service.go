package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TaskNotificationService定时任务服务
type TaskNotificationService struct {
	UserID           uint   `json:"user_id"`           // 对应定时任务的UserID
	TaskID           uint   `json:"task_id"`           // 对应定时任务的TaskID
	TaskTitle        string `json:"task_title"`        // 对应定时任务的TaskTitle
	NotificationType string `json:"notification_type"` // 通知类型
	db               *gorm.DB
}

func NewTaskNotificationService(db *gorm.DB) *TaskNotificationService {
	return &TaskNotificationService{
		db: db,
	}
}

// 可直接绑定到 Gin 路由
func HandleTaskNotification(c *gin.Context) {
	var req TaskNotificationService
	// 解析定时任务发送的 JSON 请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数解析失败",
			"error":   err.Error(),
		})
		return
	}

	// 核心逻辑：在终端打印定时任务筛选出的任务信息
	printTaskNotification(req)

	// 返回成功响应，让定时任务确认通知发送成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "通知接收成功",
		"data":    req,
	})
}

// printTaskNotification 封装终端打印逻辑（可复用）
func printTaskNotification(req TaskNotificationService) {
	fmt.Println("======================================")
	fmt.Printf("【定时任务自动触发通知】\n")
	fmt.Printf("接收时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("任务创建者ID：%d\n", req.UserID)
	fmt.Printf("任务ID：%d\n", req.TaskID)
	fmt.Printf("任务标题：%s\n", req.TaskTitle)
	fmt.Printf("通知类型：%s\n", req.NotificationType)
	fmt.Println("======================================\n")
}
