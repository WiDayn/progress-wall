package user

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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

// UploadAvatar 上传头像
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	// 检查文件大小 (例如限制为 2MB)
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过2MB"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 JPG, PNG, GIF 格式的图片"})
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	// 确保目录存在
	savePath := filepath.Join("uploads", "avatars", filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 生成访问 URL
	// 注意：这里假设静态资源通过 /uploads 路径访问
	avatarURL := "/uploads/avatars/" + filename

	// 更新用户信息
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	user.Avatar = avatarURL
	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新头像信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"url":     avatarURL,
	})
}

// UpdateProfileRequest 更新用户信息请求结构
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// UpdateProfile 更新用户信息
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	// Email 更新可能需要验证唯一性等逻辑，这里暂时简化
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
