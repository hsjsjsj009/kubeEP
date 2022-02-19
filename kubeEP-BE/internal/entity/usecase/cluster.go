package UCEntity

import (
	"github.com/google/uuid"
)

type ClusterData struct {
	ID             uuid.UUID
	Name           string
	Certificate    string
	ServerEndpoint string
	Datacenter     DatacenterDetailedData
}
