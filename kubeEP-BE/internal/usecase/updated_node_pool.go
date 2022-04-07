package useCase

import (
	"github.com/google/uuid"
	UCEntity "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/usecase"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/repository"
	"gorm.io/gorm"
)

type UpdatedNodePool interface {
	GetAllUpdatedNodePoolByEvent(
		tx *gorm.DB,
		eventID uuid.UUID,
	) ([]*UCEntity.UpdatedNodePoolData, error)
}

type updatedNodePool struct {
	updatedNodePoolRepo repository.UpdatedNodePool
}

func newUpdatedNodePool(updatedNodePoolRepo repository.UpdatedNodePool) UpdatedNodePool {
	return &updatedNodePool{
		updatedNodePoolRepo: updatedNodePoolRepo,
	}
}

func (u *updatedNodePool) GetAllUpdatedNodePoolByEvent(
	tx *gorm.DB,
	eventID uuid.UUID,
) ([]*UCEntity.UpdatedNodePoolData, error) {
	var output []*UCEntity.UpdatedNodePoolData
	data, err := u.updatedNodePoolRepo.GetAllUpdatedNodePoolByEventID(tx, eventID)
	if err != nil {
		return nil, err
	}
	for _, d := range data {
		output = append(
			output, &UCEntity.UpdatedNodePoolData{
				UpdatedNodePoolID: d.ID.GetUUID(),
				NodePoolName:      d.NodePoolName,
			},
		)
	}
	return output, nil
}
