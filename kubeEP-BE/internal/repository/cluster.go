package repository

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
)

type Cluster interface {
	GetClusterByID(id uuid.UUID) (*model.Cluster, error)
	ListClusterByDatacenterID(id uuid.UUID) ([]*model.Cluster, error)
	InsertCluster(data *model.Cluster) error
}

type cluster struct {
	db *gorm.DB
}

func NewCluster(db *gorm.DB) Cluster {
	return &cluster{db: db}
}

func (d *cluster) GetClusterByID(id uuid.UUID) (*model.Cluster, error) {
	data := &model.Cluster{}
	tx := d.db.Model(data).First(data, id)
	return data, tx.Error
}

func (d *cluster) ListClusterByDatacenterID(id uuid.UUID) ([]*model.Cluster, error) {
	var data []*model.Cluster
	tx := d.db.Model(&model.Cluster{}).Where("datacenter_id = ?", id).Scan(&data)
	return data, tx.Error
}

func (d cluster) InsertCluster(data *model.Cluster) error {
	return d.db.Create(data).Error
}
