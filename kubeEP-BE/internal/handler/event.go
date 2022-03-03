package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/request"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
	"gorm.io/gorm"
	"time"
)

type Event interface {
	RegisterEvents(c *fiber.Ctx) error
	ListEventByCluster(c *fiber.Ctx) error
	UpdateEvent(c *fiber.Ctx) error
}

type event struct {
	kubernetesBaseHandler
	validatorInst        *validator.Validate
	db                   *gorm.DB
	eventUC              useCase.Event
	scheduledHPAConfigUC useCase.ScheduledHPAConfig
	hpaConfigStatusUC    useCase.HPAConfigStatus
}

func newEventHandler(
	validatorInst *validator.Validate,
	eventUC useCase.Event,
	scheduledHPAConfigUC useCase.ScheduledHPAConfig,
	hpaConfigStatusUC useCase.HPAConfigStatus,
	db *gorm.DB,
	kubeHandler kubernetesBaseHandler,
) Event {
	return &event{
		kubernetesBaseHandler: kubeHandler,
		validatorInst:         validatorInst,
		eventUC:               eventUC,
		scheduledHPAConfigUC:  scheduledHPAConfigUC,
		hpaConfigStatusUC:     hpaConfigStatusUC,
		db:                    db,
	}
}

func (e *event) RegisterEvents(c *fiber.Ctx) error {
	reqData := &request.EventDataRequest{}

	err := c.BodyParser(reqData)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}
	err = e.validatorInst.Struct(reqData)
	if err != nil {
		return e.errorResponse(c, errorConstant.InvalidRequestBody)
	}

	utcNow := time.Now().UTC()
	if utcNow.After(*reqData.StartTime) || utcNow.After(*reqData.EndTime) {
		return e.errorResponse(c, errorConstant.InvalidRequestBody)
	}

	ctx := c.Context()
	db := e.db.WithContext(ctx)
	tx := db.Begin()

	existingCluster, err := e.eventUC.GetEventByName(db, *reqData.Name)
	if existingCluster != nil {
		return e.errorResponse(c, errorConstant.EventExist)
	}

	kubernetesClient, latestHPAAPIVersion, err := e.getClusterKubernetesClient(
		ctx,
		db,
		*reqData.ClusterID,
	)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	HPAs, err := e.generalClusterUC.GetAllHPAInCluster(
		ctx,
		kubernetesClient,
		*reqData.ClusterID,
		latestHPAAPIVersion,
	)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	eventData := &UCEntity.Event{
		Name:      *reqData.Name,
		StartTime: *reqData.StartTime,
		EndTime:   *reqData.EndTime,
	}
	eventData.Cluster.ID = *reqData.ClusterID

	eventID, err := e.eventUC.RegisterEvents(tx, eventData)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	var HPAConfigs []UCEntity.EventModifiedHPAConfigData
	for _, hpaConfig := range reqData.ModifiedHPAConfigs {
		found := false
		for _, HPA := range HPAs {
			if *hpaConfig.Name == HPA.Name && *hpaConfig.Namespace == HPA.Namespace {
				found = true
				break
			}
		}
		if found {
			HPAConfigs = append(
				HPAConfigs, UCEntity.EventModifiedHPAConfigData{
					Name:        *hpaConfig.Name,
					Namespace:   *hpaConfig.Namespace,
					MinReplicas: hpaConfig.MinReplicas,
					MaxReplicas: *hpaConfig.MaxReplicas,
				},
			)
		}
	}

	hpaConfigIDs, err := e.scheduledHPAConfigUC.RegisterModifiedHPAConfigs(tx, HPAConfigs, eventID)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	_, err = e.hpaConfigStatusUC.CreateHPAConfigStatusForScheduledConfigIDs(tx, hpaConfigIDs)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	tx.Commit()

	return e.successResponse(c, response.EventCreationResponse{EventID: eventID})

}

func (e *event) ListEventByCluster(c *fiber.Ctx) error {
	reqData := &request.EventListRequest{}
	err := c.QueryParser(reqData)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	err = e.validatorInst.Struct(reqData)
	if err != nil {
		return e.errorResponse(c, errorConstant.InvalidQueryParam)
	}

	ctx := c.Context()
	tx := e.db.WithContext(ctx)

	events, err := e.eventUC.ListEventByClusterID(tx, *reqData.ClusterID)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	responseData := make([]response.EventSimpleResponse, 0)
	for _, event := range events {
		responseData = append(
			responseData, response.EventSimpleResponse{
				ID:        event.ID,
				Name:      event.Name,
				StartTime: event.StartTime,
				EndTime:   event.EndTime,
			},
		)
	}

	return e.successResponse(c, responseData)

}

func (e *event) UpdateEvent(c *fiber.Ctx) error {
	req := &request.UpdateEventDataRequest{}
	if err := c.BodyParser(req); err != nil {
		return e.errorResponse(c, err.Error())
	}

	if err := e.validatorInst.Struct(req); err != nil {
		return e.errorResponse(c, errorConstant.InvalidRequestBody)
	}

	ctx := c.Context()
	db := e.db.WithContext(ctx)
	tx := db.Begin()

	eventData, err := e.eventUC.GetEventByID(db, *req.EventID)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	if eventData.Name != *req.Name {
		existingCluster, err := e.eventUC.GetEventByName(db, *req.Name)
		if existingCluster != nil || err == nil {
			return e.errorResponse(c, errorConstant.EventExist)
		}
		eventData.Name = *req.Name
	}

	eventData.Cluster.ID = *req.ClusterID

	if err := e.eventUC.UpdateEvent(tx, eventData); err != nil {
		return e.errorResponse(c, err.Error())
	}

	if err := e.scheduledHPAConfigUC.DeleteEventModifiedHPAConfigs(tx, eventData.ID); err != nil {
		return e.errorResponse(c, err.Error())
	}

	var newModifiedHPAConfigs []UCEntity.EventModifiedHPAConfigData
	for _, hpaConfig := range req.ModifiedHPAConfigs {
		newModifiedHPAConfigs = append(
			newModifiedHPAConfigs, UCEntity.EventModifiedHPAConfigData{
				Name:        *hpaConfig.Name,
				Namespace:   *hpaConfig.Namespace,
				MinReplicas: hpaConfig.MinReplicas,
				MaxReplicas: *hpaConfig.MaxReplicas,
			},
		)
	}

	newHPAConfigIDs, err := e.scheduledHPAConfigUC.RegisterModifiedHPAConfigs(
		tx,
		newModifiedHPAConfigs,
		eventData.ID,
	)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	_, err = e.hpaConfigStatusUC.CreateHPAConfigStatusForScheduledConfigIDs(tx, newHPAConfigIDs)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}

	tx.Commit()

	res := &response.EventCreationResponse{EventID: eventData.ID}
	return e.successResponse(c, res)
}
