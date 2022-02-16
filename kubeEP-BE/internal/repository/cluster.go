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
	InsertClusterBatch(tx *gorm.DB, data []*model.Cluster) error
	ListAllRegisteredCluster(tx *gorm.DB) ([]*model.Cluster, error)
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

func (d *cluster) ListAllRegisteredCluster(tx *gorm.DB) ([]*model.Cluster, error) {
	var data []*model.Cluster
	rows, err := tx.Raw(`
		SELECT 
		       c.id, 
		       c.datacenter_id, 
		       c.metadata, 
		       c.name, 
		       c.certificate, 
		       c.server_endpoint,
		       d.datacenter
		from clusters c
		join datacenters d on d.id = c.datacenter_id and d.deleted_at is null
		where c.datacenter_id = ? and c.deleted_at is null
	`).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		cluster := &model.Cluster{}
		err = rows.Scan(
			&cluster.ID,
			&cluster.DatacenterID,
			&cluster.Metadata,
			&cluster.Name,
			&cluster.Certificate,
			&cluster.ServerEndpoint,
			&cluster.Datacenter.Datacenter,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, cluster)
	}
	return data, tx.Error
}

func (d cluster) InsertCluster(tx *gorm.DB, data *model.Cluster) error {
	return tx.Create(data).Error
}

func (d cluster) InsertClusterBatch(tx *gorm.DB, data []*model.Cluster) error {
	return tx.Create(&data).Error
}
