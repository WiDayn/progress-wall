package auth

import (
	"net/http"
	"progress-wall-backend/config"
	"progress-wall-backend/models"
	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterHandler 注册处理器
type RegisterHandler struct {
	authService *services.AuthService
}

// NewRegisterHandler 创建注册处理器
func NewRegisterHandler(db *gorm.DB, cfg *config.Config) *RegisterHandler {
	return &RegisterHandler{
		authService: services.NewAuthService(db, cfg),
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// RegisterResponse 注册响应结构
type RegisterResponse struct {
	User *models.User `json:"user"`
}

// Register 处理注册请求
func (h *RegisterHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 调用service层处理业务逻辑
	result, err := h.authService.Register(services.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	})

	if err != nil {
		// 根据错误类型返回相应的HTTP状态码
		switch err {
		case services.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case services.ErrInvalidPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, RegisterResponse{
		User: result.User,
	})
}
