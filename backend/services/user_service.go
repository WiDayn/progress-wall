package services

import (
	"errors"
	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, errors.New("查询用户失败")
	}

	// 清除敏感信息
	user.Password = ""

	return &user, nil
}
