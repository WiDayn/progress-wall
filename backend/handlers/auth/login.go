package auth

import (
	"net/http"
	"progress-wall-backend/config"
	"progress-wall-backend/models"
	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginHandler 登录处理器
type LoginHandler struct {
	authService *services.AuthService
}

// NewLoginHandler 创建登录处理器
func NewLoginHandler(db *gorm.DB, cfg *config.Config) *LoginHandler {
	return &LoginHandler{
		authService: services.NewAuthService(db, cfg),
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	AccessToken string       `json:"accessToken"`
	User        *models.User `json:"user,omitempty"`
}

// Login 处理登录请求
func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 调用service层处理业务逻辑
	result, err := h.authService.Login(services.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		// 根据错误类型返回相应的HTTP状态码
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case services.ErrUserDisabled:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: result.AccessToken,
		User:        result.User,
	})
}
