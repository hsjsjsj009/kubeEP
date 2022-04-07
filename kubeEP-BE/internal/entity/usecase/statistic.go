package UCEntity

import (
	"github.com/google/uuid"
	"time"
)

type UpdatedNodePoolData struct {
	UpdatedNodePoolID uuid.UUID
	NodePoolName      string
}

type NodePoolStatusData struct {
	CreatedAt time.Time
	Count     int32
}

type HPAStatusData struct {
	CreatedAt time.Time
	Count     int32
}
