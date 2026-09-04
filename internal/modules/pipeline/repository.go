package pipeline

import (
	"errors"
	"gorm.io/gorm"
)

type PipelineRepository struct{}

func NewPipelineRepository() *PipelineRepository {
	return &PipelineRepository{}
}

func (r *PipelineRepository) Create(tx *gorm.DB, pipeline *Pipeline) error {
	return tx.Create(pipeline).Error
}

func (r *PipelineRepository) Update(tx *gorm.DB, pipeline *Pipeline) error {
	return tx.Save(pipeline).Error
}

func (r *PipelineRepository) FindById(tx *gorm.DB, id string) (*Pipeline, error) {
	var pipeline Pipeline

	if err := tx.Where("id = ?", id).First(&pipeline).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &pipeline, nil
}

func (r *PipelineRepository) FindList(tx *gorm.DB, page, limit int) ([]Pipeline, int64, error) {
	var pipelines []Pipeline
	var total int64

	offset := (page - 1) * limit

	if err := tx.Model(&Pipeline{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := tx.Preload("Stages").Offset(offset).Limit(limit).Order("created_at DESC").Find(&pipelines).Error

	return pipelines, total, err
}

