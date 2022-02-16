package useCase

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	gcpUseCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase/gcp"
)

type UseCases struct {
	GcpUseCaseDC      gcpUseCase.Datacenter
	GcpUseCaseCluster gcpUseCase.Cluster
	UseCaseCluster    Cluster
}

func BuildUseCases(resources *config.KubeEPResources, repositories *repository.Repositories) *UseCases {
	return &UseCases{
		GcpUseCaseCluster: gcpUseCase.NewCluster(resources.ValidatorInst, repositories.Cluster),
		GcpUseCaseDC:      gcpUseCase.NewDatacenter(repositories.Datacenter, resources.ValidatorInst),
		UseCaseCluster:    NewCluster(resources.ValidatorInst, repositories.Cluster),
	}
}
