package useCase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"gorm.io/gorm"
	v1hpa "k8s.io/api/autoscaling/v1"
	"k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/autoscaling/v2beta2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

type Cluster interface {
	GetAllClustersInLocalByDatacenterID(
		tx *gorm.DB,
		datacenterID uuid.UUID,
	) ([]UCEntity.ClusterData, error)
	GetAllClustersInLocal(tx *gorm.DB) ([]UCEntity.ClusterData, error)
	GetAllHPAInCluster(
		ctx context.Context,
		client kubernetes.Interface,
		clusterID uuid.UUID,
	) (output []UCEntity.SimpleHPAData, err error)
	GetClusterAndDatacenterDataByClusterID(tx *gorm.DB, id uuid.UUID) (*UCEntity.ClusterData, error)
	GetLatestHPAAPIVersion(
		k8sClient kubernetes.Interface,
	) (
		constant.HPAVersion,
		error,
	)
}

type cluster struct {
	validatorInst *validator.Validate
	clusterRepo   repository.Cluster
	hpaRepo       repository.K8sHPA
	namespaceRepo repository.K8sNamespace
	discoveryRepo repository.K8SDiscovery
}

func newCluster(
	validatorInst *validator.Validate,
	clusterRepo repository.Cluster,
	hpaRepo repository.K8sHPA,
	namespaceRepo repository.K8sNamespace,
	discoveryRepo repository.K8SDiscovery,
) Cluster {
	return &cluster{
		validatorInst: validatorInst,
		clusterRepo:   clusterRepo,
		hpaRepo:       hpaRepo,
		namespaceRepo: namespaceRepo,
		discoveryRepo: discoveryRepo,
	}
}

func (c *cluster) GetLatestHPAAPIVersion(
	k8sClient kubernetes.Interface,
) (
	constant.HPAVersion,
	error,
) {
	response, err := c.discoveryRepo.GetServerGroups(k8sClient)
	if err != nil {
		return "", err
	}
	var autoscalingAPIGroup v1.APIGroup
	for _, apiGroup := range response.Groups {
		if apiGroup.Name == "autoscaling" {
			autoscalingAPIGroup = apiGroup
			break
		}
	}
	versionCount := len(autoscalingAPIGroup.Versions)
	latestVersion := autoscalingAPIGroup.Versions[versionCount-1]
	return latestVersion.GroupVersion, nil
}

func (c *cluster) GetAllClustersInLocalByDatacenterID(
	tx *gorm.DB,
	datacenterID uuid.UUID,
) ([]UCEntity.ClusterData, error) {
	clusters, err := c.clusterRepo.ListClusterByDatacenterID(tx, datacenterID)
	if err != nil {
		return nil, err
	}
	var output []UCEntity.ClusterData
	for _, cluster := range clusters {
		output = append(
			output, UCEntity.ClusterData{
				ID:             cluster.ID.GetUUID(),
				Name:           cluster.Name,
				Certificate:    cluster.Certificate,
				ServerEndpoint: cluster.ServerEndpoint,
				Datacenter: UCEntity.DatacenterDetailedData{
					Datacenter: cluster.Datacenter.Datacenter,
				},
			},
		)
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
		output = append(
			output, UCEntity.ClusterData{
				ID:             cluster.ID.GetUUID(),
				Name:           cluster.Name,
				Certificate:    cluster.Certificate,
				ServerEndpoint: cluster.ServerEndpoint,
				Datacenter: UCEntity.DatacenterDetailedData{
					Datacenter: cluster.Datacenter.Datacenter,
				},
			},
		)
	}
	return output, nil
}

func (c *cluster) GetClusterAndDatacenterDataByClusterID(
	tx *gorm.DB,
	id uuid.UUID,
) (*UCEntity.ClusterData, error) {
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

func (c *cluster) GetAllHPAInCluster(
	ctx context.Context,
	client kubernetes.Interface,
	clusterID uuid.UUID,
) (output []UCEntity.SimpleHPAData, err error) {
	namespaces, err := c.namespaceRepo.GetAllNamespace(ctx, client)
	if err != nil {
		return nil, err
	}
	var v1HPAs []v1hpa.HorizontalPodAutoscaler
	var v2b2HPAs []v2beta2.HorizontalPodAutoscaler
	var v2b1HPAs []v2beta1.HorizontalPodAutoscaler
	var errV1, errV2b2, errV2b1 error
	for _, namespace := range namespaces {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			v2b2HPAs, errV2b2 = c.hpaRepo.GetAllV2beta2HPA(ctx, client, namespace, clusterID)
		}()
		go func() {
			defer wg.Done()
			v2b1HPAs, errV2b1 = c.hpaRepo.GetAllV2beta1HPA(ctx, client, namespace, clusterID)
		}()
		wg.Wait()
		if errV1 == nil {
			for _, v1HPA := range v1HPAs {
				output = append(
					output, UCEntity.SimpleHPAData{
						APIVersion:      constant.AutoscalingV1,
						Name:            v1HPA.Name,
						Namespace:       namespace.Name,
						MinReplicas:     v1HPA.Spec.MinReplicas,
						MaxReplicas:     v1HPA.Spec.MaxReplicas,
						CurrentReplicas: v1HPA.Status.CurrentReplicas,
						ScaleTargetRef: UCEntity.HPAScaleTargetRef{
							Name: v1HPA.Spec.ScaleTargetRef.Name,
							Kind: v1HPA.Spec.ScaleTargetRef.Kind,
						},
					},
				)
			}
		}
	}
	return output, nil
}
