package useCase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"gorm.io/gorm"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

type Cluster interface {
	GetAllClustersInLocalByDatacenterID(tx *gorm.DB, datacenterID uuid.UUID) ([]UCEntity.ClusterData, error)
	GetAllClustersInLocal(tx *gorm.DB) ([]UCEntity.ClusterData, error)
	GetAllHPAInCluster(ctx context.Context, client kubernetes.Interface, clusterID uuid.UUID) (output []UCEntity.SimpleHPAData, err error)
	GetClusterAndDatacenterDataByClusterID(tx *gorm.DB, id uuid.UUID) (*UCEntity.ClusterData, error)
}

type cluster struct {
	validatorInst *validator.Validate
	clusterRepo   repository.Cluster
	hpaRepo       repository.K8sHPA
	namespaceRepo repository.K8sNamespace
}

func newCluster(validatorInst *validator.Validate, clusterRepo repository.Cluster, hpaRepo repository.K8sHPA, namespaceRepo repository.K8sNamespace) Cluster {
	return &cluster{validatorInst: validatorInst, clusterRepo: clusterRepo, hpaRepo: hpaRepo, namespaceRepo: namespaceRepo}
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
			Datacenter: UCEntity.DatacenterDetailedData{
				Datacenter: cluster.Datacenter.Datacenter,
			},
		})
	}
	return output, nil
}

func (c *cluster) GetAllClustersInLocal(tx *gorm.DB) ([]UCEntity.ClusterData, error) {
	clusters, err := c.clusterRepo.ListAllRegisteredCluster(tx)
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
			Datacenter: UCEntity.DatacenterDetailedData{
				Datacenter: cluster.Datacenter.Datacenter,
			},
		})
	}
	return output, nil
}

func (c *cluster) GetClusterAndDatacenterDataByClusterID(tx *gorm.DB, id uuid.UUID) (*UCEntity.ClusterData, error) {
	data, err := c.clusterRepo.GetClusterWithDatacenterByID(tx, id)
	if err != nil {
		return nil, err
	}
	datacenterModelData := data.Datacenter
	clusterData := &UCEntity.ClusterData{
		ID:             data.ID.GetUUID(),
		Name:           data.Name,
		Certificate:    data.Certificate,
		ServerEndpoint: data.ServerEndpoint,
		Datacenter: UCEntity.DatacenterDetailedData{
			ID:          datacenterModelData.ID.GetUUID(),
			Name:        datacenterModelData.Name,
			Credentials: datacenterModelData.Credentials.GetRawMessage(),
			Metadata:    datacenterModelData.Metadata.GetRawMessage(),
			Datacenter:  datacenterModelData.Datacenter,
		},
	}
	return clusterData, nil
}

func (c *cluster) GetAllHPAInCluster(ctx context.Context, client kubernetes.Interface, clusterID uuid.UUID) (output []UCEntity.SimpleHPAData, err error) {
	namespaces, err := c.namespaceRepo.GetAllNamespace(ctx, client)
	if err != nil {
		return nil, err
	}
	var lock sync.Mutex
	var wg sync.WaitGroup
	for _, namespace := range namespaces {
		wg.Add(2)
		go func(ns v1.Namespace) {
			defer wg.Done()
			v1HPAs, err := c.hpaRepo.GetAllV1HPA(ctx, client, ns, clusterID)
			if err == nil {
				lock.Lock()
				for _, v1HPA := range v1HPAs {
					output = append(output, UCEntity.SimpleHPAData{
						APIVersion:      constant.AutoscalingV1,
						Name:            v1HPA.Name,
						Namespace:       ns.Name,
						MinReplicas:     v1HPA.Spec.MinReplicas,
						MaxReplicas:     v1HPA.Spec.MaxReplicas,
						CurrentReplicas: v1HPA.Status.CurrentReplicas,
					})
				}
				lock.Unlock()
			}
		}(namespace)
		go func(ns v1.Namespace) {
			defer wg.Done()
			v2HPAs, err := c.hpaRepo.GetAllV2HPA(ctx, client, ns, clusterID)
			if err == nil {
				lock.Lock()
				for _, v2HPA := range v2HPAs {
					output = append(output, UCEntity.SimpleHPAData{
						APIVersion:      constant.AutoscalingV2,
						Name:            v2HPA.Name,
						Namespace:       ns.Name,
						MinReplicas:     v2HPA.Spec.MinReplicas,
						MaxReplicas:     v2HPA.Spec.MaxReplicas,
						CurrentReplicas: v2HPA.Status.CurrentReplicas,
					})
				}
				lock.Unlock()
			}
		}(namespace)
	}
	wg.Wait()
	return output, nil
}
