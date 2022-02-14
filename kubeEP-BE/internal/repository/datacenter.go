package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	gormDatatype "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/gorm/datatype"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository/model"
	"gorm.io/gorm"
	"time"
)

type Datacenter interface {
	GetDatacenterByID(tx *gorm.DB, id uuid.UUID) (*model.Datacenter, error)
	InsertDatacenter(tx *gorm.DB, data *model.Datacenter) error
	InsertTemporaryDatacenter(ctx context.Context, data *model.Datacenter, exp time.Duration) error
	GetTemporaryDatacenterByID(ctx context.Context, id uuid.UUID) (*model.Datacenter, error)
}

type datacenter struct {
	redisClient *redis.Client
}

func NewDatacenter(redisClient *redis.Client) Datacenter {
	return &datacenter{
		redisClient: redisClient,
	}
}

func (d *datacenter) InsertTemporaryDatacenter(ctx context.Context, data *model.Datacenter, exp time.Duration) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	rdClient := d.redisClient

	data.ID = gormDatatype.UUID(id)

	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	status := rdClient.Set(ctx, fmt.Sprintf("datacenter_%s", id.String()), byteData, exp)
	if err := status.Err(); err != nil {
		return err
	}

	return nil
}

func (d *datacenter) GetTemporaryDatacenterByID(ctx context.Context, id uuid.UUID) (*model.Datacenter, error) {
	rdClient := d.redisClient
	val := rdClient.Get(ctx, fmt.Sprintf("datacenter_%s", id.String()))
	if err := val.Err(); err != nil {
		return nil, err
	}
	dataByte, err := val.Bytes()
	if err != nil {
		return nil, err
	}
	data := &model.Datacenter{}
	err = json.Unmarshal(dataByte, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *datacenter) GetDatacenterByID(tx *gorm.DB, id uuid.UUID) (*model.Datacenter, error) {
	data := &model.Datacenter{}
	tx = tx.First(data, id)
	return data, tx.Error
}

func (d *datacenter) InsertDatacenter(tx *gorm.DB, data *model.Datacenter) error {
	return tx.Create(data).Error
}
