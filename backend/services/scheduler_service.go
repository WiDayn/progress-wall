package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"progress-wall-backend/models"

	"gorm.io/gorm"

	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	db                  *gorm.DB // SQLite数据库连接
	notificationBaseURL string   // 通知服务(B8)基础URL
}

// NewScheduler 创建调度器实例
func NewScheduler(db *gorm.DB, notificationBaseURL string) *Scheduler {
	return &Scheduler{
		db:                  db,
		notificationBaseURL: notificationBaseURL,
	}
}

// Start 启动定时任务（每小时执行一次）
func (s *Scheduler) Start() *cron.Cron {
	c := cron.New()

	// 注册定时任务：每小时第0分钟执行（与需求中的"0 * * * *"一致）
	_, err := c.AddFunc("1 * * * *", s.Execute)
	if err != nil {
		log.Fatalf("注册定时任务失败: %v", err)
	}

	c.Start()
	log.Println("定时任务调度器已启动，每小时执行一次")
	return c
}

// Execute 执行任务扫描和通知逻辑（核心方法）
func (s *Scheduler) Execute() {
	log.Println("=== 开始执行任务扫描 ===")

	// 1. 查询符合条件的任务
	tasks, err := s.queryPendingTasks()
	if err != nil {
		log.Printf("查询任务失败: %v", err)
		return
	}

	log.Printf("发现 %d 个符合条件的任务", len(tasks))
	if len(tasks) == 0 {
		log.Println("=== 任务扫描完成，无需要处理的任务 ===")
		return
	}

	// 2. 处理每个任务：发送通知并更新状态
	for _, task := range tasks {
		log.Printf("处理任务: %s (ID: %d, 截止时间: %v)", task.Title, task.ID, task.DueDate)

		// 发送通知
		if err := s.sendNotification(task); err != nil {
			log.Printf("任务 %d 通知发送失败: %v", task.ID, err)
			continue
		}

		// 更新提醒状态（确保幂等性）
		if err := s.updateAlertStatus(task.ID); err != nil {
			log.Printf("任务 %d 状态更新失败: %v", task.ID, err)
		}
	}

	log.Println("=== 任务扫描完成 ===")
}

// queryPendingTasks 查询符合条件的任务（未完成、24小时内截止、未发送提醒）
func (s *Scheduler) queryPendingTasks() ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()
	// 计算24小时后的时间点（SQLite的时间比较需要使用UTC或统一时区）
	deadlineEnd := now.Add(24 * time.Hour)

	// SQL条件：状态不是已完成/已归档，deadline在未来24小时内，且未发送提醒
	err := s.db.Where(
		"status NOT IN (?) AND due_date IS NOT NULL AND due_date > ? AND due_date < ? AND deadline_alert_sent != ?",
		[]int{3, 5}, now, deadlineEnd, 1,
	).Find(&tasks).Error

	return tasks, err
}

// sendNotification 调用B8通知服务发送提醒
func (s *Scheduler) sendNotification(task models.Task) error {
	// 构建通知 payload（按需求包含必要字段）
	payload := struct {
		UserID           uint   `json:"user_id"`
		TaskID           uint   `json:"task_id"`
		TaskTitle        string `json:"task_title"`
		NotificationType string `json:"notification_type"`
	}{
		UserID:           task.CreatorID,
		TaskID:           task.ID,
		TaskTitle:        task.Title,
		NotificationType: "TASK_DEADLINE_APPROACHING", // 固定通知类型
	}

	// 序列化payload为JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化通知数据失败: %v", err)
	}

	// 调用通知服务API
	resp, err := http.Post(
		fmt.Sprintf("%s/api/notifications", s.notificationBaseURL),
		"application/json",
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		return fmt.Errorf("调用通知服务失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("通知服务返回错误状态码: %d", resp.StatusCode)
	}

	return nil
}

// updateAlertStatus 更新任务的提醒状态（避免重复发送）
func (s *Scheduler) updateAlertStatus(taskID uint) error {
	return s.db.Model(&models.Task{}).
		Where("id = ?", taskID).
		Update("deadline_alert_sent", 1).Error
}
