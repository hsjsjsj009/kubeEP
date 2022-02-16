package gcpResponse

import "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response"

type Cluster struct {
	response.Cluster
	Location string `json:"location"`
}

type DatacenterClusters struct {
	Clusters              []Cluster `json:"clusters"`
	IsTemporaryDatacenter bool      `json:"is_temporary_datacenter"`
}
