package useCase

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
)

type UseCases struct {
	GcpDatacenter      GCPDatacenter
	GcpCluster         GCPCluster
	Cluster            Cluster
	Datacenter         Datacenter
	Event              Event
	ScheduledHPAConfig ScheduledHPAConfig
}

func BuildUseCases(
	resources *config.KubeEPResources,
	repositories *repository.Repositories,
) *UseCases {
	return &UseCases{
		GcpCluster: newGCPCluster(
			resources.ValidatorInst, repositories.Cluster,
			repositories.GCPCluster, repositories.K8SDiscovery,
		),
		GcpDatacenter: newGCPDatacenter(repositories.Datacenter, resources.ValidatorInst),
		Cluster: newCluster(
			resources.ValidatorInst,
			repositories.Cluster,
			repositories.K8sHPA,
			repositories.K8sNamespace,
			repositories.K8SDiscovery,
			repositories.K8sDeployment,
		),
		Datacenter: newDatacenter(resources.ValidatorInst, repositories.Datacenter),
		Event: newEvent(
			resources.ValidatorInst,
			repositories.Event,
			repositories.ScheduledHPAConfig,
			repositories.Cluster,
		),
		ScheduledHPAConfig: newScheduledHPAConfig(repositories.ScheduledHPAConfig),
	}
}
