package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type ScheduledHPAConfig interface {
	GetScheduledHPAConfigByID(tx *gorm.DB, id uuid.UUID) (*model.ScheduledHPAConfig, error)
	ListScheduledHPAConfigByEventID(tx *gorm.DB, id uuid.UUID) ([]*model.ScheduledHPAConfig, error)
	InsertScheduledHPAConfig(tx *gorm.DB, data *model.ScheduledHPAConfig) error
}

type scheduledHPAConfig struct {
}

func NewScheduledHPAConfig() ScheduledHPAConfig {
	return &scheduledHPAConfig{}
}

func (s *scheduledHPAConfig) GetScheduledHPAConfigByID(tx *gorm.DB, id uuid.UUID) (*model.ScheduledHPAConfig, error) {
	data := &model.ScheduledHPAConfig{}
	tx = tx.Model(data).First(data, id)
	return data, tx.Error
}

func (s *scheduledHPAConfig) ListScheduledHPAConfigByEventID(tx *gorm.DB, id uuid.UUID) ([]*model.ScheduledHPAConfig, error) {
	var data []*model.ScheduledHPAConfig
	tx = tx.Model(&model.ScheduledHPAConfig{}).Where("event_id = ?", id).Find(&data)
	return data, tx.Error
}

func (s *scheduledHPAConfig) InsertScheduledHPAConfig(tx *gorm.DB, data *model.ScheduledHPAConfig) error {
	return tx.Create(data).Error
}
