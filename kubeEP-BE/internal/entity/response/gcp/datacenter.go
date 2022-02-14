package gcpResponse

import "github.com/google/uuid"

type DatacenterData struct {
	DatacenterID uuid.UUID `json:"datacenter_id"`
	IsTemporary  bool      `json:"is_temporary"`
}
