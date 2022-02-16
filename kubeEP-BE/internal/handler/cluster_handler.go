package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
	"gorm.io/gorm"
)

type ClusterHandler interface {
	GetAllClusters(c *fiber.Ctx) error
}

type clusterHandler struct {
	baseHandler
	validatorInst    *validator.Validate
	generalClusterUC useCase.Cluster
	db               *gorm.DB
}

func newClusterHandler(validatorInst *validator.Validate, generalClusterUC useCase.Cluster, db *gorm.DB) ClusterHandler {
	return &clusterHandler{validatorInst: validatorInst, generalClusterUC: generalClusterUC, db: db}
}

func (ch *clusterHandler) GetAllClusters(c *fiber.Ctx) error {
	existingClusters, err := ch.generalClusterUC.GetAllClustersInLocal(ch.db)
	if err != nil {
		return ch.errorResponse(c, err.Error())
	}
	var responses []response.Cluster
	for _, cluster := range existingClusters {
		responses = append(responses, response.Cluster{
			ID:         &cluster.ID,
			Name:       cluster.Name,
			Datacenter: constant.GCP,
		})
	}
	return ch.successResponse(c, responses)
}
