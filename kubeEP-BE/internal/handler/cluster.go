package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/request"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
	"gorm.io/gorm"
)

type Cluster interface {
	GetAllRegisteredClusters(c *fiber.Ctx) error
	GetAllHPA(c *fiber.Ctx) error
}

type cluster struct {
	kubernetesBaseHandler
	validatorInst       *validator.Validate
	generalDatacenterUC useCase.Datacenter
	db                  *gorm.DB
}

func newClusterHandler(
	validatorInst *validator.Validate,
	db *gorm.DB,
	generalDatacenterUC useCase.Datacenter,
	kubeHandler kubernetesBaseHandler,
) Cluster {
	return &cluster{
		kubernetesBaseHandler: kubeHandler,
		validatorInst:         validatorInst,
		db:                    db,
		generalDatacenterUC:   generalDatacenterUC,
	}
}

func (ch *cluster) GetAllRegisteredClusters(c *fiber.Ctx) error {

	ctx := c.Context()
	tx := ch.db.WithContext(ctx)

	existingClusters, err := ch.generalClusterUC.GetAllClustersInLocal(tx)
	if err != nil {
		return ch.errorResponse(c, err.Error())
	}
	responses := make([]response.Cluster, 0)
	for _, cluster := range existingClusters {
		responses = append(
			responses, response.Cluster{
				ID:             &cluster.ID,
				Name:           cluster.Name,
				Datacenter:     cluster.Datacenter.Datacenter,
				DatacenterName: cluster.Datacenter.Name,
			},
		)
	}
	return ch.successResponse(c, responses)
}

func (ch *cluster) GetAllHPA(c *fiber.Ctx) error {
	requestData := &request.ExistingClusterData{}
	err := c.QueryParser(requestData)
	if err != nil {
		return ch.errorResponse(c, errorConstant.InvalidQueryParam)
	}
	err = ch.validatorInst.Struct(requestData)
	if err != nil {
		return ch.errorResponse(c, errorConstant.InvalidQueryParam)
	}

	ctx := c.Context()

	tx := ch.db.WithContext(ctx)

	kubernetesClient, latestHPAAPIVersion, err := ch.getClusterKubernetesClient(
		ctx,
		tx,
		*requestData.ClusterID,
	)
	if err != nil {
		return ch.errorResponse(c, err.Error())
	}

	HPAs, err := ch.generalClusterUC.GetAllHPAInCluster(
		ctx,
		kubernetesClient,
		*requestData.ClusterID,
		latestHPAAPIVersion,
	)
	if err != nil {
		return ch.errorResponse(c, err.Error())
	}
	responses := make([]response.SimpleHPA, 0)
	for _, hpa := range HPAs {
		responses = append(
			responses, response.SimpleHPA{
				Name:            hpa.Name,
				Namespace:       hpa.Namespace,
				MinReplicas:     hpa.MinReplicas,
				MaxReplicas:     hpa.MaxReplicas,
				CurrentReplicas: hpa.CurrentReplicas,
			},
		)
	}

	return ch.successResponse(c, responses)
}
