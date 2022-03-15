package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
	"time"
)

type Event interface {
	GetEventByID(tx *gorm.DB, id uuid.UUID) (*model.Event, error)
	GetEventByName(tx *gorm.DB, name string) (*model.Event, error)
	ListEventByClusterID(tx *gorm.DB, id uuid.UUID) ([]*model.Event, error)
	InsertEvent(tx *gorm.DB, data *model.Event) error
	SaveEvent(tx *gorm.DB, data *model.Event) error
	DeleteEvent(tx *gorm.DB, id uuid.UUID) error
	FindPendingEventWithIntervalMinute(tx *gorm.DB, minute int, now time.Time) (
		[]*model.Event,
		error,
	)
}

type event struct {
}

func newEvent() Event {
	return &event{}
}

func (e *event) GetEventByID(tx *gorm.DB, id uuid.UUID) (*model.Event, error) {
	data := &model.Event{}
	tx = tx.Model(data).First(data, id)
	return data, tx.Error
}

func (e *event) GetEventByName(tx *gorm.DB, name string) (*model.Event, error) {
	data := &model.Event{}
	tx = tx.Model(data).Where("name = ?", name).First(data)
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

func (e *event) SaveEvent(tx *gorm.DB, data *model.Event) error {
	return tx.Save(data).Error
}

func (e *event) DeleteEvent(tx *gorm.DB, id uuid.UUID) error {
	return tx.Delete(&model.Event{}, "id = ?", id).Error
}

func (e *event) FindPendingEventWithIntervalMinute(tx *gorm.DB, minute int, now time.Time) (
	[]*model.Event,
	error,
) {
	var data []*model.Event
	tx = tx.Model(&model.Event{}).Where(
		"date_trunc('minutes', start_time) - date_trunc('minutes', ?) = interval '? minutes' and status = ?",
		now.UTC(),
		minute,
		model.EventPending,
	).Find(&data)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return data, nil
}
