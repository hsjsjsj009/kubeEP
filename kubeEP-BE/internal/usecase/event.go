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
}

type event struct {
	validatorInst   *validator.Validate
	eventRepository repository.Event
}

func NewEvent(
	validatorInst *validator.Validate,
	eventRepository repository.Event,
) Event {
	return &event{
		validatorInst:   validatorInst,
		eventRepository: eventRepository,
	}
}

func (e *event) RegisterEvents(tx *gorm.DB, eventData *UCEntity.Event) (uuid.UUID, error) {
	data := &model.Event{
		StartTime: eventData.StartTime,
		EndTime:   eventData.EndTime,
	}
	data.ID.SetUUID(eventData.ID)
	data.ClusterID.SetUUID(eventData.Cluster.ID)

	err := e.eventRepository.InsertEvent(tx, data)
	if err != nil {
		return uuid.UUID{}, err
	}
	return data.ID.GetUUID(), nil
}

//func (e *event) RegisterModifiedHPAs(
//	tx *gorm.DB,
//	modifiedHPAs []UCEntity.EventModifiedHPAData,
//	eventID uuid.UUID,
//) ([]UCEntity.EventModifiedHPAData, error) {
//	var data []*model.ScheduledHPAConfig
//	for _, modifiedHPA := range modifiedHPAs {
//		data = append(data, &model.ScheduledHPAConfig{})
//	}
//}
