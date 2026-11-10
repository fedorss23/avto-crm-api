package pipeline

import (
	"avto-crm-api/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PipelineRepository struct{
	db *gorm.DB
}

func NewPipelineRepository(db *gorm.DB) *PipelineRepository {
	return &PipelineRepository{
		db: db,
	}
}

func (r *PipelineRepository) Create(tx *gorm.DB, pipeline *Pipeline) error {
	return tx.Create(pipeline).Error
}

func (r *PipelineRepository) FindList(tx *gorm.DB, page, limit int, scope *utils.AuthScope) ([]Pipeline, int64, error) {
	var pipelines []Pipeline
	var total int64

	offset := (page - 1) * limit

	query := tx.Model(&Pipeline{}).
		Joins("JOIN deals ON deals.id = pipelines.deal_id").
		Where("deals.company_id = ?", scope.CompanyId)

	if scope.ManagerId != nil {
		query = query.Where("deals.manager_id = ?", *scope.ManagerId)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Stages").Offset(offset).Limit(limit).Order("pipelines.created_at DESC").Find(&pipelines).Error

	return pipelines, total, err
}

func (r *PipelineRepository) DeleteWithDeal(tx *gorm.DB, id uuid.UUID) error {
	query := tx.Where("id = ?", id)

	result := query.Delete(&Pipeline{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PipelineRepository) Delete(id string, scope *utils.AuthScope) error {
	result := r.db.
		Where("id = ?", id).
		Where("EXISTS(?)", r.checkAccessQuery(r.db, scope)).
		Unscoped().
		Delete(&Pipeline{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}

func (r *PipelineRepository) FindForUpdate(tx *gorm.DB, id string, scope *utils.AuthScope, isFull bool) (*Pipeline, error) {
	var pipeline Pipeline
	query := tx.Clauses(&clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		Where("EXISTS(?)", r.checkAccessQuery(tx, scope))

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

func (r *PipelineRepository) checkAccessQuery(tx *gorm.DB, scope *utils.AuthScope) *gorm.DB {
	subQuery := tx.Table("deals").
		Select("1").
		Where("deals.id = pipelines.deal_id").
		Where("deals.company_id = ?", scope.CompanyId)

	if scope.ManagerId != nil {
		subQuery = subQuery.Where("deals.manager_id = ?", *scope.ManagerId)
	}
	return subQuery
}
