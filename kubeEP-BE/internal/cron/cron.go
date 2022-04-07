package cron

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	errorConstant "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant/errors"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	useCase "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/usecase"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	containerGCP "google.golang.org/genproto/googleapis/container/v1"
	"gorm.io/gorm"
	v1Apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/autoscaling/v1"
	"k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/autoscaling/v2beta2"
	v1Option "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"math"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Cron interface {
	Start()
}

type cron struct {
	eventUC              useCase.Event
	clusterUC            useCase.Cluster
	gcpClusterUC         useCase.GCPCluster
	gcpDatacenterUC      useCase.GCPDatacenter
	scheduledHPAConfigUC useCase.ScheduledHPAConfig
	updatedNodePoolUC    useCase.UpdatedNodePool
	tx                   *gorm.DB
}

func newCron(
	eventUC useCase.Event,
	clusterUC useCase.Cluster,
	gcpClusterUC useCase.GCPCluster,
	gcpDatacenterUC useCase.GCPDatacenter,
	scheduledHPAConfigUC useCase.ScheduledHPAConfig,
	updatedNodePoolUC useCase.UpdatedNodePool,
	tx *gorm.DB,
) Cron {
	return &cron{
		eventUC:              eventUC,
		tx:                   tx,
		clusterUC:            clusterUC,
		gcpClusterUC:         gcpClusterUC,
		gcpDatacenterUC:      gcpDatacenterUC,
		scheduledHPAConfigUC: scheduledHPAConfigUC,
		updatedNodePoolUC:    updatedNodePoolUC,
	}
}

func (c *cron) handleExecEventError(db *gorm.DB, e *UCEntity.Event, errMsg string) {
	e.Status = model.EventFailed
	e.Message = errMsg
	err := c.eventUC.UpdateEvent(db, e)
	if err != nil {
		log.Errorf("[EventCronJob] Error Update Event : %s", err.Error())
	}
	log.Errorf("[EventCronJob] Event : %s, Error : %s", e.Name, errMsg)
}

func (c *cron) handleWatchEvent(db *gorm.DB, e *UCEntity.Event, errMsg string) {
	e.Message = errMsg
	err := c.eventUC.UpdateEvent(db, e)
	if err != nil {
		log.Errorf("[EventCronJob] Error Update Event : %s", err.Error())
	}
	log.Errorf("[EventCronJob] Watching event : %s, Error : %s", e.Name, errMsg)
}

func (c *cron) execEvent(e *UCEntity.Event, db *gorm.DB, ctx context.Context) {
	log.Infof("[EventCronJob] Executing event %s", e.Name)
	e.Status = model.EventExecuting

	err := c.eventUC.UpdateEvent(db, e)
	if err != nil {
		log.Errorf("[EventCronJob] Error Update Event : %s", err.Error())
		return
	}

	clusterID := e.Cluster.ID
	clusterData, err := c.clusterUC.GetClusterAndDatacenterDataByClusterID(db, clusterID)
	if err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}
	datacenter := clusterData.Datacenter.Datacenter
	var kubernetesClient kubernetes.Interface
	var googleContainerClient *container.ClusterManagerClient

	// Get Clients
	switch datacenter {
	case model.GCP:
		kubernetesClient, googleContainerClient, err = c.getAllGCPClient(ctx, clusterData)
		if err != nil {
			c.handleExecEventError(db, e, err.Error())
			return
		}
	}

	modifiedHPAs, err := c.scheduledHPAConfigUC.ListScheduledHPAConfigByEventID(db, e.ID)
	if err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}

	// Check HPA
	log.Infof("[EventCronJob] Event : %s, Checking HPAs", e.Name)
	existingK8sHPA, err := c.clusterUC.GetAllK8sHPAObjectInCluster(
		ctx,
		kubernetesClient,
		clusterID,
		clusterData.LatestHPAAPIVersion,
	)
	if err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}

	var selectedK8sHPAs []interface{}
	var existingModifiedHPAs []*UCEntity.EventModifiedHPAConfigData
	for _, modifiedHPA := range modifiedHPAs {
		hpaExist := false
		for _, hpa := range existingK8sHPA {
			switch h := hpa.(type) {
			case v1.HorizontalPodAutoscaler:
				if modifiedHPA.Name == h.Name && modifiedHPA.Namespace == h.Namespace {
					hpaExist = true
					selectedK8sHPAs = append(selectedK8sHPAs, h.DeepCopy())
					existingModifiedHPAs = append(existingModifiedHPAs, modifiedHPA)
					break
				}
			case v2beta1.HorizontalPodAutoscaler:
				if modifiedHPA.Name == h.Name && modifiedHPA.Namespace == h.Namespace {
					hpaExist = true
					selectedK8sHPAs = append(selectedK8sHPAs, h.DeepCopy())
					existingModifiedHPAs = append(existingModifiedHPAs, modifiedHPA)
					break
				}
			case v2beta2.HorizontalPodAutoscaler:
				if modifiedHPA.Name == h.Name && modifiedHPA.Namespace == h.Namespace {
					hpaExist = true
					selectedK8sHPAs = append(selectedK8sHPAs, h.DeepCopy())
					existingModifiedHPAs = append(existingModifiedHPAs, modifiedHPA)
					break
				}
			}
		}
		if !hpaExist {
			err := c.scheduledHPAConfigUC.UpdateScheduledHPAConfigStatusMessage(
				db,
				modifiedHPA.ID,
				model.HPAUpdateFailed,
				"hpa not found",
			)
			if err != nil {
				log.Errorf(
					"[EventCronJob] Error Update HPA %s Namespace %s : %s",
					modifiedHPA.Name,
					modifiedHPA.Namespace,
					err.Error(),
				)
			}
			continue
		}
	}

	if len(selectedK8sHPAs) == 0 {
		c.handleExecEventError(db, e, "no hpa exist")
		return
	}

	// Parse GCP Cluster Name
	clusterMetadata := strings.Split(clusterData.Name, "_")
	project := clusterMetadata[1]
	location := clusterMetadata[3]
	name := clusterMetadata[2]

	// Get GCP Node Pools
	googleClusterData, err := googleContainerClient.GetCluster(
		ctx, &containerGCP.GetClusterRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/cluster/%s", project, location, name),
		},
	)
	if err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}

	// Get Linux Daemonsets and Calculate Required Resources
	log.Infof("[EventCronJob] Event : %s, Calculate daemonsets resources", e.Name)
	daemonsets, err := kubernetesClient.AppsV1().DaemonSets("kube-system").List(
		ctx,
		v1Option.ListOptions{},
	)
	if err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}
	linuxDaemonSetsPodAmount := int64(0)
	daemonSetsRequestedMemory := float64(0)
	daemonSetsRequestedCPU := float64(0)
	for _, daemonSet := range daemonsets.Items {
		if daemonSet.Spec.Template.Spec.NodeSelector["kubernetes.io/os"] == "linux" {
			for _, containerSpec := range daemonSet.Spec.Template.Spec.Containers {
				daemonSetsRequestedMemory += containerSpec.Resources.Requests.Memory().AsApproximateFloat64()
				daemonSetsRequestedCPU += containerSpec.Resources.Requests.Cpu().AsApproximateFloat64()
			}
			linuxDaemonSetsPodAmount++
		}
	}

	// Get Maximum Resources each Node Pools
	log.Infof(
		"[EventCronJob] Event : %s, Calculate maximum available resources in node pools",
		e.Name,
	)
	nodePools := googleClusterData.NodePools
	nodePoolsMaxResources := map[string]*NodePoolResourceData{}
	nodePoolsRequestedResources := map[string]*NodePoolRequestedResourceData{}
	nodePoolsMap := map[string]*containerGCP.NodePool{}
	var nodePoolsList []string
	errGroup, ctxEg := errgroup.WithContext(ctx)
	dbEg := db.WithContext(ctxEg).Begin()
	for _, nodePool := range nodePools {
		nodePoolsList = append(nodePoolsList, nodePool.Name)
		nodePoolsRequestedResources[nodePool.Name] = &NodePoolRequestedResourceData{}
		resourceData := &NodePoolResourceData{}
		nodePoolsMaxResources[nodePool.Name] = resourceData
		nodePoolsMap[nodePool.Name] = nodePool
		loadFunc := func(nP *containerGCP.NodePool, rD *NodePoolResourceData) func() error {
			return func() error {
				updatedNodePool := &model.UpdatedNodePool{
					NodePoolName: nP.Name,
				}
				updatedNodePool.EventID.SetUUID(e.ID)
				err := dbEg.Create(updatedNodePool).Error
				if err != nil {
					if ctxEg.Err() != nil {
						return nil
					}
					return err
				}

				nodePoolMaxPods := nP.MaxPodsConstraint.MaxPodsPerNode
				nodePoolMaxNode := nP.Autoscaling.MaxNodeCount
				availablePods := nodePoolMaxPods - linuxDaemonSetsPodAmount
				rD.AvailablePods = availablePods
				rD.MaxAvailablePods = availablePods * int64(nodePoolMaxNode)
				nodes, err := kubernetesClient.CoreV1().Nodes().List(
					ctxEg, v1Option.ListOptions{
						LabelSelector: fmt.Sprintf(
							"%s=%s",
							constant.GCPNodePoolLabel,
							nP.Name,
						),
						Limit: 1,
					},
				)
				if err != nil {
					if ctxEg.Err() != nil {
						return nil
					}
					return err
				}
				node := nodes.Items[0]
				allocatableCPU := node.Status.Allocatable.Cpu().AsApproximateFloat64()
				allocatableMemory := node.Status.Allocatable.Memory().AsApproximateFloat64()
				availableCPU := allocatableCPU - daemonSetsRequestedCPU
				availableMemory := allocatableMemory - daemonSetsRequestedMemory
				rD.AvailableCPU = availableCPU
				rD.AvailableMemory = availableMemory
				rD.MaxAvailableCPU = availableCPU * float64(nodePoolMaxNode)
				rD.MaxAvailableMemory = availableMemory * float64(nodePoolMaxNode)

				return nil
			}
		}
		errGroup.Go(loadFunc(nodePool, resourceData))
	}

	if err := errGroup.Wait(); err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}
	dbEg.Commit()

	// Calculate Required Resource
	log.Infof("[EventCronJob] Event : %s, Calculate required resources", e.Name)
	var nodePoolRequestedResourceLock sync.Mutex
	errGroup, ctxEg = errgroup.WithContext(ctx)

	for idx, selectedHPA := range selectedK8sHPAs {
		errGroup.Go(
			func(i int, hpa interface{}) func() error {
				return func() error {
					requestedModification := existingModifiedHPAs[i]
					var scaleTargetRef interface{}
					namespace := requestedModification.Namespace
					maxReplicas := requestedModification.MaxReplicas

					// Modify HPA, Get Target Ref and Namespace
					switch h := hpa.(type) {
					case *v1.HorizontalPodAutoscaler:
						h.Spec.MinReplicas = requestedModification.MinReplicas
						h.Spec.MaxReplicas = maxReplicas
						scaleTargetRef = h.Spec.ScaleTargetRef
					case *v2beta1.HorizontalPodAutoscaler:
						h.Spec.MinReplicas = requestedModification.MinReplicas
						h.Spec.MaxReplicas = maxReplicas
						scaleTargetRef = h.Spec.ScaleTargetRef
					case *v2beta2.HorizontalPodAutoscaler:
						h.Spec.MinReplicas = requestedModification.MinReplicas
						h.Spec.MaxReplicas = maxReplicas
						scaleTargetRef = h.Spec.ScaleTargetRef
					default:
						return errors.New(errorConstant.HPAVersionUnknown)
					}

					// Resolve Target Ref to Get Pods
					resolveResult, err := c.clusterUC.ResolveScaleTargetRef(
						ctxEg,
						kubernetesClient,
						scaleTargetRef,
						namespace,
					)
					if err != nil {
						if ctxEg.Err() != nil {
							return nil
						}
						return err
					}

					// Get Labels Selector and Calculate Requested Resource
					var nodeSelectorLabels []string
					var maxRequestedCPU, maxRequestedMemory float64
					switch resolveRes := resolveResult.(type) {
					case *v1Apps.Deployment:
						nodeSelector := resolveRes.Spec.Template.Spec.NodeSelector
						for label, value := range nodeSelector {
							nodeSelectorLabels = append(
								nodeSelectorLabels,
								fmt.Sprintf("%s=%s", label, value),
							)
						}

						totalCpuRequested := float64(0)
						totalMemoryRequested := float64(0)
						containers := resolveRes.Spec.Template.Spec.Containers
						for _, containerSpec := range containers {
							totalCpuRequested += containerSpec.Resources.Requests.Cpu().AsApproximateFloat64()
							totalMemoryRequested += containerSpec.Resources.Requests.Memory().AsApproximateFloat64()
						}
						maxRequestedCPU = totalCpuRequested * float64(maxReplicas)
						maxRequestedMemory = totalMemoryRequested * float64(maxReplicas)
					}

					//Resolve node selector and Find all node pools
					var hpaNodePools map[string]bool

					if len(nodeSelectorLabels) > 0 {
						selectedNode, err := kubernetesClient.CoreV1().Nodes().List(
							ctxEg, v1Option.ListOptions{
								LabelSelector: strings.Join(nodeSelectorLabels, ","),
							},
						)
						if err != nil {
							if ctxEg.Err() != nil {
								return nil
							}
							return err
						}
						for _, node := range selectedNode.Items {
							nodePoolName := node.Labels[constant.GCPNodePoolLabel]
							hpaNodePools[nodePoolName] = true
						}
					}

					// Calculate & Save requested resource data
					nodePoolRequestedResourceLock.Lock()
					defer nodePoolRequestedResourceLock.Unlock()

					if len(hpaNodePools) == 0 {
						for _, nodePoolName := range nodePoolsList {
							requestedResourceData := nodePoolsRequestedResources[nodePoolName]
							requestedResourceData.MaxPods += int64(maxReplicas)
							requestedResourceData.MaxCPU += maxRequestedCPU
							requestedResourceData.MaxMemory += maxRequestedMemory
						}
						return nil
					}

					for nodePoolName := range hpaNodePools {
						requestedResourceData := nodePoolsRequestedResources[nodePoolName]
						requestedResourceData.MaxPods += int64(maxReplicas)
						requestedResourceData.MaxCPU += maxRequestedCPU
						requestedResourceData.MaxMemory += maxRequestedMemory
					}

					return nil
				}
			}(idx, selectedHPA),
		)
	}

	if err := errGroup.Wait(); err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}

	// Calculate Requested Resource Each Node Pool and Update the Node Pool
	log.Infof(
		"[EventCronJob] Event : %s, Calculate needed pool based on requested resources",
		e.Name,
	)
	errGroup, ctxEg = errgroup.WithContext(ctx)
	var updateNodePoolLock sync.Mutex
	for _, nodePoolName := range nodePoolsList {
		requestedResourceData := nodePoolsRequestedResources[nodePoolName]
		maxResourceData := nodePoolsMaxResources[nodePoolName]
		nodePool := nodePoolsMap[nodePoolName]
		errGroup.Go(
			func(
				reqResources *NodePoolRequestedResourceData,
				maxResources *NodePoolResourceData,
				nodePoolObj *containerGCP.NodePool,
			) func() error {
				return func() error {
					unfulfilledCPU := float64(0)
					if reqResources.MaxCPU > maxResources.MaxAvailableCPU {
						unfulfilledCPU = reqResources.MaxCPU - maxResources.MaxAvailableCPU
					}

					unfulfilledMemory := float64(0)
					if reqResources.MaxMemory > maxResources.MaxAvailableMemory {
						unfulfilledMemory = reqResources.MaxMemory - maxResources.MaxAvailableMemory
					}

					unfulfilledPods := int64(0)
					if reqResources.MaxPods > maxResources.MaxAvailablePods {
						unfulfilledPods = reqResources.MaxPods - maxResources.MaxAvailablePods
					}

					neededNodeBasedOnCPU := math.Ceil(unfulfilledCPU / maxResources.AvailableCPU)
					neededNodeBasedOnMemory := math.Ceil(unfulfilledMemory / maxResources.AvailableMemory)
					neededNodeBasedOnPods := math.Ceil(float64(unfulfilledPods) / float64(maxResources.AvailablePods))

					autoscalingData := nodePoolObj.Autoscaling

					maxNeededNode := int32(
						math.Max(
							neededNodeBasedOnCPU,
							math.Max(neededNodeBasedOnMemory, neededNodeBasedOnPods),
						),
					)

					newMaxNode := autoscalingData.MaxNodeCount

					if maxNeededNode > 0 {
						newMaxNode += maxNeededNode + 5
					}

					updateNodePoolLock.Lock()
					defer updateNodePoolLock.Unlock()
					log.Infof(
						"[EventCronJob] Event : %s, Updating GCP node pool %s with new max node size %d (before : %d)",
						e.Name,
						nodePoolObj.Name,
						newMaxNode,
						autoscalingData.MaxNodeCount,
					)

					autoscalingData.MaxNodeCount = newMaxNode

					op, err := googleContainerClient.SetNodePoolAutoscaling(
						ctxEg, &containerGCP.SetNodePoolAutoscalingRequest{
							Name: fmt.Sprintf(
								"projects/%s/locations/%s/clusters/%s/nodePools/%s",
								project,
								location,
								name,
								nodePoolObj.Name,
							),
							Autoscaling: autoscalingData,
						},
					)
					if err != nil {
						if ctxEg.Err() != nil {
							return nil
						}
						return err
					}
					for {
						opIns, err := googleContainerClient.GetOperation(
							ctx,
							&containerGCP.GetOperationRequest{Name: fmt.Sprintf(
								"projects/%s/locations/%s/operations/%s",
								project,
								location,
								op.Name,
							)},
						)
						if err != nil {
							return err
						}
						if opIns.Status == containerGCP.Operation_DONE {
							if opIns.Error != nil {
								return errors.New(opIns.Error.String())
							}
							return nil
						}
						time.Sleep(100 * time.Millisecond)
					}
				}
			}(requestedResourceData, maxResourceData, nodePool),
		)
	}

	if err := errGroup.Wait(); err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}

	// Update K8s HPA
	log.Infof("[EventCronJob] Event : %s, Updating K8s HPA with new configuration", e.Name)
	err = c.clusterUC.UpdateHPAK8sObjectBatch(ctx, kubernetesClient, clusterID, selectedK8sHPAs)
	if err != nil {
		c.handleExecEventError(db, e, err.Error())
		return
	}

	for _, existingModifiedHPA := range existingModifiedHPAs {
		err := c.scheduledHPAConfigUC.UpdateScheduledHPAConfigStatusMessage(
			db,
			existingModifiedHPA.ID,
			model.HPAUpdateSuccess,
			"",
		)
		if err != nil {
			c.handleExecEventError(
				db, e, fmt.Sprintf(
					"Error Update HPA %s Namespace %s : %s", existingModifiedHPA.Name,
					existingModifiedHPA.Namespace,
					err.Error(),
				),
			)
			return
		}
	}

	e.Status = model.EventPrescaled

	err = c.eventUC.UpdateEvent(db, e)
	if err != nil {
		log.Errorf("[EventCronJob] Error Update Event : %s", err.Error())
	}

	log.Infof("[EventCronJob] Event : %s, Done executing update and calculation", e.Name)
}

func (c *cron) watchNodePool(
	client kubernetes.Interface,
	db *gorm.DB,
	provider model.DatacenterProvider,
	event *UCEntity.Event,
	now time.Time,
	ctx context.Context,
	updatedNodePoolMap map[string]uuid.UUID,
) {
	nodes, err := client.CoreV1().Nodes().List(ctx, v1Option.ListOptions{})
	if err != nil {
		log.Errorf(
			"[EventCronJob] Watching event : %s, Watch node pool error : %s",
			event.Name,
			err.Error(),
		)
		return
	}
	nodeCounts := map[string]int32{}
	for _, node := range nodes.Items {
		nodeLabels := node.Labels
		var nodePoolName string
		switch provider {
		case model.GCP:
			nodePoolName = nodeLabels[constant.GCPNodePoolLabel]
			nodeCounts[nodePoolName] += 1
		}
	}

	var nodePoolStatusObjects []model.NodePoolStatus
	for nodePoolName, nodeCount := range nodeCounts {
		nodePoolStatus := model.NodePoolStatus{
			CreatedAt: now,
			NodeCount: nodeCount,
		}
		nodePoolStatus.UpdatedNodePoolID.SetUUID(updatedNodePoolMap[nodePoolName])
		nodePoolStatusObjects = append(nodePoolStatusObjects, nodePoolStatus)
	}

	err = db.Create(&nodePoolStatusObjects).Error
	if err != nil {
		log.Errorf(
			"[EventCronJob] Watching event : %s, Watch node pool error : %s",
			event.Name,
			err.Error(),
		)
	}
	log.Infof(
		"[EventCronJob] Watching event : %s, Watching node pool at : %s",
		event.Name,
		now,
	)
}

func (c *cron) watchHPA(
	hpaListFunc func(ctx context.Context) ([]UCEntity.SimpleHPAData, error),
	db *gorm.DB,
	event *UCEntity.Event,
	scheduledHPAConfigs []*UCEntity.EventModifiedHPAConfigData,
	now time.Time,
	ctx context.Context,
) {
	hpaList, err := hpaListFunc(ctx)
	if err != nil {
		log.Errorf(
			"[EventCronJob] Watching event : %s, Watch hpa error : %s",
			event.Name,
			err.Error(),
		)
		return
	}

	var selectedHPAStatuses []model.HPAStatus
	for _, hpa := range hpaList {
		for _, scheduledHPAConfig := range scheduledHPAConfigs {
			if hpa.Name == scheduledHPAConfig.Name && scheduledHPAConfig.Namespace == hpa.Namespace {
				hpaStatus := model.HPAStatus{
					CreatedAt: now,
					PodCount:  hpa.CurrentReplicas,
				}
				hpaStatus.ScheduledHPAConfigID.SetUUID(scheduledHPAConfig.ID)
				selectedHPAStatuses = append(selectedHPAStatuses, hpaStatus)
			}
		}
	}

	err = db.Create(&selectedHPAStatuses).Error
	if err != nil {
		log.Errorf(
			"[EventCronJob] Watching event : %s, Watch hpa error : %s",
			event.Name,
			err.Error(),
		)
		return
	}

	log.Infof(
		"[EventCronJob] Watching event : %s, Watching hpa at : %s",
		event.Name,
		now,
	)
}

func (c *cron) watchEvent(e *UCEntity.Event, db *gorm.DB, ctx context.Context) {
	log.Infof("[EventCronJob] Watching event %s", e.Name)
	e.Status = model.EventWatching

	err := c.eventUC.UpdateEvent(db, e)
	if err != nil {
		log.Errorf("[EventCronJob] Error update event : %s", err.Error())
		return
	}

	clusterID := e.Cluster.ID
	clusterData, err := c.clusterUC.GetClusterAndDatacenterDataByClusterID(db, clusterID)
	if err != nil {
		c.handleWatchEvent(db, e, err.Error())
		return
	}
	datacenter := clusterData.Datacenter.Datacenter
	var kubernetesClient kubernetes.Interface

	// Get Clients
	switch datacenter {
	case model.GCP:
		kubernetesClient, _, err = c.getAllGCPClient(ctx, clusterData)
		if err != nil {
			c.handleWatchEvent(db, e, err.Error())
			return
		}
	}

	scheduledHPAConfigs, err := c.scheduledHPAConfigUC.ListScheduledHPAConfigByEventID(db, e.ID)
	if err != nil {
		c.handleWatchEvent(db, e, err.Error())
		return
	}

	getAllHPAFunc := func(ctx context.Context) ([]UCEntity.SimpleHPAData, error) {
		return c.clusterUC.GetAllHPAInCluster(
			ctx,
			kubernetesClient,
			clusterID,
			clusterData.LatestHPAAPIVersion,
		)
	}

	updatedNodePools, err := c.updatedNodePoolUC.GetAllUpdatedNodePoolByEvent(db, e.ID)
	if err != nil {
		c.handleWatchEvent(db, e, err.Error())
		return
	}
	updatedNodePoolMap := map[string]uuid.UUID{}
	for _, updatedNodePool := range updatedNodePools {
		updatedNodePoolMap[updatedNodePool.NodePoolName] = updatedNodePool.UpdatedNodePoolID
	}

	endTime := e.EndTime
	watcherTicker := time.NewTicker(30 * time.Second)
	defer watcherTicker.Stop()
	for {
		select {
		case now := <-watcherTicker.C:
			if now.After(endTime) {
				return
			}

			go c.watchNodePool(kubernetesClient, db, datacenter, e, now, ctx, updatedNodePoolMap)
			go c.watchHPA(getAllHPAFunc, db, e, scheduledHPAConfigs, now, ctx)
		case <-ctx.Done():
			return
		}
	}

}

func (c *cron) Start() {
	log.Infof("Starting event cron job")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	db := c.tx.WithContext(ctx)
	mainTicker := time.NewTicker(1 * time.Minute)
	defer mainTicker.Stop()
	for {
		select {
		case now := <-mainTicker.C:
			go func() {
				pendingEvents, err := c.eventUC.GetAllPendingExecutableEvent(db, now)
				if err != nil {
					log.Errorf(
						"[EventCronJob] Error getting pending executable events : %s",
						err.Error(),
					)
				}
				if len(pendingEvents) != 0 && err == nil {
					for _, pendingEvent := range pendingEvents {
						go c.execEvent(pendingEvent, db, ctx)
					}
				}
			}()

			go func() {
				prescaledEvents, err := c.eventUC.GetAllPrescaledEvent(db, now.Add(-5*time.Minute))
				if err != nil {
					log.Errorf("[EventCronJob] Error getting prescaled events : %s", err.Error())
				}
				if len(prescaledEvents) != 0 && err == nil {
					for _, prescaledEvent := range prescaledEvents {
						go c.watchEvent(prescaledEvent, db, ctx)
					}
				}
			}()

			go func() {
				err := c.eventUC.FinishAllWatchedEvent(db, now)
				if err != nil {
					log.Errorf("[EventCronJob] Error update watched events : %s", err.Error())
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}
