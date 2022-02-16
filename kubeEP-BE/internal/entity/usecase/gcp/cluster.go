package gcpUCEntity

import UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"

type ClusterData struct {
	UCEntity.ClusterData
	Location string
}

type ClusterMetaData struct {
	Location string `json:"location"`
}
