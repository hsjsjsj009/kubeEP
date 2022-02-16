package handler

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
)

type Handlers struct {
	GcpHandler     GcpHandler
	ClusterHandler ClusterHandler
}

func BuildHandlers(useCases *useCase.UseCases, resources *config.KubeEPResources) *Handlers {
	return &Handlers{
		GcpHandler: newGCPHandler(
			resources.ValidatorInst,
			useCases.GcpUseCaseCluster,
			useCases.GcpUseCaseDC,
			resources.DB,
			useCases.UseCaseCluster,
		),
		ClusterHandler: newClusterHandler(
			resources.ValidatorInst,
			useCases.UseCaseCluster,
			resources.DB,
		),
	}

}
