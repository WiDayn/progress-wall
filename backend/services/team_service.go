package services

import (
	"errors"
	"fmt"

	"progress-wall-backend/models"

	"gorm.io/gorm"
)

var (
	ErrTeamNotFound      = errors.New("team not found")
	ErrUserAlreadyMember = errors.New("user is already a member of the team")
)

type TeamService struct {
	db *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{db: db}
}

// CreateTeam 创建团队并自动将创建者设为管理员
func (s *TeamService) CreateTeam(name, description string, creatorID uint) (*models.Team, error) {
	team := &models.Team{
		Name:        name,
		Description: description,
		CreatorID:   creatorID,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}

		// Add creator as admin
		member := &models.TeamMember{
			TeamID: team.ID,
			UserID: creatorID,
			Role:   models.TeamRoleAdmin,
		}
		if err := tx.Create(member).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %v", err)
	}

	return team, nil
}

// GetUserTeams 获取用户参与的所有团队
func (s *TeamService) GetUserTeams(userID uint) ([]models.Team, error) {
	var teams []models.Team
	// 通过 team_members 关联查询
	err := s.db.Model(&models.Team{}).
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userID).
		Preload("Creator").
		Find(&teams).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user teams: %v", err)
	}

	return teams, nil
}

// GetTeamByID 获取团队详情
func (s *TeamService) GetTeamByID(teamID uint) (*models.Team, error) {
	var team models.Team
	err := s.db.Preload("Creator").Preload("Members").First(&team, teamID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return &team, nil
}

// GetMembers 获取团队成员列表
func (s *TeamService) GetMembers(teamID uint) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := s.db.Where("team_id = ?", teamID).Preload("User").Find(&members).Error
	return members, err
}

// AddTeamMember 添加成员到团队
func (s *TeamService) AddTeamMember(teamID, userID uint, role models.TeamRole) error {
	// Check if already exists
	var count int64
	s.db.Model(&models.TeamMember{}).Where("team_id = ? AND user_id = ?", teamID, userID).Count(&count)
	if count > 0 {
		return ErrUserAlreadyMember
	}

	member := models.TeamMember{
		TeamID: teamID,
		UserID: userID,
		Role:   role,
	}

	return s.db.Create(&member).Error
}
