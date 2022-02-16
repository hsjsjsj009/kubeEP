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
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/contactcenterinsights/v1"
	"gorm.io/gorm"
	"time"
)

type Datacenter interface {
	SaveDatacenterDetailedData(tx *gorm.DB, data *gcpUCEntity.DatacenterDetailedData) (uuid.UUID, error)
	SaveDatacenter(tx *gorm.DB, data gcpUCEntity.DatacenterData, SACredentials *gcpUCEntity.SAKeyCredentials) (uuid.UUID, error)
	ParseServiceAccountKey(data gcpUCEntity.DatacenterData) (*gcpUCEntity.SAKeyCredentials, error)
	GetGoogleCredentials(ctx context.Context, data gcpUCEntity.DatacenterData) (*google.Credentials, error)
	SaveTemporaryDatacenter(ctx context.Context, data gcpUCEntity.DatacenterData, SACredentials *gcpUCEntity.SAKeyCredentials) (uuid.UUID, error)
	GetTemporaryDatacenterData(ctx context.Context, id uuid.UUID) (*gcpUCEntity.DatacenterDetailedData, error)
	GetDatacenterData(tx *gorm.DB, id uuid.UUID) (*gcpUCEntity.DatacenterDetailedData, error)
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
		Name:       data.Name,
		Datacenter: constant.GCP,
	}
	datacenterModel.Credentials.SetRawMessage(data.Credentials)
	datacenterModel.Metadata.SetRawMessage(metaDataByte)
	err = d.datacenterRepo.InsertTemporaryDatacenter(ctx, datacenterModel, time.Hour)
	return datacenterModel.ID.GetUUID(), err
}

func (d *datacenter) GetTemporaryDatacenterData(ctx context.Context, id uuid.UUID) (*gcpUCEntity.DatacenterDetailedData, error) {
	data, err := d.datacenterRepo.GetTemporaryDatacenterByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &gcpUCEntity.DatacenterDetailedData{
		ID:          data.ID.GetUUID(),
		Name:        data.Name,
		Credentials: data.Credentials.GetRawMessage(),
		Metadata:    data.Metadata.GetRawMessage(),
		Datacenter:  data.Datacenter,
	}, nil
}

func (d *datacenter) GetDatacenterData(tx *gorm.DB, id uuid.UUID) (*gcpUCEntity.DatacenterDetailedData, error) {
	data, err := d.datacenterRepo.GetDatacenterByID(tx, id)
	if err != nil {
		return nil, err
	}
	return &gcpUCEntity.DatacenterDetailedData{
		ID:          data.ID.GetUUID(),
		Name:        data.Name,
		Credentials: data.Credentials.GetRawMessage(),
		Metadata:    data.Metadata.GetRawMessage(),
		Datacenter:  data.Datacenter,
	}, nil

}

func (d *datacenter) SaveDatacenterDetailedData(tx *gorm.DB, data *gcpUCEntity.DatacenterDetailedData) (uuid.UUID, error) {
	datacenterData := &model.Datacenter{
		Name:       data.Name,
		Datacenter: data.Datacenter,
	}
	datacenterData.ID.SetUUID(data.ID)
	datacenterData.Credentials.SetRawMessage(data.Credentials)
	datacenterData.Metadata.SetRawMessage(data.Metadata)
	err := d.datacenterRepo.InsertDatacenter(tx, datacenterData)
	return datacenterData.ID.GetUUID(), err
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
		Name:       data.Name,
		Datacenter: constant.GCP,
	}
	datacenterModel.Credentials.SetRawMessage(data.Credentials)
	datacenterModel.Metadata.SetRawMessage(metaDataByte)
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
