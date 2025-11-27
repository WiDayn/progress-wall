package services

import (
	"errors"
	"fmt"

	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// BoardService 看板服务
type BoardService struct {
	db *gorm.DB
}

// NewBoardService 创建看板服务
func NewBoardService(db *gorm.DB) *BoardService {
	return &BoardService{
		db: db,
	}
}

// GetBoardByID 根据ID获取看板（包含嵌套的列和任务）
func (s *BoardService) GetBoardByID(boardID uint) (*models.Board, error) {
	var board models.Board
	result := s.db.
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Columns.Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC").Preload("Assignee")
		}).
		Preload("Columns.Tasks.Creator").
		Preload("Owner").
		First(&board, boardID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, fmt.Errorf("查询看板失败: %v", result.Error)
	}

	return &board, nil
}

// GetBoardsByUserID 获取用户的所有看板
func (s *BoardService) GetBoardsByUserID(userID uint) ([]models.Board, error) {
	var boards []models.Board
	result := s.db.
		Where("owner_id = ?", userID).
		Order("position ASC").
		Find(&boards)

	if result.Error != nil {
		return nil, errors.New("查询看板列表失败")
	}

	return boards, nil
}

// GetBoardsByProjectID 获取项目的所有看板
func (s *BoardService) GetBoardsByProjectID(projectID uint) ([]models.Board, error) {
	var boards []models.Board
	err := s.db.
		Where("project_id = ?", projectID).
		Order("position ASC").
		Preload("Columns").
		Find(&boards).Error

	if err != nil {
		return nil, fmt.Errorf("查询项目看板失败: %v", err)
	}

	return boards, nil
}

// CreateBoard 创建看板（带默认列）
func (s *BoardService) CreateBoard(board *models.Board) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(board).Error; err != nil {
			return fmt.Errorf("创建看板失败: %v", err)
		}

		// 初始化默认列
		defaultColumns := []models.Column{
			{
				Name:        "Backlog",
				Description: "待办事项",
				Color:       "#6B7280", // Gray
				Position:    1000,
				BoardID:     board.ID,
				Status:      models.ColumnStatusActive,
			},
			{
				Name:        "Ready",
				Description: "准备就绪",
				Color:       "#3B82F6", // Blue
				Position:    2000,
				BoardID:     board.ID,
				Status:      models.ColumnStatusActive,
			},
			{
				Name:        "In processing",
				Description: "进行中",
				Color:       "#F59E0B", // Yellow
				Position:    3000,
				BoardID:     board.ID,
				Status:      models.ColumnStatusActive,
			},
			{
				Name:        "In review",
				Description: "审核中",
				Color:       "#8B5CF6", // Purple
				Position:    4000,
				BoardID:     board.ID,
				Status:      models.ColumnStatusActive,
			},
			{
				Name:        "Done",
				Description: "已完成",
				Color:       "#10B981", // Green
				Position:    5000,
				BoardID:     board.ID,
				Status:      models.ColumnStatusActive,
			},
		}

		if err := tx.Create(&defaultColumns).Error; err != nil {
			return fmt.Errorf("创建默认列失败: %v", err)
		}

		return nil
	})
}

// UpdateBoard 更新看板
func (s *BoardService) UpdateBoard(boardID uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Board{}).Where("id = ?", boardID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新看板失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrBoardNotFound
	}
	return nil
}

// DeleteBoard 删除看板（软删除）
func (s *BoardService) DeleteBoard(boardID uint) error {
	result := s.db.Delete(&models.Board{}, boardID)
	if result.Error != nil {
		return fmt.Errorf("删除看板失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrBoardNotFound
	}
	return nil
}
