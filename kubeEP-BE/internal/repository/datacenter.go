package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type Datacenter interface {
	GetDatacenterByID(id uuid.UUID) (*model.Datacenter, error)
	InsertDatacenter(data *model.Datacenter) error
}

type datacenter struct {
	db *gorm.DB
}

func NewDatacenter(db *gorm.DB) Datacenter {
	return &datacenter{
		db: db,
	}
}

func (d *datacenter) GetDatacenterByID(id uuid.UUID) (*model.Datacenter, error) {
	data := &model.Datacenter{}
	tx := d.db.First(data, id)
	return data, tx.Error
}

func (d *datacenter) InsertDatacenter(data *model.Datacenter) error {
	return d.db.Create(data).Error
}
