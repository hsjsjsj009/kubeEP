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
		GcpUseCaseCluster: newGCPCluster(resources.ValidatorInst, repositories.Cluster),
		GcpUseCaseDC:      newGCPDatacenter(repositories.Datacenter, resources.ValidatorInst),
		UseCaseCluster:    newCluster(resources.ValidatorInst, repositories.Cluster, repositories.K8sHPA, repositories.K8sNamespace),
		UseCaseDC:         newDatacenter(resources.ValidatorInst, repositories.Datacenter),
	}
}
