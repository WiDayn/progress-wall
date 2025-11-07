package services

import (
	"errors"
	"progress-wall-backend/config"
	"progress-wall-backend/models"
	"progress-wall-backend/repository"
	"progress-wall-backend/utils"
)

type UserService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register 注册新用户
func (s *userService) Register(username, email, password string) (*models.User, error) {
	// 检查用户名或邮箱是否已存在
	if user, err := s.userRepo.FindByUsername(username); err == nil {
		return user, errors.New("username already exists")
	}
	if user, err := s.userRepo.FindByEmail(email); err == nil {
		return user, errors.New("email already exists")
	}

	// 创建用户实例
	user := &models.User{
		Username: username,
		Email:    email,
	}

	// 加密密码
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// 保存用户
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *userService) Login(email, password string) (string, *models.User, error) {
	// 查找用户
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("user not found")
	}

	// 检查密码
	if !user.CheckPassword(password) {
		return "", nil, errors.New("invalid password")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, s.cfg)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	// 不返回密码
	user.Password = ""

	return token, user, nil
}
