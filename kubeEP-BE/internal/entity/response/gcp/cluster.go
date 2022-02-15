package gcpResponse

type Cluster struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type DatacenterClusters struct {
	Clusters              []Cluster `json:"clusters"`
	IsTemporaryDatacenter bool      `json:"is_temporary_datacenter"`
}
