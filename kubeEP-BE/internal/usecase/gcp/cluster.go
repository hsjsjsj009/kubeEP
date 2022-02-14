package gcpUseCase

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	gcpUCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase/gcp"
	gcpCustomAuth "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/k8s/auth/gcp_custom"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	containerProto "google.golang.org/genproto/googleapis/container/v1"
)

type Cluster interface {
	RegisterGoogleCredentials(credentialsName string, gcpCredentials *google.Credentials)
	GetAllClustersInProject(ctx context.Context, projectID string, clusterClient *container.ClusterManagerClient) ([]*gcpUCEntity.ClusterData, error)
	GetGoogleClusterClient(ctx context.Context, googleCredential *google.Credentials) (*container.ClusterManagerClient, error)
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

func (c cluster) GetGoogleClusterClient(ctx context.Context, googleCredential *google.Credentials) (*container.ClusterManagerClient, error) {
	return container.NewClusterManagerClient(ctx, option.WithCredentials(googleCredential))
}

func (c *cluster) GetAllClustersInProject(ctx context.Context, projectID string, clusterClient *container.ClusterManagerClient) ([]*gcpUCEntity.ClusterData, error) {
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
			Name:           cluster.GetName(),
			Certificate:    cluster.GetMasterAuth().GetClusterCaCertificate(),
			ServerEndpoint: fmt.Sprintf("https://%s", cluster.GetEndpoint()),
			Location:       cluster.GetLocation(),
		})
	}
	return clusterData, nil
}
