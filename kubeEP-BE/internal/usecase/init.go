package useCase

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
)

type UseCases struct {
	GcpUseCaseDC      GCPDatacenter
	GcpUseCaseCluster GCPCluster
	UseCaseCluster    Cluster
	UseCaseDC         Datacenter
}

func BuildUseCases(resources *config.KubeEPResources, repositories *repository.Repositories) *UseCases {
	return &UseCases{
		GcpUseCaseCluster: NewGCPCluster(resources.ValidatorInst, repositories.Cluster),
		GcpUseCaseDC:      NewGCPDatacenter(repositories.Datacenter, resources.ValidatorInst),
		UseCaseCluster:    NewCluster(resources.ValidatorInst, repositories.Cluster, repositories.K8sHPA, repositories.K8sNamespace),
		UseCaseDC:         NewDatacenter(resources.ValidatorInst, repositories.Datacenter),
	}
}
