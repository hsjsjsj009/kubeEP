package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type Datacenter interface {
	GetDatacenterByID(tx *gorm.DB, id uuid.UUID) (*model.Datacenter, error)
	InsertDatacenter(tx *gorm.DB, data *model.Datacenter) error
}

type datacenter struct {
}

func NewDatacenter() Datacenter {
	return &datacenter{}
}

func (d *datacenter) GetDatacenterByID(tx *gorm.DB, id uuid.UUID) (*model.Datacenter, error) {
	data := &model.Datacenter{}
	tx = tx.First(data, id)
	return data, tx.Error
}

func (d *datacenter) InsertDatacenter(tx *gorm.DB, data *model.Datacenter) error {
	return tx.Create(data).Error
}
