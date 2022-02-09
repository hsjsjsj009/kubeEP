package model

import "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"

type Datacenter struct {
	BaseModel
	Name        string
	Credentials datatype.JSON
	Metadata    datatype.JSON
}

func (d *Datacenter) TableName() string {
	return "datacenters"
}
