package request

import (
	"github.com/google/uuid"
	"time"
)

type EventDataRequest struct {
	Name               *string                      `json:"name" validate:"required"`
	StartTime          *time.Time                   `json:"start_time" validate:"required"`
	EndTime            *time.Time                   `json:"end_time" validate:"required,gtefield=StartTime"`
	ClusterID          *uuid.UUID                   `json:"cluster_id" validate:"required"`
	ModifiedHPAConfigs []EventModifiedHPAConfigData `json:"modified_hpa_configs" validate:"required,min=1,dive"`
}

type EventListRequest struct {
	ClusterID *uuid.UUID `query:"cluster_id" validate:"required"`
}

type UpdateEventDataRequest struct {
	EventDataRequest
	EventID *uuid.UUID `json:"event_id" validator:"required"`
}
