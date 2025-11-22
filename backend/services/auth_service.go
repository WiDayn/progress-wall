package services

import (
	"errors"
	"time"

	"progress-wall-backend/config"
	"progress-wall-backend/models"
	"progress-wall-backend/utils"

	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string
	Password string
}

// LoginResult 登录结果
type LoginResult struct {
	AccessToken string
	User        *models.User
}

// Login 处理登录业务逻辑
func (s *AuthService) Login(req LoginRequest) (*LoginResult, error) {
	// 查找用户（支持用户名或邮箱登录）
	var user models.User
	result := s.db.Where("username = ? OR email = ?", req.Username, req.Username).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, ErrUserNotFound
	}

	// 检查用户状态
	if user.Status != models.UserStatusEnabled {
		return nil, ErrUserDisabled
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, s.cfg)
	if err != nil {
		return nil, ErrGenerateToken
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	if err := s.db.Save(&user).Error; err != nil {
		return nil, ErrUpdateLoginTime
	}

	// 清除敏感信息
	user.Password = ""

	return &LoginResult{
		AccessToken: token,
		User:        &user,
	}, nil
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string
	Email    string
	Password string
	Nickname string
}

// RegisterResult 注册结果
type RegisterResult struct {
	User *models.User
}

// Register 处理注册业务逻辑
func (s *AuthService) Register(req RegisterRequest) (*RegisterResult, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	result := s.db.Where("username = ?", req.Username).First(&existingUser)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("查询用户失败")
	}
	if result.Error == nil {
		return nil, ErrUserExists
	}

	// 检查邮箱是否已存在
	result = s.db.Where("email = ?", req.Email).First(&existingUser)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("查询用户失败")
	}
	if result.Error == nil {
		return nil, ErrUserExists
	}

	// 验证密码强度
	if err := utils.ValidatePasswordStrength(req.Password); err != nil {
		return nil, ErrInvalidPassword
	}

	// 验证邮箱格式
	if !utils.ValidateEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 验证用户名格式（3-20个字符，只能包含字母、数字、下划线）
	if len(req.Username) < 3 || len(req.Username) > 20 {
		return nil, errors.New("用户名长度必须在3-20个字符之间")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("加密密码失败")
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   models.UserStatusEnabled,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 清除敏感信息
	user.Password = ""

	return &RegisterResult{
		User: &user,
	}, nil
}
