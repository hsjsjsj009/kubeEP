package UCEntity

import "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"

type HPAScaleTargetRef struct {
	Name string
	Kind string
}

type SimpleHPAData struct {
	APIVersion      constant.HPAVersion
	Name            string
	Namespace       string
	MinReplicas     *int32
	MaxReplicas     int32
	CurrentReplicas int32
	ScaleTargetRef  HPAScaleTargetRef
}

type EventModifiedHPAData struct {
	APIVersion  constant.HPAVersion
	Name        string
	Namespace   string
	MinReplicas *int32
	MaxReplicas int32
}
