package model

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
)

type Datacenter struct {
	BaseModel
	Name        string                  `json:"name"`
	Credentials gormDatatype.JSON       `json:"credentials"`
	Metadata    gormDatatype.JSON       `json:"metadata"`
	Datacenter  constant.DatacenterType `json:"datacenter"`
}

func (d *Datacenter) TableName() string {
	return "datacenters"
}
