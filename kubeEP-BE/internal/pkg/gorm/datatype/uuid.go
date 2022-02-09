package datatype

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UUID uuid.UUID

func (UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "UUID"
}

func (UUID) GormDataType() string {
	return "uuid"
}
