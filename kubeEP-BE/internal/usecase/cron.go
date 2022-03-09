package useCase

import "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"

type Cron interface {
}

type cron struct {
	datacenterRepo repository.Datacenter
	eventRepo      repository.Event
}

func newCron(datacenterRepo repository.Datacenter, eventRepo repository.Event) Cron {
	return &cron{datacenterRepo: datacenterRepo, eventRepo: eventRepo}
}
