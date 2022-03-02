package useCase

import (
	"github.com/google/uuid"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type ScheduledHPAConfig interface {
	RegisterModifiedHPAConfigs(
		tx *gorm.DB,
		modifiedHPAs []UCEntity.EventModifiedHPAConfigData,
		eventID uuid.UUID,
	) ([]uuid.UUID, error)
	DeleteEventModifiedHPAConfigs(tx *gorm.DB, eventID uuid.UUID) error
}

type scheduledHPAConfig struct {
	scheduledHPAConfigRepo repository.ScheduledHPAConfig
}

func newScheduledHPAConfig(scheduledHPAConfigRepo repository.ScheduledHPAConfig) ScheduledHPAConfig {
	return &scheduledHPAConfig{scheduledHPAConfigRepo: scheduledHPAConfigRepo}
}

func (s *scheduledHPAConfig) RegisterModifiedHPAConfigs(
	tx *gorm.DB,
	modifiedHPAs []UCEntity.EventModifiedHPAConfigData,
	eventID uuid.UUID,
) ([]uuid.UUID, error) {
	var data []*model.ScheduledHPAConfig
	for _, modifiedHPA := range modifiedHPAs {
		modelData := &model.ScheduledHPAConfig{
			Name:      modifiedHPA.Name,
			MinPods:   modifiedHPA.MinReplicas,
			MaxPods:   modifiedHPA.MaxReplicas,
			Namespace: modifiedHPA.Namespace,
		}
		modelData.EventID.SetUUID(eventID)
		data = append(
			data, modelData,
		)
	}
	err := s.scheduledHPAConfigRepo.InsertBatchScheduledHPAConfig(tx, data)
	if err != nil {
		return nil, err
	}
	var uuids []uuid.UUID
	for _, datum := range data {
		uuids = append(uuids, datum.ID.GetUUID())
	}
	return uuids, nil
}

func (s *scheduledHPAConfig) DeleteEventModifiedHPAConfigs(tx *gorm.DB, eventID uuid.UUID) error {
	return s.scheduledHPAConfigRepo.DeletePermanentAllHPAConfigByEventID(tx, eventID)
}
