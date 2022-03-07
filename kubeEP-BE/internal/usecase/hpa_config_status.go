package useCase

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type HPAConfigStatus interface {
	CreateHPAConfigStatusForScheduledConfigIDs(
		tx *gorm.DB,
		scheduledConfigIDs []uuid.UUID,
	) ([]uuid.UUID, error)
	SoftDeleteHPAConfigStatusByEventID(tx *gorm.DB, eventID uuid.UUID) error
}

type hpaConfigStatus struct {
	hpaConfigStatusRepo repository.HPAConfigStatus
}

func newHpaConfigStatus(hpaConfigStatusRepo repository.HPAConfigStatus) HPAConfigStatus {
	return &hpaConfigStatus{hpaConfigStatusRepo: hpaConfigStatusRepo}
}

func (h *hpaConfigStatus) CreateHPAConfigStatusForScheduledConfigIDs(
	tx *gorm.DB,
	scheduledConfigIDs []uuid.UUID,
) ([]uuid.UUID, error) {
	var hpaConfigStatusData []*model.HPAConfigStatus
	for _, scheduledConfigID := range scheduledConfigIDs {
		modelData := &model.HPAConfigStatus{
			Status: model.HPAConfigCreated,
		}
		modelData.ScheduledHPAConfigID.SetUUID(scheduledConfigID)
		hpaConfigStatusData = append(hpaConfigStatusData, modelData)
	}

	err := h.hpaConfigStatusRepo.InsertBatchHPAConfigStatus(tx, hpaConfigStatusData)
	if err != nil {
		return nil, err
	}

	var uuids []uuid.UUID
	for _, hpaConfigStatus := range hpaConfigStatusData {
		uuids = append(uuids, hpaConfigStatus.ID.GetUUID())
	}
	return uuids, nil
}

func (h *hpaConfigStatus) SoftDeleteHPAConfigStatusByEventID(tx *gorm.DB, eventID uuid.UUID) error {
	return h.hpaConfigStatusRepo.DeleteHPAConfigStatusByEventID(tx, eventID)
}
