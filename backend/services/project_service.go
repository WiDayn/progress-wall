package services

import (
	"errors"
	"fmt"

	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// ProjectService 项目服务
type ProjectService struct {
	db *gorm.DB
}

// NewProjectService 创建项目服务
func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{
		db: db,
	}
}

// GetProjectByID 根据ID获取项目
func (s *ProjectService) GetProjectByID(projectID uint) (*models.Project, error) {
	var project models.Project
	result := s.db.
		Preload("Owner").
		Preload("Boards").
		Preload("Members").
		First(&project, projectID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("查询项目失败: %v", result.Error)
	}

	return &project, nil
}

// GetProjectsByUserID 获取用户的所有项目
func (s *ProjectService) GetProjectsByUserID(userID uint) ([]models.Project, error) {
	var projects []models.Project
	result := s.db.
		Where("owner_id = ?", userID).
		Preload("Owner").
		Order("created_at DESC").
		Find(&projects)

	if result.Error != nil {
		return nil, errors.New("查询项目列表失败")
	}

	return projects, nil
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(project *models.Project) error {
	if err := s.db.Create(project).Error; err != nil {
		return fmt.Errorf("创建项目失败: %v", err)
	}
	return nil
}

// UpdateProject 更新项目
func (s *ProjectService) UpdateProject(projectID uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Project{}).Where("id = ?", projectID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新项目失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}
	return nil
}

// DeleteProject 删除项目（软删除）
func (s *ProjectService) DeleteProject(projectID uint) error {
	result := s.db.Delete(&models.Project{}, projectID)
	if result.Error != nil {
		return fmt.Errorf("删除项目失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}
	return nil
}
