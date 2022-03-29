package cron

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"errors"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"k8s.io/client-go/kubernetes"
)

func (c *cron) getAllGCPClient(
	ctx context.Context,
	clusterData *UCEntity.ClusterData,
) (kubernetes.Interface, *container.ClusterManagerClient, error) {
	datacenter := clusterData.Datacenter.Datacenter
	if datacenter != model.GCP {
		return nil, nil, errors.New(errorConstant.DatacenterMismatch)
	}
	datacenterName := clusterData.Datacenter.Name
	datacenterData := UCEntity.DatacenterData{
		Credentials: clusterData.Datacenter.Credentials,
		Name:        datacenterName,
	}
	googleCredential, err := c.gcpDatacenterUC.GetGoogleCredentials(
		ctx,
		datacenterData,
	)
	gcpClusterClient, err := c.gcpClusterUC.GetGoogleClusterClient(ctx, googleCredential)
	if err != nil {
		return nil, nil, err
	}
	c.gcpClusterUC.RegisterGoogleCredentials(datacenterName, googleCredential)
	kubernetesClient, err := c.gcpClusterUC.GetKubernetesClusterClient(
		datacenterName,
		clusterData,
	)
	if err != nil {
		return nil, nil, err
	}
	return kubernetesClient, gcpClusterClient, nil
}
