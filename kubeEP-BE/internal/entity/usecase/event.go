package UCEntity

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID        uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	Cluster   ClusterData
}