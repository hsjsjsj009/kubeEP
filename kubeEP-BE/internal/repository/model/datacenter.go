package model

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
)

type Datacenter struct {
	BaseModel
	Name        string
	Credentials gormDatatype.JSON
	Metadata    gormDatatype.JSON
	Datacenter  constant.DatacenterType
}

func (d *Datacenter) TableName() string {
	return "datacenters"
}
