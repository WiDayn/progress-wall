package project

import (
	"net/http"
	"strconv"
	"time"

	"progress-wall-backend/models"
	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	projectService *services.ProjectService
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{
		projectService: services.NewProjectService(db),
	}
}

// GetProject 获取单个项目
// GET /api/projects/:projectId
func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	project, err := h.projectService.GetProjectByID(uint(projectID))
	if err != nil {
		if err == services.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// GetProjects 获取用户的所有项目 (Deprecated or kept for "All Projects" view)
// GET /api/projects
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	projects, err := h.projectService.GetProjectsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// CreateProject 创建项目
// POST /api/teams/:teamId/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	var createProjectRequest struct {
		Name        string     `json:"name" binding:"required,min=1,max=100"`
		Description string     `json:"description" binding:"max=500"`
		Status      *int       `json:"status"`
		StartDate   *time.Time `json:"start_date"`
		EndDate     *time.Time `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&createProjectRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// Validate dates
	if createProjectRequest.StartDate != nil && createProjectRequest.EndDate != nil {
		if createProjectRequest.EndDate.Before(*createProjectRequest.StartDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能早于开始时间"})
			return
		}
	}

	status := models.ProjectStatusActive
	if createProjectRequest.Status != nil {
		status = models.ProjectStatus(*createProjectRequest.Status)
	}

	project := &models.Project{
		Name:        createProjectRequest.Name,
		Description: createProjectRequest.Description,
		Status:      status,
		StartDate:   createProjectRequest.StartDate,
		EndDate:     createProjectRequest.EndDate,
		OwnerID:     userID,
		TeamID:      uint(teamID),
	}

	if err := h.projectService.CreateProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetTeamProjects Gets all projects for a specific team.
// GET /api/teams/:teamId/projects
func (h *ProjectHandler) GetTeamProjects(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
		return
	}

	projects, err := h.projectService.GetTeamProjects(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// UpdateProject 更新项目
// PUT /api/projects/:projectId
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	var updateProjectRequest struct {
		Name        *string               `json:"name" binding:"omitempty,min=1,max=100"`
		Description *string               `json:"description" binding:"omitempty,max=500"`
		Status      *models.ProjectStatus `json:"status"`
		StartDate   *time.Time            `json:"start_date"`
		EndDate     *time.Time            `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&updateProjectRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if updateProjectRequest.Name != nil {
		updates["name"] = *updateProjectRequest.Name
	}
	if updateProjectRequest.Description != nil {
		updates["description"] = *updateProjectRequest.Description
	}
	if updateProjectRequest.Status != nil {
		updates["status"] = *updateProjectRequest.Status
	}
	if updateProjectRequest.StartDate != nil {
		updates["start_date"] = *updateProjectRequest.StartDate
	}
	if updateProjectRequest.EndDate != nil {
		updates["end_date"] = *updateProjectRequest.EndDate
	}

	if err := h.projectService.UpdateProject(uint(projectID), updates); err != nil {
		if err == services.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteProject 删除项目
// DELETE /api/projects/:projectId
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	if err := h.projectService.DeleteProject(uint(projectID)); err != nil {
		if err == services.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
