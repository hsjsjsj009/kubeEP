package response

type SimpleHPA struct {
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	MinReplicas     *int32 `json:"min_replicas,omitempty"`
	MaxReplicas     int32  `json:"max_replicas"`
	CurrentReplicas int32  `json:"current_replicas"`
}

type ModifiedHPAConfig struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	MinReplicas *int32 `json:"min_replicas,omitempty"`
	MaxReplicas int32  `json:"max_replicas"`
}
