package stage

import (
	"gorm.io/gorm"
)

type StageRepository struct {}

func NewStageRepository() *StageRepository {
	return &StageRepository{}
}

func (r *StageRepository) Create(tx *gorm.DB, stage *Stage) error {
	return tx.Create(stage).Error
}
