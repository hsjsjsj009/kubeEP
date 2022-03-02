package useCase

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type Event interface {
	RegisterEvents(tx *gorm.DB, eventData *UCEntity.Event) (uuid.UUID, error)
	GetEventByName(tx *gorm.DB, eventName string) (*UCEntity.Event, error)
	ListEventByClusterID(tx *gorm.DB, clusterID uuid.UUID) ([]UCEntity.Event, error)
	UpdateEvent(tx *gorm.DB, eventData *UCEntity.Event) error
}

type event struct {
	validatorInst                *validator.Validate
	eventRepository              repository.Event
	scheduledHPAConfigRepository repository.ScheduledHPAConfig
}

func newEvent(
	validatorInst *validator.Validate,
	eventRepository repository.Event,
	scheduledHPAConfigRepository repository.ScheduledHPAConfig,
) Event {
	return &event{
		validatorInst:                validatorInst,
		eventRepository:              eventRepository,
		scheduledHPAConfigRepository: scheduledHPAConfigRepository,
	}
}

func (e *event) RegisterEvents(tx *gorm.DB, eventData *UCEntity.Event) (uuid.UUID, error) {
	data := &model.Event{
		Name:      eventData.Name,
		StartTime: eventData.StartTime,
		EndTime:   eventData.EndTime,
	}
	data.ClusterID.SetUUID(eventData.Cluster.ID)

	err := e.eventRepository.InsertEvent(tx, data)
	if err != nil {
		return uuid.UUID{}, err
	}
	return data.ID.GetUUID(), nil
}

func (e *event) GetEventByName(tx *gorm.DB, eventName string) (*UCEntity.Event, error) {
	data, err := e.eventRepository.GetEventByName(tx, eventName)
	if err != nil {
		return nil, err
	}
	return &UCEntity.Event{
		ID:        data.ID.GetUUID(),
		Name:      data.Name,
		StartTime: data.StartTime,
		EndTime:   data.EndTime,
	}, nil
}

func (e *event) ListEventByClusterID(tx *gorm.DB, clusterID uuid.UUID) ([]UCEntity.Event, error) {
	events, err := e.eventRepository.ListEventByClusterID(tx, clusterID)
	if err != nil {
		return nil, err
	}
	var output []UCEntity.Event
	for _, event := range events {
		output = append(
			output, UCEntity.Event{
				ID:        event.ID.GetUUID(),
				Name:      event.Name,
				StartTime: event.StartTime,
				EndTime:   event.EndTime,
			},
		)
	}
	return output, nil
}

func (e *event) UpdateEvent(tx *gorm.DB, eventData *UCEntity.Event) error {
	data := &model.Event{
		Name:      eventData.Name,
		StartTime: eventData.StartTime,
		EndTime:   eventData.EndTime,
	}
	data.ID.SetUUID(eventData.ID)
	data.ClusterID.SetUUID(eventData.Cluster.ID)
	return e.eventRepository.SaveEvent(tx, data)
}
