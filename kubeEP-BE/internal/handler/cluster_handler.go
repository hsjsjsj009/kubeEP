package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/request"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
)

type ClusterHandler interface {
	GetAllRegisteredClusters(c *fiber.Ctx) error
	GetAllHPA(c *fiber.Ctx) error
}

type clusterHandler struct {
	baseHandler
	validatorInst       *validator.Validate
	generalClusterUC    useCase.Cluster
	gcpClusterUC        useCase.GCPCluster
	gcpDatacenterUC     useCase.GCPDatacenter
	generalDatacenterUC useCase.Datacenter
	db                  *gorm.DB
}

func newClusterHandler(
	validatorInst *validator.Validate,
	generalClusterUC useCase.Cluster,
	db *gorm.DB,
	gcpClusterUC useCase.GCPCluster,
	gcpDatacenterUC useCase.GCPDatacenter,
	generalDatacenterUC useCase.Datacenter,
) ClusterHandler {
	return &clusterHandler{
		validatorInst:       validatorInst,
		generalClusterUC:    generalClusterUC,
		db:                  db,
		gcpClusterUC:        gcpClusterUC,
		gcpDatacenterUC:     gcpDatacenterUC,
		generalDatacenterUC: generalDatacenterUC,
	}
}

func (ch *clusterHandler) GetAllRegisteredClusters(c *fiber.Ctx) error {
	existingClusters, err := ch.generalClusterUC.GetAllClustersInLocal(ch.db)
	if err != nil {
		return ch.errorResponse(c, err.Error())
	}
	var responses []response.Cluster
	for _, cluster := range existingClusters {
		responses = append(
			responses, response.Cluster{
				ID:         &cluster.ID,
				Name:       cluster.Name,
				Datacenter: model.GCP,
			},
		)
	}
	return ch.successResponse(c, responses)
}

func (ch *clusterHandler) GetAllHPA(c *fiber.Ctx) error {
	requestData := &request.ExistingClusterData{}
	err := c.QueryParser(requestData)
	if err != nil {
		return ch.errorResponse(c, errors.New(errorConstant.InvalidQueryParam))
	}
	err = ch.validatorInst.Struct(requestData)
	if err != nil {
		return ch.errorResponse(c, errors.New(errorConstant.InvalidQueryParam))
	}
	clusterData, err := ch.generalClusterUC.GetClusterAndDatacenterDataByClusterID(
		ch.db,
		*requestData.ClusterID,
	)
	if err != nil {
		return ch.errorResponse(c, err.Error())
	}
	var kubernetesClient kubernetes.Interface
	switch clusterData.Datacenter.Datacenter {
	case model.GCP:
		datacenterName := clusterData.Datacenter.Name
		datacenterData := UCEntity.DatacenterData{
			Credentials: clusterData.Datacenter.Credentials,
			Name:        datacenterName,
		}
		googleCredential, err := ch.gcpDatacenterUC.GetGoogleCredentials(
			c.Context(),
			datacenterData,
		)
		if err != nil {
			return ch.errorResponse(c, err.Error())
		}
		ch.gcpClusterUC.RegisterGoogleCredentials(datacenterName, googleCredential)
		kubernetesClient, err = ch.gcpClusterUC.GetKubernetesClusterClient(
			datacenterName,
			clusterData,
		)
		if err != nil {
			return ch.errorResponse(c, err.Error())
		}
	default:
		return ch.errorResponse(c, errors.New(errorConstant.DatacenterTypeNotFound))
	}

	HPAs, err := ch.generalClusterUC.GetAllHPAInCluster(
		c.Context(),
		kubernetesClient,
		*requestData.ClusterID,
	)
	if err != nil {
		return ch.errorResponse(c, err.Error())
	}
	var responses []response.SimpleHPA
	for _, hpa := range HPAs {
		responses = append(
			responses, response.SimpleHPA{
				APIVersion:      hpa.APIVersion,
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
