package user

import (
	"net/http"

	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProfileHandler 用户信息处理器
type ProfileHandler struct {
	userService *services.UserService
}

// NewProfileHandler 创建用户信息处理器
func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{
		userService: services.NewUserService(db),
	}
}

// GetProfile 获取当前登录用户信息
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// 从中间件设置的上下文中获取用户ID
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	// 调用service层获取用户信息
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
