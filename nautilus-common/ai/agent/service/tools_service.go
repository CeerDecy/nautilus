package service

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"nautilus/nautilus-common/mysql/model"
)

type ToolsService struct {
	db *gorm.DB
}

func NewToolsService(db *gorm.DB) *ToolsService {
	return &ToolsService{
		db: db,
	}
}

func (t *ToolsService) GetToolsByRole(role string) []model.AITool {
	var res []model.AITool
	if err := t.db.Model(&model.AITool{}).Where("role = ? or role = '*'", role).Find(&res).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("get tools by role error: %v", err)
		}
	}
	return res
}
