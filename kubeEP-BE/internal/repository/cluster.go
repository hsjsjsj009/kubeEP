package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type Cluster interface {
	GetClusterByID(tx *gorm.DB, id uuid.UUID) (*model.Cluster, error)
	ListClusterByDatacenterID(tx *gorm.DB, id uuid.UUID) ([]*model.Cluster, error)
	InsertCluster(tx *gorm.DB, data *model.Cluster) error
}

type cluster struct {
}

func NewCluster() Cluster {
	return &cluster{}
}

func (d *cluster) GetClusterByID(tx *gorm.DB, id uuid.UUID) (*model.Cluster, error) {
	data := &model.Cluster{}
	tx = tx.Model(data).First(data, id)
	return data, tx.Error
}

func (d *cluster) ListClusterByDatacenterID(tx *gorm.DB, id uuid.UUID) ([]*model.Cluster, error) {
	var data []*model.Cluster
	tx = tx.Model(&model.Cluster{}).Where("datacenter_id = ?", id).Find(&data)
	return data, tx.Error
}

func (d cluster) InsertCluster(tx *gorm.DB, data *model.Cluster) error {
	return tx.Create(data).Error
}
