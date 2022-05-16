package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	"k8s.io/api/core/v1"
	v1Option "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sNode interface {
	GetNodesFromGCPNodePool(
		ctx context.Context,
		k8sClient kubernetes.Interface,
		nodePoolName string,
	) (*v1.NodeList, error)
}

type k8sNode struct {
}

func newK8sNode() K8sNode {
	return &k8sNode{}
}

func (n *k8sNode) GetNodesFromGCPNodePool(
	ctx context.Context,
	k8sClient kubernetes.Interface,
	nodePoolName string,
) (*v1.NodeList, error) {
	data, err := k8sClient.CoreV1().Nodes().List(
		ctx, v1Option.ListOptions{
			LabelSelector: fmt.Sprintf(
				"%s=%s",
				constant.GCPNodePoolLabel,
				nodePoolName,
			),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(data.Items) == 0 {
		return nil, errors.New(errorConstant.NoExistingNode)
	}
	return data, nil
}
