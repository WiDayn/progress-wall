package services

import (
	"errors"
	"progress-wall-backend/models"

	"gorm.io/gorm"
)

type TeamService struct {
	db *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{db: db}
}

// Creates a new team and assigns the creator as TeamAdmin.
func (s *TeamService) CreateTeam(name, description string, creatorID uint) (*models.Team, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the Team
	team := &models.Team{
		Name:        name,
		Description: description,
		CreatorID:   creatorID,
	}
	if err := tx.Create(team).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Add Creator as TeamAdmin
	member := &models.TeamMember{
		TeamID: team.ID,
		UserID: creatorID,
		Role:   models.TeamRoleAdmin,
	}
	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return team, nil
}

// Adds a user to the team.
func (s *TeamService) AddMember(teamID, userID uint, role models.TeamRole) error {
	// Check if user exists
	var count int64
	s.db.Model(&models.User{}).Where("id = ?", userID).Count(&count)
	if count == 0 {
		return errors.New("user not found")
	}

	// Check if already a member
	var member models.TeamMember
	err := s.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member).Error
	if err == nil {
		return errors.New("user is already a member of this team")
	}

	// Add member
	newMember := models.TeamMember{
		TeamID: teamID,
		UserID: userID,
		Role:   role,
	}
	return s.db.Create(&newMember).Error
}

// Gets all teams a user belongs to.
func (s *TeamService) GetUserTeams(userID uint) ([]models.Team, error) {
	var teams []models.Team
	// Select teams where user_id in team_members matches
	err := s.db.Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userID).
		Find(&teams).Error

	if err != nil {
		return nil, err
	}
	return teams, err
}

// Gets team members for a specific team
func (s *TeamService) GetTeamMembers(teamID uint) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := s.db.Where("team_id = ?", teamID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, err
}
