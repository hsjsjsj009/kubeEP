package handler

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
)

type Handlers struct {
	GcpHandler     Gcp
	ClusterHandler Cluster
	EventHandler   Event
}

func BuildHandlers(useCases *useCase.UseCases, resources *config.KubeEPResources) *Handlers {
	return &Handlers{
		GcpHandler: newGCPHandler(
			resources.ValidatorInst,
			useCases.GcpCluster,
			useCases.GcpDatacenter,
			resources.DB,
			useCases.Cluster,
		),
		ClusterHandler: newClusterHandler(
			resources.ValidatorInst,
			useCases.Cluster,
			resources.DB,
			useCases.GcpCluster,
			useCases.GcpDatacenter,
			useCases.Datacenter,
		),
		EventHandler: newEventHandler(
			resources.ValidatorInst,
			useCases.Event,
			useCases.ScheduledHPAConfig,
			useCases.HPAConfigStatus,
			resources.DB,
		),
	}

}
