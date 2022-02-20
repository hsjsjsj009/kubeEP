package repository

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type Repositories struct {
	Cluster            Cluster
	Datacenter         Datacenter
	Event              Event
	ScheduledHPAConfig ScheduledHPAConfig
	K8sHPA             K8sHPA
	K8sNamespace       K8sNamespace
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Datacenter{},
		&model.Cluster{},
		&model.Event{},
		&model.ScheduledHPAConfig{},
	)
}

func BuildRepositories(resources *config.KubeEPResources) *Repositories {
	return &Repositories{
		Cluster:            newCluster(),
		Datacenter:         newDatacenter(resources.Redis),
		Event:              newEvent(),
		ScheduledHPAConfig: newScheduledHPAConfig(),
		K8sHPA:             newK8sHPA(resources.Redis),
		K8sNamespace:       newK8sNamespace(),
	}
}
