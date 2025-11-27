package notification

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"progress-wall-backend/services"
)

// NotificationHandler 通知服务处理器
type NotificationHandler struct {
	notificationService *services.TaskNotificationService // 若有服务层可关联，无则留空
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	return &NotificationHandler{
		notificationService: services.NewTaskNotificationService(db), // 若无需服务层可简化为 nil
	}
}

// 通知请求体结构（与定时任务格式匹配）
type TaskNotificationReq struct {
	UserID           uint   `json:"user_id"`           // 任务创建者ID
	TaskID           uint   `json:"task_id"`           // 任务ID
	TaskTitle        string `json:"task_title"`        // 任务标题
	NotificationType string `json:"notification_type"` // 通知类型
}

// ReceiveTaskNotification 接收定时任务发送的通知
// 对应路由：POST /api/notifications
func (h *NotificationHandler) ReceiveTaskNotification(c *gin.Context) {
	var req TaskNotificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的通知参数",
			"error":   err.Error(),
		})
		return
	}

	// 终端打印定时任务信息（核心逻辑）
	h.printNotification(req)

	// 若有业务逻辑可调用服务层处理（如存储通知记录）
	// err := h.notificationService.SaveNotification(req)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"message": "保存通知失败"})
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "通知接收成功",
	})
}

// 终端打印通知信息
func (h *NotificationHandler) printNotification(req TaskNotificationReq) {
	fmt.Println("======================================")
	fmt.Printf("【定时任务通知】\n")
	fmt.Printf("时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("用户ID：%d\n", req.UserID)
	fmt.Printf("任务ID：%d\n", req.TaskID)
	fmt.Printf("任务标题：%s\n", req.TaskTitle)
	fmt.Printf("类型：%s\n", req.NotificationType)
	fmt.Println("======================================\n")
}
