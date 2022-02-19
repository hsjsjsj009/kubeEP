package UCEntity

type GCPClusterData struct {
	ClusterData
	Location string
}

type GCPClusterMetaData struct {
	Location string `json:"location"`
}
