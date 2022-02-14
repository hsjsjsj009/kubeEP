package handler

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
)

type Handlers struct {
	GcpHandler GcpHandler
}

func BuildHandlers(useCases *usecase.UseCases, resources *config.KubeEPResources) *Handlers {
	return &Handlers{
		GcpHandler: newGCPHandler(
			resources.ValidatorInst,
			useCases.GcpUseCaseCluster,
			useCases.GcpUseCaseDC,
			resources.DB,
		),
	}

}
