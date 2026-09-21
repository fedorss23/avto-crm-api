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

func (r *DealRepository) FindById(id, ownerId string, isFull bool) (*Deal, error) {
	var deal Deal

	query := r.db.Where("owner_id = ? AND id = ?", ownerId, id)

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

func (r *DealRepository) FindFullForUpdate(tx *gorm.DB, id, ownerId string) (*Deal, error) {
	var deal Deal

	err := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Preload("Pipeline.Stages").First(&deal, "id = ? AND owner_id = ?", id, ownerId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &deal, nil

}

func (r *DealRepository) GetTotalForStatus(ownerId string, status string) (int64, error) {
	var total int64

	if err := r.db.Model(&Deal{}).Where("owner_id = ? AND status = ?", ownerId, status).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *DealRepository) FindList(filters *DealFilters) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	offset := (filters.page - 1) * filters.limit

	query := r.db.Model(&Deal{})

	if filters.ownerId != "" {
		query = query.Where("owner_id = ?", filters.ownerId)
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

func (r *DealRepository) Delete(ownerId, dealId string) error {
	result := r.db.Delete(
		&Deal{},
		"id = ? AND owner_id = ?",
		dealId,
		ownerId,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}
