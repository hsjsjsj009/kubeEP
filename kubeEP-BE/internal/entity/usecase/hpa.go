package UCEntity

type HPAScaleTargetRef struct {
	Name string
	Kind string
}

type SimpleHPAData struct {
	Name            string
	Namespace       string
	MinReplicas     *int32
	MaxReplicas     int32
	CurrentReplicas int32
	ScaleTargetRef  HPAScaleTargetRef
}

type EventModifiedHPAData struct {
	Name        string
	Namespace   string
	MinReplicas *int32
	MaxReplicas int32
}
