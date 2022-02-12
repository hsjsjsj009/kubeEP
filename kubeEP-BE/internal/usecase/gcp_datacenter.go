package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	useCaseEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/contactcenterinsights/v1"
	"gorm.io/gorm"
)

type Datacenter interface {
	RegisterDatacenter(db *gorm.DB, data useCaseEntity.DatacenterData) error
}

type datacenter struct {
	datacenterRepo repository.Datacenter
}

func NewDatacenter(datacenterRepo repository.Datacenter) Datacenter {
	return &datacenter{
		datacenterRepo: datacenterRepo,
	}
}

func (d *datacenter) RegisterDatacenter(tx *gorm.DB, data useCaseEntity.DatacenterData) error {
	metaDataByte, err := json.Marshal(data.Metadata)
	if err != nil {
		return err
	}
	datacenterModel := model.Datacenter{
		Name:        data.Name,
		Credentials: gormDatatype.JSON(data.Credentials),
		Metadata:    gormDatatype.JSON(metaDataByte),
		Datacenter:  data.Datacenter,
	}
	err = d.datacenterRepo.InsertDatacenter(tx, &datacenterModel)
	return err
}

func (d *datacenter) getGCPClient(ctx context.Context, data useCaseEntity.DatacenterData) (*google.Credentials, error) {
	if data.Datacenter != constant.GCP {
		return nil, errors.New(errorConstant.DatacenterMismatch)
	}
	credentials, err := google.CredentialsFromJSON(ctx, data.Credentials, contactcenterinsights.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}
