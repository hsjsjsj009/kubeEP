package repository

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	containerEntity "google.golang.org/genproto/googleapis/container/v1"
)

type GCPCluster interface {
	GetAllCluster(
		ctx context.Context,
		clusterClient *container.ClusterManagerClient,
		projectID string,
	) (
		*containerEntity.ListClustersResponse, error,
	)
}

type gcpCluster struct {
}

func newGcpCluster() GCPCluster {
	return &gcpCluster{}
}

func (g *gcpCluster) GetAllCluster(
	ctx context.Context,
	clusterClient *container.ClusterManagerClient,
	projectID string,
) (
	*containerEntity.ListClustersResponse, error,
) {
	return clusterClient.ListClusters(
		ctx,
		&containerEntity.ListClustersRequest{
			Parent: fmt.Sprintf("projects/%s/locations/-", projectID),
		},
	)
}
