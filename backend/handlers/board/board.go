package board

import (
	"net/http"
	"strconv"

	"progress-wall-backend/models"
	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BoardHandler 看板处理器
type BoardHandler struct {
	boardService *services.BoardService
}

// NewBoardHandler 创建看板处理器
func NewBoardHandler(db *gorm.DB) *BoardHandler {
	return &BoardHandler {
		boardService: services.NewBoardService(db),
	}
}

// GetBoard 获取单个看板（包含嵌套的列和任务）
func (h *BoardHandler) GetBoard(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	board, err := h.boardService.GetBoardByID(uint(boardID))
	if err != nil {
		if err == services.ErrBoardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}

// GetBoards 获取用户的所有看板
func (h *BoardHandler) GetBoards(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	boards, err := h.boardService.GetBoardsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"boards": boards})
}

// GET /api/projects/:projectId/boards
func (h *BoardHandler) GetBoardsByProject(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID"})
		return
	}

	boards, err := h.boardService.GetBoardsByProjectID(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"boards": boards})
}

// CreateBoard 创建看板
func (h *BoardHandler) CreateBoard(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID"})
		return
	}

	var createBoardRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
		// ProjectID   uint   `json:"project_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&createBoardRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	board := &models.Board{
		Name:        createBoardRequest.Name,
		Description: createBoardRequest.Description,
		Color:       createBoardRequest.Color,
		ProjectID:   uint(projectID),
		OwnerID:     userID,
		Status:      models.BoardStatusActive,
	}

	if err := h.boardService.CreateBoard(board); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, board)
}

// UpdateBoard 更新看板
func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	var updateBoardRequest struct {
		Name        *string             `json:"name"`
		Description *string             `json:"description"`
		Color       *string             `json:"color"`
		Status      *models.BoardStatus `json:"status"`
		Position    *int                `json:"position"`
	}

	if err := c.ShouldBindJSON(&updateBoardRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if updateBoardRequest.Name != nil {
		updates["name"] = *updateBoardRequest.Name
	}
	if updateBoardRequest.Description != nil {
		updates["description"] = *updateBoardRequest.Description
	}
	if updateBoardRequest.Color != nil {
		updates["color"] = *updateBoardRequest.Color
	}
	if updateBoardRequest.Status != nil {
		updates["status"] = *updateBoardRequest.Status
	}
	if updateBoardRequest.Position != nil {
		updates["position"] = *updateBoardRequest.Position
	}

	if err := h.boardService.UpdateBoard(uint(boardID), updates); err != nil {
		if err == services.ErrBoardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteBoard 删除看板
func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	if err := h.boardService.DeleteBoard(uint(boardID)); err != nil {
		if err == services.ErrBoardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
