package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type ScheduledHPAConfig interface {
	GetScheduledHPAConfigByID(id uuid.UUID) (*model.ScheduledHPAConfig, error)
	ListScheduledHPAConfigByEventID(id uuid.UUID) ([]*model.ScheduledHPAConfig, error)
	InsertScheduledHPAConfig(data *model.ScheduledHPAConfig) error
}

type scheduledHPAConfig struct {
	db *gorm.DB
}

func NewScheduledHPAConfig(db *gorm.DB) ScheduledHPAConfig {
	return &scheduledHPAConfig{
		db: db,
	}
}

func (s *scheduledHPAConfig) GetScheduledHPAConfigByID(id uuid.UUID) (*model.ScheduledHPAConfig, error) {
	data := &model.ScheduledHPAConfig{}
	tx := s.db.Model(data).First(data, id)
	return data, tx.Error
}

func (s *scheduledHPAConfig) ListScheduledHPAConfigByEventID(id uuid.UUID) ([]*model.ScheduledHPAConfig, error) {
	var data []*model.ScheduledHPAConfig
	tx := s.db.Model(&model.ScheduledHPAConfig{}).Where("event_id = ?", id).Scan(&data)
	return data, tx.Error
}

func (s *scheduledHPAConfig) InsertScheduledHPAConfig(data *model.ScheduledHPAConfig) error {
	return s.db.Create(data).Error
}
