package model

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        gormDatatype.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time         `json:"-"`
	UpdatedAt time.Time         `json:"-"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`
}
