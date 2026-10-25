package deal

import (
	"avto-crm-api/internal/utils"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DealRepository struct {
	db *gorm.DB
}

func NewDealRepository(db *gorm.DB) *DealRepository {
	return &DealRepository{
		db: db,
	}
}

func (r *DealRepository) Create(tx *gorm.DB, deal *Deal) error {
	return tx.Create(deal).Error
}

func (r *DealRepository) FindById(id string, scope *DealScope, isFull bool) (*Deal, error) {
	var deal Deal

	query := r.db.Where("id = ? AND company_id = ?", id, scope.CompanyId)

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", *scope.ManagerId)
	}

	if isFull {
		query = query.Preload("Pipeline.Stages").Preload("Car").Preload("Client")
	}

	err := query.First(&deal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &deal, nil
}

func (r *DealRepository) FindFullForUpdate(tx *gorm.DB, id string, scope *DealScope) (*Deal, error) {
	var deal Deal

	query := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Preload("Pipeline.Stages").Where("id = ? AND company_id = ?", id, scope.CompanyId)

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", *scope.ManagerId)
	}

	err := query.First(&deal).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &deal, nil
}

func (r *DealRepository) GetTotalForStatus(scope *DealScope, status string) (int64, error) {
	var total int64

	query := r.db.Model(&Deal{}).Where("company_id = ?", scope.CompanyId)

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", *scope.ManagerId)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Model(&Deal{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *DealRepository) FindList(filters *DealFilters) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	offset := (filters.page - 1) * filters.limit

	query := r.db.Model(&Deal{}).Where("company_id = ?", filters.scope.CompanyId)

	if filters.scope.ManagerId != nil{
		query = query.Where("manager_id = ?", *filters.scope.ManagerId)
	}

	if filters.clientId != "" {
		query = query.Where("client_id = ?", filters.clientId)
	}

	if filters.search != "" {
		query = query.Where("name ILIKE ?", "%"+filters.search+"%")
	}

	if filters.status != "" {
		query = query.Where("status = ?", filters.status)
	}

	if filters.isFull {
		query = query.Preload("Pipeline.Stages").Preload("Car").Preload("Client")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(filters.limit).Order("created_at DESC").Find(&deals).Error

	return deals, total, err
}

func (r *DealRepository) Delete(dealId string, scope *DealScope) error {
	query := r.db.Where(
		"id = ? AND company_id = ?",
		dealId,
		scope.CompanyId,
	)

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", *scope.ManagerId)
	}

	result := query.Delete(&Deal{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}
