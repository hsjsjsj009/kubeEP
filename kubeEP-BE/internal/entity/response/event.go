package response

import "github.com/google/uuid"

type EventCreationResponse struct {
	EventID uuid.UUID `json:"event_id"`
}
