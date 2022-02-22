package useCase

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
	v1hpa "k8s.io/api/autoscaling/v1"
	"k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/autoscaling/v2beta2"
	v1Core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sort"
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
		latestHPAVersion constant.HPAVersion,
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
		LatestHPAAPIVersion: data.LatestHPAAPIVersion,
	}
	return clusterData, nil
}

func (c *cluster) GetAllHPAInCluster(
	ctx context.Context,
	client kubernetes.Interface,
	clusterID uuid.UUID,
	latestHPAVersion constant.HPAVersion,
) (output []UCEntity.SimpleHPAData, err error) {
	namespaces, err := c.namespaceRepo.GetAllNamespace(ctx, client)
	if err != nil {
		return nil, err
	}
	HPAs := map[string]interface{}{}
	var lock sync.Mutex
	eg, _ := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(4)
	//Get data
	for _, namespace := range namespaces {
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
		loadFunc := func(ns v1Core.Namespace) func() error {
			return func() error {
				defer sem.Release(1)

				var response interface{}
				var err error

				switch latestHPAVersion {
				case constant.AutoscalingV1:
					response, err = c.hpaRepo.GetAllV1HPA(ctx, client, ns, clusterID)
				case constant.AutoscalingV2Beta1:
					response, err = c.hpaRepo.GetAllV2beta1HPA(ctx, client, ns, clusterID)
				case constant.AutoscalingV2Beta2:
					response, err = c.hpaRepo.GetAllV2beta2HPA(ctx, client, ns, clusterID)
				default:
					return errors.New(errorConstant.HPAVersionUnknown)
				}
				if err != nil {
					return err
				}
				lock.Lock()
				HPAs[ns.Name] = response
				lock.Unlock()

				return nil
			}
		}
		eg.Go(loadFunc(namespace))
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var keys []string

	for k := range HPAs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, ns := range keys {
		v := HPAs[ns]
		switch chosenHPAs := v.(type) {
		case []v1hpa.HorizontalPodAutoscaler:
			for _, hpa := range chosenHPAs {
				output = append(
					output, UCEntity.SimpleHPAData{
						Name:            hpa.Name,
						Namespace:       ns,
						MinReplicas:     hpa.Spec.MinReplicas,
						MaxReplicas:     hpa.Spec.MaxReplicas,
						CurrentReplicas: hpa.Status.CurrentReplicas,
						ScaleTargetRef: UCEntity.HPAScaleTargetRef{
							Name: hpa.Spec.ScaleTargetRef.Name,
							Kind: hpa.Spec.ScaleTargetRef.Kind,
						},
					},
				)
			}
		case []v2beta1.HorizontalPodAutoscaler:
			for _, hpa := range chosenHPAs {
				output = append(
					output, UCEntity.SimpleHPAData{
						Name:            hpa.Name,
						Namespace:       ns,
						MinReplicas:     hpa.Spec.MinReplicas,
						MaxReplicas:     hpa.Spec.MaxReplicas,
						CurrentReplicas: hpa.Status.CurrentReplicas,
						ScaleTargetRef: UCEntity.HPAScaleTargetRef{
							Name: hpa.Spec.ScaleTargetRef.Name,
							Kind: hpa.Spec.ScaleTargetRef.Kind,
						},
					},
				)
			}
		case []v2beta2.HorizontalPodAutoscaler:
			for _, hpa := range chosenHPAs {
				output = append(
					output, UCEntity.SimpleHPAData{
						Name:            hpa.Name,
						Namespace:       ns,
						MinReplicas:     hpa.Spec.MinReplicas,
						MaxReplicas:     hpa.Spec.MaxReplicas,
						CurrentReplicas: hpa.Status.CurrentReplicas,
						ScaleTargetRef: UCEntity.HPAScaleTargetRef{
							Name: hpa.Spec.ScaleTargetRef.Name,
							Kind: hpa.Spec.ScaleTargetRef.Kind,
						},
					},
				)
			}
		}
	}
	return output, nil
}
