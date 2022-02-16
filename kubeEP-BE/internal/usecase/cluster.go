package useCase

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"gorm.io/gorm"
)

type Cluster interface {
	GetAllClustersInLocalByDatacenterID(tx *gorm.DB, datacenterID uuid.UUID) ([]UCEntity.ClusterData, error)
}

type cluster struct {
	validatorInst *validator.Validate
	clusterRepo   repository.Cluster
}

func NewCluster(validatorInst *validator.Validate, clusterRepo repository.Cluster) Cluster {
	return &cluster{validatorInst: validatorInst, clusterRepo: clusterRepo}
}

func (c *cluster) GetAllClustersInLocalByDatacenterID(tx *gorm.DB, datacenterID uuid.UUID) ([]UCEntity.ClusterData, error) {
	clusters, err := c.clusterRepo.ListClusterByDatacenterID(tx, datacenterID)
	if err != nil {
		return nil, err
	}
	var output []UCEntity.ClusterData
	for _, cluster := range clusters {
		output = append(output, UCEntity.ClusterData{
			ID:             cluster.ID.GetUUID(),
			Name:           cluster.Name,
			Certificate:    cluster.Certificate,
			ServerEndpoint: cluster.ServerEndpoint,
			Datacenter:     cluster.Datacenter.Datacenter,
		})
	}
	return output, nil
}
