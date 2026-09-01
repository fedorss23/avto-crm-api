package stage

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StageRepository struct {}

func NewStageRepository() *StageRepository {
	return &StageRepository{}
}

func (r *StageRepository) Create(tx *gorm.DB, stage *Stage) error {
	return tx.Create(stage).Error
}

func (r *StageRepository) Update(tx *gorm.DB, stage *Stage) error {
	return tx.Save(stage).Error
}

func (r *StageRepository) Delete(tx *gorm.DB, stageID uuid.UUID) error {
	return tx.Delete(&Stage{}, "WHERE id = ?", stageID).Error
}