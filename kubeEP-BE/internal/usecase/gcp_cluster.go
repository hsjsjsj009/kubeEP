package useCase

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	gcpCustomAuth "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/k8s/auth/gcp_custom"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/k8s/client"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	containerProto "google.golang.org/genproto/googleapis/container/v1"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
)

type GCPCluster interface {
	RegisterGoogleCredentials(credentialsName string, gcpCredentials *google.Credentials)
	GetAllClustersInGCPProject(ctx context.Context, projectID string, clusterClient *container.ClusterManagerClient) ([]*UCEntity.GCPClusterData, error)
	GetGoogleClusterClient(ctx context.Context, googleCredential *google.Credentials) (*container.ClusterManagerClient, error)
	RegisterClusters(tx *gorm.DB, datacenterID uuid.UUID, listCluster []*UCEntity.GCPClusterData) error
	GetKubernetesClusterClient(credentialsName string, clusterData *UCEntity.ClusterData) (*kubernetes.Clientset, error)
}

type gcpCluster struct {
	validatorInst *validator.Validate
	clusterRepo   repository.Cluster
}

func NewGCPCluster(validatorInst *validator.Validate, clusterRepo repository.Cluster) GCPCluster {
	return &gcpCluster{validatorInst: validatorInst, clusterRepo: clusterRepo}
}

func (c *gcpCluster) RegisterGoogleCredentials(credentialsName string, gcpCredentials *google.Credentials) {
	gcpCustomAuth.RegisterGoogleCredentials(credentialsName, gcpCredentials)
}

func (c *gcpCluster) GetGoogleClusterClient(ctx context.Context, googleCredential *google.Credentials) (*container.ClusterManagerClient, error) {
	return container.NewClusterManagerClient(ctx, option.WithCredentials(googleCredential))
}

func (c *gcpCluster) GetAllClustersInGCPProject(ctx context.Context, projectID string, clusterClient *container.ClusterManagerClient) ([]*UCEntity.GCPClusterData, error) {
	clusters, err := clusterClient.ListClusters(
		ctx,
		&containerProto.ListClustersRequest{
			Parent: fmt.Sprintf("projects/%s/locations/-", projectID),
		},
	)
	if err != nil {
		return nil, err
	}
	var clusterData []*UCEntity.GCPClusterData
	for _, cluster := range clusters.GetClusters() {
		clusterData = append(clusterData, &UCEntity.GCPClusterData{
			ClusterData: UCEntity.ClusterData{
				Name:           fmt.Sprintf("gke_%s_%s_%s", projectID, cluster.GetName(), cluster.GetLocation()),
				Certificate:    cluster.GetMasterAuth().GetClusterCaCertificate(),
				ServerEndpoint: fmt.Sprintf("https://%s", cluster.GetEndpoint()),
				Datacenter: UCEntity.DatacenterDetailedData{
					Datacenter: constant.GCP,
				},
			},
			Location: cluster.GetLocation(),
		})
	}
	return clusterData, nil
}

func (c *gcpCluster) RegisterClusters(tx *gorm.DB, datacenterID uuid.UUID, listCluster []*UCEntity.GCPClusterData) error {
	var clusters []*model.Cluster
	for _, cluster := range listCluster {
		metadata := UCEntity.GCPClusterMetaData{Location: cluster.Location}
		metadataByte, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		clusterModel := &model.Cluster{
			Name:           cluster.Name,
			ServerEndpoint: cluster.ServerEndpoint,
			Certificate:    cluster.Certificate,
		}
		clusterModel.DatacenterID.SetUUID(datacenterID)
		clusterModel.Metadata.SetRawMessage(metadataByte)
		clusters = append(clusters, clusterModel)
	}

	err := c.clusterRepo.InsertClusterBatch(tx, clusters)
	if err != nil {
		return err
	}

	for idx, cluster := range clusters {
		listCluster[idx].ID = cluster.ID.GetUUID()
	}

	return nil
}

func (c gcpCluster) GetKubernetesClusterClient(credentialsName string, clusterData *UCEntity.ClusterData) (*kubernetes.Clientset, error) {
	if clusterData.Datacenter.Datacenter != constant.GCP {
		return nil, errors.New(errorConstant.DatacenterMismatch)
	}

	credentials := &k8sClient.Credentials{
		Certificate:    clusterData.Certificate,
		Name:           clusterData.Name,
		ServerEndpoint: clusterData.ServerEndpoint,
		AuthProviderConfig: &api.AuthProviderConfig{
			Name: "gcp_custom",
			Config: map[string]string{
				gcpCustomAuth.CredentialsNameConfigKey: credentialsName,
			},
		},
	}

	return k8sClient.GetClient(credentials)
}
