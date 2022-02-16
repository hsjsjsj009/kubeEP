package gcpRequest

import (
	"encoding/json"
	"github.com/google/uuid"
)

type DatacenterData struct {
	Name             *string          `json:"name" validate:"required"`
	SAKeyCredentials *json.RawMessage `json:"sa_key_credentials" validate:"required"`
	IsTemporary      *bool            `json:"is_temporary" validate:"required"`
}

type ExistingDatacenterData struct {
	DatacenterID *uuid.UUID `json:"datacenter_id" query:"datacenter_id" validate:"required"`
}
