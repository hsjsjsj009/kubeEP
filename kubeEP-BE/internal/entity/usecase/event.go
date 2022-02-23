package UCEntity

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID        uuid.UUID
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Cluster   ClusterData
}
