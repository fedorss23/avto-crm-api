package deal

import (
	"avto-crm-api/internal/utils"
	"errors"

	"github.com/google/uuid"
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

func (r *DealRepository) FindById(id string, scope *utils.AuthScope, isFull bool) (*Deal, error) {
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

func (r *DealRepository) FindForUpdate(tx *gorm.DB, id string, scope *utils.AuthScope) (*Deal, error) {
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

func (r *DealRepository) FindForDelete(tx *gorm.DB, id string, scope *utils.AuthScope) (*Deal, error) {
	var deal Deal

	query := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Preload("Pipeline.Stages").Preload("Client").Preload("Car").Where("id = ? AND company_id = ?", id, scope.CompanyId)

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

func (r *DealRepository) GetTotalForStatus(scope *utils.AuthScope, status string) (int64, error) {
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

func (r *DealRepository) Delete(tx *gorm.DB, dealId string, scope *utils.AuthScope, fullDelete bool) error {
	query := tx.Where(
		"id = ? AND company_id = ?",
		dealId,
		scope.CompanyId,
	)

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", *scope.ManagerId)
	}

	var result *gorm.DB
	if fullDelete == true {
		result = query.Unscoped().Delete(&Deal{})
	} else {
		result = query.Delete(&Deal{})
	}

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}

func (r *DealRepository) GetFreeManager(tx *gorm.DB, companyId uuid.UUID) (string, error) {
	if err := r.CheckManagers(tx, companyId); err != nil {
		return "", err
	}

	mid := struct{
		ID uuid.UUID
	}{}
	err := tx.Table("managers").
		Joins("LEFT JOIN deals ON deals.manager_id = managers.id").
		Where("managers.company_id = ?", companyId).
		Group("managers.id").
		Order("COUNT(deals.id) ASC",).
		First(&mid).Error

	if err != nil {
		return "", err
	}


	return mid.ID.String(), nil
}

func (r *DealRepository) CheckManagers(tx *gorm.DB, companyId uuid.UUID) error {
	var count int64
	if err := tx.Table("managers").Where("company_id = ?::uuid", companyId).Count(&count).Error; err != nil || count == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotManagers
		}
		return err
	}
	return nil
}