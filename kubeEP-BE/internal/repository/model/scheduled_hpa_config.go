package model

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
)

type HPAUpdateStatus int8

const (
	HPAUpdateFailed HPAUpdateStatus = iota
	HPAUpdateSuccess
	HPAUpdatePending
	HPAUpdateScheduled
)

type ScheduledHPAConfig struct {
	BaseModel
	Name                   string
	MinPods                *int32
	MaxPods                int32
	ServiceTestingEndpoint *string
	Namespace              string
	EventID                gormDatatype.UUID
	Status                 HPAUpdateStatus
	StatusDetailMessage    string
	Event                  Event `gorm:"ForeignKey:EventID;constraint:OnDelete:CASCADE"`
}

func (s *ScheduledHPAConfig) TableName() string {
	return "scheduled_hpa_configs"
}
