package gcpUseCase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	gcpUCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase/gcp"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/contactcenterinsights/v1"
	"gorm.io/gorm"
	"time"
)

type Datacenter interface {
	SaveDatacenter(tx *gorm.DB, data gcpUCEntity.DatacenterData, SACredentials *gcpUCEntity.SAKeyCredentials) (uuid.UUID, error)
	ParseServiceAccountKey(data gcpUCEntity.DatacenterData) (*gcpUCEntity.SAKeyCredentials, error)
	GetGoogleCredentials(ctx context.Context, data gcpUCEntity.DatacenterData) (*google.Credentials, error)
	SaveTemporaryDatacenter(ctx context.Context, data gcpUCEntity.DatacenterData, SACredentials *gcpUCEntity.SAKeyCredentials) (uuid.UUID, error)
	GetTemporaryDatacenterData(ctx context.Context, id uuid.UUID) (*model.Datacenter, error)
	GetDatacenterData(tx *gorm.DB, id uuid.UUID) (*model.Datacenter, error)
}

type datacenter struct {
	datacenterRepo repository.Datacenter
	validatorInst  *validator.Validate
}

func NewDatacenter(datacenterRepo repository.Datacenter, validatorInst *validator.Validate) Datacenter {
	return &datacenter{
		datacenterRepo: datacenterRepo,
		validatorInst:  validatorInst,
	}
}

func (d *datacenter) ParseServiceAccountKey(data gcpUCEntity.DatacenterData) (*gcpUCEntity.SAKeyCredentials, error) {
	SACredentials := &gcpUCEntity.SAKeyCredentials{}
	err := json.Unmarshal(data.Credentials, SACredentials)
	if err != nil {
		return nil, err
	}
	err = d.validatorInst.Struct(SACredentials)
	if err != nil {
		return nil, errors.New(errorConstant.SAKeyInvalid)
	}
	return SACredentials, nil
}

func (d *datacenter) SaveTemporaryDatacenter(ctx context.Context, data gcpUCEntity.DatacenterData, SACredentials *gcpUCEntity.SAKeyCredentials) (uuid.UUID, error) {
	metaData := &gcpUCEntity.DatacenterMetaData{
		ProjectId: *SACredentials.ProjectId,
		SAEmail:   *SACredentials.ClientEmail,
	}
	metaDataByte, err := json.Marshal(metaData)
	if err != nil {
		return uuid.UUID{}, err
	}
	datacenterModel := &model.Datacenter{
		Name:        data.Name,
		Credentials: gormDatatype.JSON(data.Credentials),
		Metadata:    gormDatatype.JSON(metaDataByte),
		Datacenter:  constant.GCP,
	}
	err = d.datacenterRepo.InsertTemporaryDatacenter(ctx, datacenterModel, time.Hour)
	return uuid.UUID(datacenterModel.ID), err
}

func (d *datacenter) GetTemporaryDatacenterData(ctx context.Context, id uuid.UUID) (*model.Datacenter, error) {
	return d.datacenterRepo.GetTemporaryDatacenterByID(ctx, id)
}

func (d *datacenter) GetDatacenterData(tx *gorm.DB, id uuid.UUID) (*model.Datacenter, error) {
	return d.datacenterRepo.GetDatacenterByID(tx, id)
}

func (d *datacenter) SaveDatacenter(tx *gorm.DB, data gcpUCEntity.DatacenterData, SACredentials *gcpUCEntity.SAKeyCredentials) (uuid.UUID, error) {
	metaData := &gcpUCEntity.DatacenterMetaData{
		ProjectId: *SACredentials.ProjectId,
		SAEmail:   *SACredentials.ClientEmail,
	}
	metaDataByte, err := json.Marshal(metaData)
	if err != nil {
		return uuid.UUID{}, err
	}
	datacenterModel := model.Datacenter{
		Name:        data.Name,
		Credentials: gormDatatype.JSON(data.Credentials),
		Metadata:    gormDatatype.JSON(metaDataByte),
		Datacenter:  constant.GCP,
	}
	err = d.datacenterRepo.InsertDatacenter(tx, &datacenterModel)
	return uuid.UUID(datacenterModel.ID), err
}

func (d *datacenter) GetGoogleCredentials(ctx context.Context, data gcpUCEntity.DatacenterData) (*google.Credentials, error) {
	credentials, err := google.CredentialsFromJSON(ctx, data.Credentials, contactcenterinsights.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}
