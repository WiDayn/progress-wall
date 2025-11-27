package services

import (
	"errors"
	"fmt"

	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// ProjectService provides methods for managing projects.
type ProjectService struct {
	db *gorm.DB
}

// NewProjectService creates a new instance of ProjectService.
func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{
		db: db,
	}
}

// GetProjectByID retrieves a project by its ID, including its owner, boards, and members.
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

// GetProjectsByUserID retrieves all projects owned by a specific user.
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

// CreateProject creates a new project under a team and assigns the creator as ProjectAdmin.
// It executes the creation and member assignment within a transaction.
func (s *ProjectService) CreateProject(project *models.Project) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Assign the creator as ProjectAdmin in ProjectMember table
	member := &models.ProjectMember{
		ProjectID: project.ID,
		UserID:    project.OwnerID,
		Role:      models.ProjectRoleAdmin, // Important: Creator gets Admin privileges
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// GetTeamProjects retrieves all projects associated with a specific team.
func (s *ProjectService) GetTeamProjects(teamID uint) ([]models.Project, error) {
	var projects []models.Project
	err := s.db.Where("team_id = ?", teamID).Find(&projects).Error
	return projects, err
}

// UpdateProject updates the fields of a project identified by projectID.
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

// DeleteProject soft-deletes a project by its ID.
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
