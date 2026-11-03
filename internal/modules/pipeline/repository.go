package pipeline

import (
	"avto-crm-api/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PipelineRepository struct{}

func NewPipelineRepository() *PipelineRepository {
	return &PipelineRepository{}
}

func (r *PipelineRepository) Create(tx *gorm.DB, pipeline *Pipeline) error {
	return tx.Create(pipeline).Error
}


func (r *PipelineRepository) FindById(tx *gorm.DB, id uuid.UUID) (*Pipeline, error) {
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

func (r *PipelineRepository) Delete(tx *gorm.DB, id uuid.UUID, scope *utils.AuthScope) (error) {
	query := tx.Where("id = ?", id)

	if scope.CompanyId != nil {
		query = query.Where("company_id = ?", scope.CompanyId)
	}

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", scope.ManagerId)
	}

	result := query.Delete(&Pipeline{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}

func (r *PipelineRepository) FindForUpdate(tx *gorm.DB, id uuid.UUID, scope *utils.AuthScope, isFull bool) (*Pipeline, error) {
	var pipeline Pipeline
	query := tx.Clauses(&clause.Locking{
		Strength: "UPDATE",
	}).Where("id = ?", id)

	if scope.CompanyId != nil {
		query = query.Where("company_id = ?", *scope.CompanyId)
	}

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", *scope.ManagerId)
	}

	if isFull {
		query = query.Preload("Stages")
	}

	err := query.First(&pipeline).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &pipeline, nil
}

