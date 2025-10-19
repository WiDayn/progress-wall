package project

import (
	"net/http"
	"strconv"

	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{
		projectService: services.NewProjectService(db),
	}
}

// POST /api/teams/:teamId/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID := c.GetUint("user_id")  // Set by AuthMiddleware

	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)  // Should already validated by RBAC Middleware
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	project, err := h.projectService.CreateProject(req.Name, req.Description, uint(teamID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// Gets all projects for a specific team.
// GET /api/teams/:teamId/projects
func (h *ProjectHandler) GetTeamProjects(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, _ := strconv.ParseUint(teamIDStr, 10, 32)  // Already validated by RBAC Middleware

	projects, err := h.projectService.GetTeamProjects(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}
