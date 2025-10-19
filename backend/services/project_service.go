package services

import (
	"progress-wall-backend/models"

	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

// Creates a new project under a team and assigns the creator as ProjectAdmin.
func (s *ProjectService) CreateProject(name, description string, teamID, creatorID uint) (*models.Project, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create the Project record
	project := &models.Project{
		Name:        name,
		Description: description,
		TeamID:      teamID,
		OwnerID:     creatorID,
		Status:      models.ProjectStatusActive,
	}

	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Assign the creator as ProjectAdmin in ProjectMember table
	member := &models.ProjectMember{
		ProjectID: project.ID,
		UserID:    creatorID,
		Role:      models.ProjectRoleAdmin, // Important: Creator gets Admin privileges
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return project, nil
}

// Gets all projects for a specific team.
func (s *ProjectService) GetTeamProjects(teamID uint) ([]models.Project, error) {
	var projects []models.Project
	err := s.db.Where("team_id = ?", teamID).Find(&projects).Error
	return projects, err
}