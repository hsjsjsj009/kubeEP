package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type Event interface {
	GetEventByID(tx *gorm.DB, id uuid.UUID) (*model.Event, error)
	ListEventByClusterID(tx *gorm.DB, id uuid.UUID) ([]*model.Event, error)
	InsertEvent(tx *gorm.DB, data *model.Event) error
}

type event struct {
}

func NewEvent() Event {
	return &event{}
}

func (e *event) GetEventByID(tx *gorm.DB, id uuid.UUID) (*model.Event, error) {
	data := &model.Event{}
	tx = tx.Model(data).First(data, id)
	return data, tx.Error
}

func (e *event) ListEventByClusterID(tx *gorm.DB, id uuid.UUID) ([]*model.Event, error) {
	var data []*model.Event
	tx = tx.Model(&model.Event{}).Where("cluster_id = ?", id).Find(&data)
	return data, tx.Error
}

func (e *event) InsertEvent(tx *gorm.DB, data *model.Event) error {
	return tx.Create(data).Error
}
