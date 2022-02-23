package request

import (
	"github.com/google/uuid"
	"time"
)

type EventRequest struct {
	Name               *string                      `json:"name" validate:"required"`
	StartTime          *time.Time                   `json:"start_time" validate:"required"`
	EndTime            *time.Time                   `json:"end_time" validate:"required"`
	ClusterID          *uuid.UUID                   `json:"cluster_id" validate:"required"`
	ModifiedHPAConfigs []EventModifiedHPAConfigData `json:"modified_hpa_configs" validate:"required,min=1,dive"`
}
