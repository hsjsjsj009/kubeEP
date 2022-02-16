package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	gcpRequest "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/request/gcp"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response"
	gcpResponse "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response/gcp"
	gcpUCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase/gcp"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
	gcpUseCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase/gcp"
	"gorm.io/gorm"
)

type GcpHandler interface {
	RegisterDatacenter(c *fiber.Ctx) error
	GetClustersByDatacenterID(c *fiber.Ctx) error
	RegisterClusterWithDatacenter(c *fiber.Ctx) error
}

type gcpHandler struct {
	baseHandler
	validatorInst    *validator.Validate
	clusterUC        gcpUseCase.Cluster
	generalClusterUC useCase.Cluster
	datacenterUC     gcpUseCase.Datacenter
	db               *gorm.DB
}

func newGCPHandler(validatorInst *validator.Validate, clusterUC gcpUseCase.Cluster, datacenterUC gcpUseCase.Datacenter, db *gorm.DB, generalClusterUC useCase.Cluster) GcpHandler {

	return &gcpHandler{
		validatorInst:    validatorInst,
		clusterUC:        clusterUC,
		datacenterUC:     datacenterUC,
		generalClusterUC: generalClusterUC,
		db:               db,
	}
}

func (g *gcpHandler) RegisterDatacenter(c *fiber.Ctx) error {
	reqData := &gcpRequest.DatacenterData{}
	err := c.BodyParser(reqData)
	if err != nil {
		return g.errorResponse(c, errorConstant.InvalidRequestBody)
	}
	err = g.validatorInst.Struct(reqData)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	datacenterData := gcpUCEntity.DatacenterData{
		Credentials: *reqData.SAKeyCredentials,
		Name:        *reqData.Name,
	}
	SAData, err := g.datacenterUC.ParseServiceAccountKey(datacenterData)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	var id uuid.UUID
	if *reqData.IsTemporary {
		id, err = g.datacenterUC.SaveTemporaryDatacenter(c.Context(), datacenterData, SAData)
	} else {
		id, err = g.datacenterUC.SaveDatacenter(g.db, datacenterData, SAData)
	}

	return g.successResponse(c, gcpResponse.DatacenterData{DatacenterID: id, IsTemporary: *reqData.IsTemporary})
}

func (g *gcpHandler) GetClustersByDatacenterID(c *fiber.Ctx) error {
	reqData := &gcpRequest.ExistingDatacenterData{}
	err := c.QueryParser(reqData)
	if err != nil {
		return g.errorResponse(c, errorConstant.InvalidQueryParam)
	}
	err = g.validatorInst.Struct(reqData)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	isTemporaryDatacenter := true
	data, err := g.datacenterUC.GetTemporaryDatacenterData(c.Context(), *reqData.DatacenterID)
	if err != nil {
		isTemporaryDatacenter = false
		data, err = g.datacenterUC.GetDatacenterData(g.db, *reqData.DatacenterID)
		if err != nil {
			return g.errorResponse(c, err.Error())
		}
	}
	datacenterData := gcpUCEntity.DatacenterData{
		Credentials: data.Credentials,
		Name:        data.Name,
	}
	googleCredentials, err := g.datacenterUC.GetGoogleCredentials(c.Context(), datacenterData)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	clusterClient, err := g.clusterUC.GetGoogleClusterClient(c.Context(), googleCredentials)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	clusters, err := g.clusterUC.GetAllClustersInGCPProject(c.Context(), googleCredentials.ProjectID, clusterClient)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}

	var clusterData []gcpResponse.Cluster
	for _, cluster := range clusters {
		clusterData = append(clusterData, gcpResponse.Cluster{
			Cluster: response.Cluster{
				Name:       cluster.Name,
				Datacenter: constant.GCP,
			},
			Location: cluster.Location,
		})
	}

	return g.successResponse(c, gcpResponse.DatacenterClusters{
		Clusters:              clusterData,
		IsTemporaryDatacenter: isTemporaryDatacenter,
	})
}

func (g *gcpHandler) RegisterClusterWithDatacenter(c *fiber.Ctx) error {
	reqData := &gcpRequest.RegisterClusterData{}
	err := c.BodyParser(reqData)
	if err != nil {
		return g.errorResponse(c, errorConstant.InvalidRequestBody)
	}
	err = g.validatorInst.Struct(reqData)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	var data *gcpUCEntity.DatacenterDetailedData
	if *reqData.IsDatacenterTemporary {
		data, err = g.datacenterUC.GetTemporaryDatacenterData(c.Context(), *reqData.DatacenterID)
	} else {
		data, err = g.datacenterUC.GetDatacenterData(g.db, *reqData.DatacenterID)
	}
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	datacenterData := gcpUCEntity.DatacenterData{
		Credentials: data.Credentials,
		Name:        data.Name,
	}
	googleCredentials, err := g.datacenterUC.GetGoogleCredentials(c.Context(), datacenterData)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	clusterClient, err := g.clusterUC.GetGoogleClusterClient(c.Context(), googleCredentials)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	clusters, err := g.clusterUC.GetAllClustersInGCPProject(c.Context(), googleCredentials.ProjectID, clusterClient)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}

	existingCluster, err := g.generalClusterUC.GetAllClustersInLocalByDatacenterID(g.db, *reqData.DatacenterID)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}

	var selectedClusters []*gcpUCEntity.ClusterData
	for _, clusterName := range reqData.ClustersName {
		for _, cluster := range existingCluster {
			if cluster.Name == clusterName {
				return g.errorResponse(c, fmt.Sprintf(errorConstant.ClusterExists, clusterName))
			}
		}

		contains := false
		for _, cluster := range clusters {
			if cluster.Name == clusterName {
				selectedClusters = append(selectedClusters, cluster)
				contains = true
				break
			}
		}
		if !contains {
			return g.errorResponse(c, fmt.Sprintf(errorConstant.ClusterNotFound, clusterName))
		}
	}

	tx := g.db.Begin()

	if *reqData.IsDatacenterTemporary {
		_, err = g.datacenterUC.SaveDatacenterDetailedData(tx, data)
		if err != nil {
			return g.errorResponse(c, err.Error())
		}
	}

	err = g.clusterUC.RegisterClusters(tx, *reqData.DatacenterID, selectedClusters)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}

	tx.Commit()

	var responses []gcpResponse.Cluster
	for _, cluster := range selectedClusters {
		responses = append(responses, gcpResponse.Cluster{
			Cluster: response.Cluster{
				ID:         &cluster.ID,
				Name:       cluster.Name,
				Datacenter: constant.GCP,
			},
			Location: cluster.Location,
		})
	}

	return g.successResponse(c, responses)
}
