package useCaseEntity

import (
	"encoding/json"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
)

type DatacenterData struct {
	Credentials json.RawMessage
	Metadata    map[string]string
	Name        string
	Datacenter  constant.DatacenterType
}
