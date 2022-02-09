package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type Event interface {
	GetEventByID(id uuid.UUID) (*model.Event, error)
	ListEventByClusterID(id uuid.UUID) ([]*model.Event, error)
	InsertEvent(data *model.Event) error
}

type event struct {
	db *gorm.DB
}

func NewEvent(db *gorm.DB) Event {
	return &event{
		db: db,
	}
}

func (e *event) GetEventByID(id uuid.UUID) (*model.Event, error) {
	data := &model.Event{}
	tx := e.db.Model(data).First(data, id)
	return data, tx.Error
}

func (e *event) ListEventByClusterID(id uuid.UUID) ([]*model.Event, error) {
	var data []*model.Event
	tx := e.db.Model(&model.Event{}).Where("cluster_id = ?", id).Scan(&data)
	return data, tx.Error
}

func (e *event) InsertEvent(data *model.Event) error {
	return e.db.Create(data).Error
}
