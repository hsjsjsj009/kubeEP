package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/request"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
	"gorm.io/gorm"
)

type Event interface {
	RegisterEvents(c *fiber.Ctx) error
}

type event struct {
	baseHandler
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
) Event {
	return &event{
		validatorInst:        validatorInst,
		eventUC:              eventUC,
		scheduledHPAConfigUC: scheduledHPAConfigUC,
		hpaConfigStatusUC:    hpaConfigStatusUC,
		db:                   db,
	}
}

func (e *event) RegisterEvents(c *fiber.Ctx) error {
	reqData := &request.EventRequest{}

	err := c.BodyParser(reqData)
	if err != nil {
		return e.errorResponse(c, err.Error())
	}
	err = e.validatorInst.Struct(reqData)
	if err != nil {
		return e.errorResponse(c, errors.New(errorConstant.InvalidRequestBody))
	}

	ctx := c.Context()
	tx := e.db.WithContext(ctx)
	tx = tx.Begin()

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
		HPAConfigs = append(
			HPAConfigs, UCEntity.EventModifiedHPAConfigData{
				Name:        *hpaConfig.Name,
				Namespace:   *hpaConfig.Namespace,
				MinReplicas: hpaConfig.MinReplicas,
				MaxReplicas: *hpaConfig.MaxReplicas,
			},
		)
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
