package stage

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StageRepository struct {
	db *gorm.DB
}

func NewStageRepository(db *gorm.DB) *StageRepository {
	return &StageRepository{
		db: db,
	}
}

func (r *StageRepository) Create(tx *gorm.DB, stage *Stage) error {
	return tx.Create(stage).Error
}

func (r *StageRepository) Update(data map[string]interface{}, stageId string, ownerId string) (*Stage, error) {
	var stage *Stage
	if err := r.db.Model(stage).Where("id = ? AND owner_id = ?", stageId, ownerId).Updates(data).Error; err != nil {
		return nil, err
	}

	return stage, nil
}

func (r *StageRepository) Delete(tx *gorm.DB, stageID uuid.UUID) error {
	return tx.Delete(&Stage{}, "WHERE id = ?", stageID).Error
}
