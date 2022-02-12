package model

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
)

type Cluster struct {
	BaseModel
	DatacenterID   gormDatatype.UUID
	Metadata       gormDatatype.JSON
	Name           string
	Certificate    []byte
	ServerEndpoint string
	Datacenter     Datacenter `gorm:"ForeignKey:DatacenterID;constraint:OnDelete:CASCADE"`
}

func (c *Cluster) TableName() string {
	return "clusters"
}
