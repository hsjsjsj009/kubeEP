package repository

import (
	"context"
	v1Autoscale "k8s.io/api/autoscaling/v1"
	v2Autoscale "k8s.io/api/autoscaling/v2"
	v1Core "k8s.io/api/core/v1"
	v1Option "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sHPA interface {
	GetAllV1HPA(ctx context.Context, client kubernetes.Interface, namespace v1Core.Namespace) ([]v1Autoscale.HorizontalPodAutoscaler, error)
	GetAllV2HPA(ctx context.Context, client kubernetes.Interface, namespace v1Core.Namespace) ([]v2Autoscale.HorizontalPodAutoscaler, error)
}

type k8sHPA struct {
}

func newK8sHPA() K8sHPA {
	return &k8sHPA{}
}

func (h *k8sHPA) GetAllV1HPA(ctx context.Context, client kubernetes.Interface, namespace v1Core.Namespace) ([]v1Autoscale.HorizontalPodAutoscaler, error) {
	data, err := client.
		AutoscalingV1().
		HorizontalPodAutoscalers(namespace.Name).
		List(
			ctx,
			v1Option.ListOptions{},
		)
	if err != nil {
		return nil, err
	}
	return data.Items, nil
}

func (h *k8sHPA) GetAllV2HPA(ctx context.Context, client kubernetes.Interface, namespace v1Core.Namespace) ([]v2Autoscale.HorizontalPodAutoscaler, error) {
	data, err := client.
		AutoscalingV2().
		HorizontalPodAutoscalers(namespace.Name).
		List(
			ctx,
			v1Option.ListOptions{},
		)
	if err != nil {
		return nil, err
	}
	return data.Items, nil
}
