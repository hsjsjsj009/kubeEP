package gcpUseCase

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	gcpUCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase/gcp"
	gcpCustomAuth "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/k8s/auth/gcp_custom"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	containerProto "google.golang.org/genproto/googleapis/container/v1"
	"gorm.io/gorm"
)

type Cluster interface {
	RegisterGoogleCredentials(credentialsName string, gcpCredentials *google.Credentials)
	GetAllClustersInGCPProject(ctx context.Context, projectID string, clusterClient *container.ClusterManagerClient) ([]*gcpUCEntity.ClusterData, error)
	GetGoogleClusterClient(ctx context.Context, googleCredential *google.Credentials) (*container.ClusterManagerClient, error)
	RegisterClusters(tx *gorm.DB, datacenterID uuid.UUID, listCluster []*gcpUCEntity.ClusterData) error
}

type cluster struct {
	validatorInst *validator.Validate
	clusterRepo   repository.Cluster
}

func NewCluster(validatorInst *validator.Validate, clusterRepo repository.Cluster) Cluster {
	return &cluster{validatorInst: validatorInst, clusterRepo: clusterRepo}
}

func (c *cluster) RegisterGoogleCredentials(credentialsName string, gcpCredentials *google.Credentials) {
	gcpCustomAuth.RegisterGoogleCredentials(credentialsName, gcpCredentials)
}

func (c *cluster) GetGoogleClusterClient(ctx context.Context, googleCredential *google.Credentials) (*container.ClusterManagerClient, error) {
	return container.NewClusterManagerClient(ctx, option.WithCredentials(googleCredential))
}

func (c *cluster) GetAllClustersInGCPProject(ctx context.Context, projectID string, clusterClient *container.ClusterManagerClient) ([]*gcpUCEntity.ClusterData, error) {
	clusters, err := clusterClient.ListClusters(
		ctx,
		&containerProto.ListClustersRequest{
			Parent: fmt.Sprintf("projects/%s/locations/-", projectID),
		},
	)
	if err != nil {
		return nil, err
	}
	var clusterData []*gcpUCEntity.ClusterData
	for _, cluster := range clusters.GetClusters() {
		clusterData = append(clusterData, &gcpUCEntity.ClusterData{
			ClusterData: UCEntity.ClusterData{
				Name:           fmt.Sprintf("gke_%s_%s_%s", projectID, cluster.GetName(), cluster.GetLocation()),
				Certificate:    cluster.GetMasterAuth().GetClusterCaCertificate(),
				ServerEndpoint: fmt.Sprintf("https://%s", cluster.GetEndpoint()),
				Datacenter:     constant.GCP,
			},
			Location: cluster.GetLocation(),
		})
	}
	return clusterData, nil
}

func (c *cluster) RegisterClusters(tx *gorm.DB, datacenterID uuid.UUID, listCluster []*gcpUCEntity.ClusterData) error {
	var clusters []*model.Cluster
	for _, cluster := range listCluster {
		metadata := gcpUCEntity.ClusterMetaData{Location: cluster.Location}
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
