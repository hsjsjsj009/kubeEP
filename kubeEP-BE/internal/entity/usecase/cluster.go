package UCEntity

import (
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
)

type ClusterData struct {
	ID             uuid.UUID
	Name           string
	Certificate    string
	ServerEndpoint string
	Datacenter     constant.DatacenterType
}
