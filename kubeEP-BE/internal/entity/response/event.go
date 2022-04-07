package response

import (
	"github.com/google/uuid"
	"time"
)

type EventCreationResponse struct {
	EventID uuid.UUID `json:"event_id"`
}

type EventSimpleResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type EventDetailedResponse struct {
	EventSimpleResponse
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	Cluster            Cluster             `json:"cluster"`
	ModifiedHPAConfigs []ModifiedHPAConfig `json:"modified_hpa_configs"`
	UpdatedNodePools   []UpdatedNodePool   `json:"updated_node_pools"`
}
