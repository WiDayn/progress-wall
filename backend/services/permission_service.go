package services

import (
    "errors"
    "progress-wall-backend/models"

    "gorm.io/gorm"
)

// PermissionService handles permission validation logic.
type PermissionService struct {
    db *gorm.DB
}

// NewPermissionService creates a new instance of PermissionService.
func NewPermissionService(db *gorm.DB) *PermissionService {
    return &PermissionService{db: db}
}

// Checks if the user is a system administrator
func (s *PermissionService) IsSysAdmin(userID uint) (bool, error) {
    var user models.User
    if err := s.db.Select("system_role").First(&user, userID).Error; err != nil {
        return false, err
    }
    return user.SystemRole == models.SystemRoleAdmin, nil
}

// Checks if the user is an admin of a specific team (Internal Use)
func (s *PermissionService) isTeamAdmin(userID, teamID uint) (bool, error) {
    var member models.TeamMember
    err := s.db.Where("user_id = ? AND team_id = ?", userID, teamID).First(&member).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }
    return member.Role == models.TeamRoleAdmin, nil
}

// Checks if the user is an admin of a specific project (Internal Use)
func (s *PermissionService) isProjectAdmin(userID, projectID uint) (bool, error) {
    var member models.ProjectMember
    err := s.db.Where("user_id = ? AND project_id = ?", userID, projectID).First(&member).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }
    return member.Role == models.ProjectRoleAdmin, nil
}

// Checks if the user has permission to manage a team (e.g., edit info, add members)
// SysAdmin OR TeamAdmin of that team can manage the team
func (s *PermissionService) CanManageTeam(userID, teamID uint) (bool, error) {
    if isAdmin, err := s.IsSysAdmin(userID); err != nil {
        return false, err
    } else if isAdmin {
        return true, nil
    }

    return s.isTeamAdmin(userID, teamID)
}

// Checks if the user can access a team (e.g., view team content)
// SysAdmin OR any member of that team (Member/Admin) can access the team
func (s *PermissionService) CanAccessTeam(userID, teamID uint) (bool, error) {
    if isAdmin, err := s.IsSysAdmin(userID); err != nil {
        return false, err
    } else if isAdmin {
        return true, nil
    }

    var count int64
    err := s.db.Model(&models.TeamMember{}).
        Where("user_id = ? AND team_id = ?", userID, teamID).
        Count(&count).Error
    return count > 0, err
}

// Checks if the user has permission to manage a project.
// SysAdmin OR TeamAdmin of the parent team OR ProjectAdmin of that project.
func (s *PermissionService) CanManageProject(userID, projectID uint) (bool, error) {
    // Check SysAdmin
    if isAdmin, err := s.IsSysAdmin(userID); err != nil {
        return false, err
    } else if isAdmin {
        return true, nil
    }

    // Retrieve the project to identify its parent TeamID
    var project models.Project
    if err := s.db.Select("team_id").First(&project, projectID).Error; err != nil {
        return false, err
    }

    // Check TeamAdmin of the parent team
    if isTeamAdmin, err := s.isTeamAdmin(userID, project.TeamID); err != nil {
        return false, err
    } else if isTeamAdmin {
        return true, nil
    }

    return s.isProjectAdmin(userID, projectID)
}

// Checks if the user can access a project.
// SysAdmin OR TeamAdmin OR Project Member
func (s *PermissionService) CanAccessProject(userID, projectID uint) (bool, error) {
    if canManage, err := s.CanManageProject(userID, projectID); err != nil {
        return false, err
    } else if canManage {
        return true, nil
    }

    var count int64
    err := s.db.Model(&models.ProjectMember{}).
        Where("user_id = ? AND project_id = ?", userID, projectID).
        Count(&count).Error
    return count > 0, err
}