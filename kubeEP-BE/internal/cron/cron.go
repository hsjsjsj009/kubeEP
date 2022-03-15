package cron

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"github.com/google/martian/log"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"time"
)

type Cron interface {
	Start()
}

type cron struct {
	eventUC              useCase.Event
	clusterUC            useCase.Cluster
	generalClusterUC     useCase.Cluster
	gcpClusterUC         useCase.GCPCluster
	gcpDatacenterUC      useCase.GCPDatacenter
	scheduledHPAConfigUC useCase.ScheduledHPAConfig
	tx                   *gorm.DB
}

func newCron(
	eventUC useCase.Event,
) Cron {
	return &cron{eventUC: eventUC}
}

func (c *cron) handleError(db *gorm.DB, e *UCEntity.Event, errMsg string) {
	e.Status = model.EventFailed
	e.Message = errMsg
	err := c.eventUC.UpdateEvent(db, e)
	if err != nil {
		log.Errorf("[EventCronJob] Error Update Event : %s", err.Error())
	}
	log.Errorf("[EventCronJob] Event (%s) Error : %s", e.Name, errMsg)
}

func (c *cron) execEvent(e *UCEntity.Event, db *gorm.DB, ctx context.Context) {
	e.Status = model.EventExecuting

	err := c.eventUC.UpdateEvent(db, e)
	if err != nil {
		log.Errorf("[EventCronJob] Error Update Event : %s", err.Error())
		return
	}

	clusterID := e.Cluster.ID
	clusterData, err := c.clusterUC.GetClusterAndDatacenterDataByClusterID(db, clusterID)
	if err != nil {
		c.handleError(db, e, err.Error())
		return
	}
	var kubernetesClient kubernetes.Interface

	// Get Clients
	switch clusterData.Datacenter.Datacenter {
	case model.GCP:
		var containerClient *container.ClusterManagerClient
		kubernetesClient, containerClient, err = c.getAllGCPClient(ctx, clusterData)
		if err != nil {
			c.handleError(db, e, err.Error())
			return
		}
	}

	modifiedHPAs, err := c.scheduledHPAConfigUC.ListScheduledHPAConfigByEventID(db, e.ID)
	if err != nil {
		c.handleError(db, e, err.Error())
		return
	}

	existingHPA, err := c.clusterUC.GetAllHPAInCluster(
		ctx,
		kubernetesClient,
		clusterID,
		clusterData.LatestHPAAPIVersion,
	)
	if err != nil {
		c.handleError(db, e, err.Error())
		return
	}

	anyHPAExist := false
	for _, modifiedHPA := range modifiedHPAs {
		hpaExist := false
		for _, hpa := range existingHPA {
			if modifiedHPA.Name == hpa.Name && modifiedHPA.Namespace == hpa.Namespace {
				hpaExist = true
				break
			}
		}
		anyHPAExist = hpaExist || anyHPAExist
		if !hpaExist {
			err := c.scheduledHPAConfigUC.UpdateScheduledHPAConfigStatusMessage(
				db,
				modifiedHPA.ID,
				model.HPAUpdateFailed,
				"hpa not found",
			)
			if err != nil {
				log.Errorf(
					"[EventCronJob] Error Update HPA %s Namespace %s : %s",
					modifiedHPA.Name,
					modifiedHPA.Namespace,
					err.Error(),
				)
			}
		}
	}

	if !anyHPAExist {
		c.handleError(db, e, "no hpa exist")
		return
	}

}

func (c *cron) Start() {
	for {
		ctx := context.Background()
		db := c.tx.WithContext(ctx)

		now := time.Now()
		events, err := c.eventUC.GetAllPendingExecutableEvent(db, now)
		if err != nil {
			log.Errorf("[EventCronJob] Error Getting Events : %s", err.Error())
			continue
		}
		if len(events) == 0 {
			continue
		}
		for _, event := range events {
			go c.execEvent(event, db, ctx)
		}
	}
}
