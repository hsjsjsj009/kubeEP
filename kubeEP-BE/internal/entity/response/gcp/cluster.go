package gcpResponse

import "github.com/google/uuid"

type Cluster struct {
	ID       *uuid.UUID `json:"id,omitempty"`
	Name     string     `json:"name"`
	Location string     `json:"location"`
}

type DatacenterClusters struct {
	Clusters              []Cluster `json:"clusters"`
	IsTemporaryDatacenter bool      `json:"is_temporary_datacenter"`
}
