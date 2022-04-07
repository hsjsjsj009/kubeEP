package response

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
)

type Cluster struct {
	ID         *uuid.UUID               `json:"id,omitempty"`
	Name       string                   `json:"name"`
	Datacenter model.DatacenterProvider `json:"datacenter"`
}

type UpdatedNodePool struct {
	ID           uuid.UUID `json:"id"`
	NodePoolName string    `json:"node_pool_name"`
}
