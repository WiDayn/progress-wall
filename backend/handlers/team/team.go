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

// POST /api/teams
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userID := c.GetUint("user_id") // Get from AuthMiddleware

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	team, err := h.teamService.CreateTeam(req.Name, req.Description, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// Gets all teams a user belongs to.
// GET /api/teams
func (h *TeamHandler) GetMyTeams(c *gin.Context) {
	userID := c.GetUint("user_id")

	teams, err := h.teamService.GetUserTeams(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// POST /api/teams/:teamId/members
func (h *TeamHandler) AddMember(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	// Who to add, and what role
	var req struct {
		UserID uint            `json:"user_id" binding:"required"`
		Role   models.TeamRole `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Default to Member if role is invalid
	if req.Role != models.TeamRoleAdmin {
		req.Role = models.TeamRoleMember
	}

	if err := h.teamService.AddMember(uint(teamID), req.UserID, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

// GET /api/teams/:teamId/members
func (h *TeamHandler) GetTeamMembers(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	members, err := h.teamService.GetTeamMembers(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team members: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}
