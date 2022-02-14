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
		Cluster:            NewCluster(),
		Datacenter:         NewDatacenter(resources.Redis),
		Event:              NewEvent(),
		ScheduledHPAConfig: NewScheduledHPAConfig(),
	}
}
