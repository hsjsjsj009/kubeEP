package model

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
)

type ScheduledHPAConfig struct {
	BaseModel
	Name                   string
	Spec                   gormDatatype.JSON
	ServiceTestingEndpoint *string
	Namespace              string
	EventID                gormDatatype.UUID
	Event                  Event `gorm:"ForeignKey:EventID;constraint:OnDelete:CASCADE"`
}

func (s *ScheduledHPAConfig) TableName() string {
	return "scheduled_hpa_configs"
}
