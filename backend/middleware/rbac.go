package middleware

import (
	"net/http"
	"progress-wall-backend/models"
	"progress-wall-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Role-based access control
type RBACMiddleware struct {
	permService *services.PermissionService
	db          *gorm.DB
}

func NewRBACMiddleware(permService *services.PermissionService, db *gorm.DB) *RBACMiddleware {
	return &RBACMiddleware{
		permService: permService,
		db:          db,
	}
}

// Ensures the user is a System Administrator.
func (m *RBACMiddleware) RequireSysAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		
		isAdmin, err := m.permService.IsSysAdmin(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify system admin privileges"})
			c.Abort()
			return
		}

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "System Admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Checks if the user has access to the project.
// level: "view" (access) or "manage" (admin).
// paramKey: the name of the URL parameter containing the ID (e.g., "projectId" or "boardId").
// idType: "project" (direct project ID) or "board" (board ID, needs resolution to project).
func (m *RBACMiddleware) RequireProjectAccess(level string, paramKey string, idType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idStr := c.Param(paramKey)
		
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			c.Abort()
			return
		}
		resourceID := uint(id)

		// Resolve ProjectID
		var projectID uint

		switch idType {
		case "project":
			projectID = resourceID

		case "board":
			var board models.Board
			if err := m.db.Select("project_id").First(&board, resourceID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
				c.Abort()
				return
			}
			projectID = board.ProjectID

		case "column":
			var column models.Column
			if err := m.db.Preload("Board").First(&column, resourceID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Column not found"})
				c.Abort()
				return
			}
			projectID = column.Board.ProjectID

		case "task":
			var task models.Task
			if err := m.db.Select("project_id").First(&task, resourceID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
				c.Abort()
				return
			}
			projectID = task.ProjectID
		
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID type for RBAC"})
			c.Abort()
			return
		}

		// Permission check
		var allowed bool
		var checkErr error

		if level == "manage" {
			allowed, checkErr = m.permService.CanManageProject(userID, projectID)
		} else {
			allowed, checkErr = m.permService.CanAccessProject(userID, projectID)
		}

		if checkErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Checks if the user has access to the team.
// level: "view" or "manage".
func (m *RBACMiddleware) RequireTeamAccess(level string, paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		teamIDStr := c.Param(paramKey)
		
		teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID"})
			c.Abort()
			return
		}

		var allowed bool
		var checkErr error

		if level == "manage" {
			allowed, checkErr = m.permService.CanManageTeam(userID, uint(teamID))
		} else {
			allowed, checkErr = m.permService.CanAccessTeam(userID, uint(teamID))
		}

		if checkErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Insufficient team permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}