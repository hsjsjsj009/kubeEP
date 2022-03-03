package UCEntity

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        uuid.UUID
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Cluster   ClusterData
}

type DetailedEvent struct {
	Event
	EventModifiedHPAConfigData []EventModifiedHPAConfigData
}
