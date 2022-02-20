package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	v1Autoscale "k8s.io/api/autoscaling/v1"
	v2Autoscale "k8s.io/api/autoscaling/v2"
	v1Core "k8s.io/api/core/v1"
	v1Option "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type K8sHPA interface {
	GetAllV1HPA(ctx context.Context, client kubernetes.Interface, namespace v1Core.Namespace, clusterID uuid.UUID) ([]v1Autoscale.HorizontalPodAutoscaler, error)
	GetAllV2HPA(ctx context.Context, client kubernetes.Interface, namespace v1Core.Namespace, clusterID uuid.UUID) ([]v2Autoscale.HorizontalPodAutoscaler, error)
}

type k8sHPA struct {
	redisClient *redis.Client
}

func newK8sHPA(redisClient *redis.Client) K8sHPA {
	return &k8sHPA{
		redisClient: redisClient,
	}
}

func (h *k8sHPA) GetAllV1HPA(ctx context.Context, client kubernetes.Interface, namespace v1Core.Namespace, clusterID uuid.UUID) ([]v1Autoscale.HorizontalPodAutoscaler, error) {
	key := fmt.Sprintf("hpa_v1_list_cluster_%s", clusterID)
	if redisResponse := h.redisClient.Get(ctx, key); redisResponse.Err() != nil {
		var HPAList v1Autoscale.HorizontalPodAutoscalerList
		b, err := redisResponse.Bytes()
		if err == nil {
			if string(b) == errorConstant.HPAListError {
				return nil, errors.New(errorConstant.HPAListError)
			}
			if err = HPAList.Unmarshal(b); err == nil {
				return HPAList.Items, nil
			}
		}
	}
	data, err := client.
		AutoscalingV1().
		HorizontalPodAutoscalers(namespace.Name).
		List(
			ctx,
			v1Option.ListOptions{},
		)
	if err != nil {
		_ = h.redisClient.Set(ctx, key, errorConstant.HPAListError, 10*time.Second).Err()
		return nil, err
	}
	if b, err := data.Marshal(); err == nil {
		_ = h.redisClient.Set(ctx, key, b, 10*time.Second).Err()
	}
	return data.Items, nil
}

func (h *k8sHPA) GetAllV2HPA(ctx context.Context, client kubernetes.Interface, namespace v1Core.Namespace, clusterID uuid.UUID) ([]v2Autoscale.HorizontalPodAutoscaler, error) {
	key := fmt.Sprintf("hpa_v2_list_cluster_%s", clusterID)
	if redisResponse := h.redisClient.Get(ctx, key); redisResponse.Err() != nil {
		var HPAList v2Autoscale.HorizontalPodAutoscalerList
		b, err := redisResponse.Bytes()
		if err == nil {
			if string(b) == errorConstant.HPAListError {
				return nil, errors.New(errorConstant.HPAListError)
			}
			if err = HPAList.Unmarshal(b); err == nil {
				return HPAList.Items, nil
			}
		}
	}
	data, err := client.
		AutoscalingV2().
		HorizontalPodAutoscalers(namespace.Name).
		List(
			ctx,
			v1Option.ListOptions{},
		)
	if err != nil {
		_ = h.redisClient.Set(ctx, key, errorConstant.HPAListError, 10*time.Second).Err()
		return nil, err
	}
	if b, err := data.Marshal(); err == nil {
		_ = h.redisClient.Set(ctx, key, b, 10*time.Second).Err()
	}
	return data.Items, nil
}
