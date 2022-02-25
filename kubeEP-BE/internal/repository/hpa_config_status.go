package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type HPAConfigStatus interface {
	GetHPAConfigStatusByID(tx *gorm.DB, id uuid.UUID) (*model.HPAConfigStatus, error)
	GetHPAConfigStatusByScheduledHPAConfigID(
		tx *gorm.DB,
		scheduledHPAConfigId uuid.UUID,
	) (*model.HPAConfigStatus, error)
	InsertHPAConfigStatus(tx *gorm.DB, data *model.HPAConfigStatus) error
	InsertBatchHPAConfigStatus(
		tx *gorm.DB,
		data []*model.HPAConfigStatus,
	) error
}

type hpaConfigStatus struct {
}

func newHpaConfigStatus() HPAConfigStatus {
	return &hpaConfigStatus{}
}

func (h *hpaConfigStatus) GetHPAConfigStatusByID(tx *gorm.DB, id uuid.UUID) (
	*model.HPAConfigStatus,
	error,
) {
	data := &model.HPAConfigStatus{}
	err := tx.First(data, id).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (h *hpaConfigStatus) GetHPAConfigStatusByScheduledHPAConfigID(
	tx *gorm.DB,
	scheduledHPAConfigId uuid.UUID,
) (*model.HPAConfigStatus, error) {
	data := &model.HPAConfigStatus{}
	err := tx.Table("hpa_config_status h").Select("h.*").Joins(
		`left join scheduled_hpa_configs s on s.id = h.scheduled_hpa_config_id and 
		s.deleted_at is null`,
	).Where(`s.id = ?`, scheduledHPAConfigId).Scan(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (h *hpaConfigStatus) InsertHPAConfigStatus(tx *gorm.DB, data *model.HPAConfigStatus) error {
	return tx.Create(data).Error
}

func (h *hpaConfigStatus) InsertBatchHPAConfigStatus(
	tx *gorm.DB,
	data []*model.HPAConfigStatus,
) error {
	return tx.Create(data).Error
}
