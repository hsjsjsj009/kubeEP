package model

import gormDatatype "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"

type HPAUpdateStatus int8

const (
	HPAUpdateFailed HPAUpdateStatus = iota
	HPAUpdateSuccess
	HPAUpdatePending
	HPAUpdateScheduled
)

type HPAConfigStatus struct {
	BaseModel
	Status               HPAUpdateStatus
	Message              *string
	CAMaxNode            *int
	CurrentPodsCount     int
	ScheduledHPAConfigID gormDatatype.UUID  `gorm:"column:scheduled_hpa_config_id"`
	ScheduledHPAConfig   ScheduledHPAConfig `gorm:"ForeignKey:ScheduledHPAConfigID;constraint:OnDelete:CASCADE"`
}

func (HPAConfigStatus) TableName() string {
	return "hpa_config_status"
}
