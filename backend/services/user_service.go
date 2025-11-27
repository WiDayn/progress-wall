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

// UpdateUser 更新用户信息
// 参数 user: 包含更新信息的User对象，必须包含ID
// 返回 error: 更新失败时返回错误
func (s *UserService) UpdateUser(user *models.User) error {
	if user == nil {
		return errors.New("用户信息不能为空")
	}
	if user.ID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	result := s.db.Save(user)
	if result.Error != nil {
		return errors.New("更新用户失败")
	}
	return nil
}