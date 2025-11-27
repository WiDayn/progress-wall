package team

import (
	"net/http"
	"strconv"

	"progress-wall-backend/models"
	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(db *gorm.DB) *TeamHandler {
	return &TeamHandler{
		teamService: services.NewTeamService(db),
	}
}

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateTeam 创建团队
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	team, err := h.teamService.CreateTeam(req.Name, req.Description, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// GetMyTeams 获取当前用户的团队列表
func (h *TeamHandler) GetMyTeams(c *gin.Context) {
	userID := c.GetUint("user_id")
	teams, err := h.teamService.GetUserTeams(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// GetTeam 获取单个团队详情
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	team, err := h.teamService.GetTeamByID(uint(teamID))
	if err != nil {
		if err == services.ErrTeamNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// GetTeamMembers 获取团队成员
func (h *TeamHandler) GetTeamMembers(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	members, err := h.teamService.GetMembers(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// AddMember 添加成员
func (h *TeamHandler) AddMember(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req struct {
		UserID uint            `json:"user_id" binding:"required"`
		Role   models.TeamRole `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == 0 {
		req.Role = models.TeamRoleMember // Default to member
	}

	if err := h.teamService.AddMember(uint(teamID), req.UserID, req.Role); err != nil {
		if err == services.ErrUserAlreadyMember {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}
