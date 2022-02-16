package response

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
)

type Cluster struct {
	ID         *uuid.UUID              `json:"id,omitempty"`
	Name       string                  `json:"name"`
	Datacenter constant.DatacenterType `json:"datacenter"`
}
