package cron

type NodePoolRequestedResourceData struct {
	MaxCPU    float64
	MaxMemory float64
	MaxPods   int64
}

type NodePoolResourceData struct {
	MaxAvailablePods   int64
	MaxAvailableCPU    float64
	MaxAvailableMemory float64
	AvailableCPU       float64
	AvailableMemory    float64
	AvailablePods      int64
}
