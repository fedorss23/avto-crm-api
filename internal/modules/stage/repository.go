package stage

import (
	"avto-crm-api/internal/utils"

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


func (r *StageRepository) DeleteWithPipeline(tx *gorm.DB, id uuid.UUID) error {
	query := tx.Where("id = ?", id)

	result := query.Delete(&Stage{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *StageRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	query := tx.Where("id = ?", id)

	result := query.Unscoped().Delete(&Stage{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}