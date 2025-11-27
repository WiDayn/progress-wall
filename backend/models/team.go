package models

import (
    "time"
    "gorm.io/gorm"
)

type TeamRole int

const (
    TeamRoleMember TeamRole = 1
    TeamRoleAdmin TeamRole = 2
)

type Team struct {
    ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    Name        string         `json:"name" gorm:"size:100;not null;comment:'Team name'"`
    Description string         `json:"description" gorm:"type:text"`
    CreatorID   uint           `json:"creator_id" gorm:"not null;index"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    // associations
    Creator  *User     `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
    Projects []Project `json:"projects,omitempty" gorm:"foreignKey:TeamID"` // A team can have multiple projects
    Members  []User    `json:"members,omitempty" gorm:"many2many:team_members;"`
}

type TeamMember struct {
    ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    TeamID    uint           `json:"team_id" gorm:"not null;index;uniqueIndex:uk_team_user"`
    UserID    uint           `json:"user_id" gorm:"not null;index;uniqueIndex:uk_team_user"`
    Role      TeamRole       `json:"role" gorm:"type:tinyint;default:1;comment:'Team role: 1=member, 2=admin'"`
    JoinedAt  time.Time      `json:"joined_at"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    Team Team `json:"team" gorm:"foreignKey:TeamID"`
    User User `json:"user" gorm:"foreignKey:UserID"`
}
