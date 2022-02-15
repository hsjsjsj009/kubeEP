package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	gcpRequest "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/request/gcp"
	gcpResponse "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response/gcp"
	gcpUCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase/gcp"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	gcpUseCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase/gcp"
	"gorm.io/gorm"
)

type GcpHandler interface {
	RegisterDatacenter(c *fiber.Ctx) error
	GetClustersByDatacenterID(c *fiber.Ctx) error
}

type gcpHandler struct {
	baseHandler
	validatorInst *validator.Validate
	clusterUC     gcpUseCase.Cluster
	datacenterUC  gcpUseCase.Datacenter
	db            *gorm.DB
}

func newGCPHandler(validatorInst *validator.Validate, clusterUC gcpUseCase.Cluster, datacenterUC gcpUseCase.Datacenter, db *gorm.DB) GcpHandler {

	return &gcpHandler{
		validatorInst: validatorInst,
		clusterUC:     clusterUC,
		datacenterUC:  datacenterUC,
		db:            db,
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
	err := c.BodyParser(reqData)
	if err != nil {
		return g.errorResponse(c, errorConstant.InvalidRequestBody)
	}
	err = g.validatorInst.Struct(reqData)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}
	var data *model.Datacenter
	isTemporaryDatacenter := true
	data, err = g.datacenterUC.GetTemporaryDatacenterData(c.Context(), *reqData.DatacenterID)
	if err != nil {
		isTemporaryDatacenter = false
		data, err = g.datacenterUC.GetDatacenterData(g.db, *reqData.DatacenterID)
		if err != nil {
			return g.errorResponse(c, err.Error())
		}
	}
	datacenterData := gcpUCEntity.DatacenterData{
		Credentials: json.RawMessage(data.Credentials),
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
	clusters, err := g.clusterUC.GetAllClustersInProject(c.Context(), googleCredentials.ProjectID, clusterClient)
	if err != nil {
		return g.errorResponse(c, err.Error())
	}

	var clusterData []gcpResponse.Cluster
	for _, cluster := range clusters {
		clusterData = append(clusterData, gcpResponse.Cluster{
			Name:     cluster.Name,
			Location: cluster.Location,
		})
	}

	return g.successResponse(c, gcpResponse.DatacenterClusters{
		Clusters:              clusterData,
		IsTemporaryDatacenter: isTemporaryDatacenter,
	})
}
