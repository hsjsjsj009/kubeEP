package model

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
	"time"
)

type Event struct {
	BaseModel
	Name      string
	StartTime time.Time
	EndTime   time.Time
	ClusterID gormDatatype.UUID
	Cluster   Cluster `gorm:"ForeignKey:ClusterID;constraint:OnDelete:CASCADE"`
}

func (e *Event) TableName() string {
	return "events"
}
