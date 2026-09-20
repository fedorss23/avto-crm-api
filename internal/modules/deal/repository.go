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

func (r *DealRepository) FindById(id, ownerId string) (*Deal, error) {
	var deal Deal

	err := r.db.First(&deal, "id = ? AND owner_id = ?", id, ownerId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &deal, nil
}

func (r *DealRepository) FindFullById(id, ownerId string) (*Deal, error) {
	var deal Deal

	err := r.db.Preload("Pipeline.Stages").Preload("Car").Preload("Client").First(&deal, "id = ? AND owner_id = ?", id, ownerId).Error
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
	
	err := r.db.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Preload("Pipeline.Stages").First(&deal, "id = ? AND owner_id = ?", id, ownerId).Error;
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &deal, nil

}

func (r *DealRepository) TransactionUpdate(tx *gorm.DB, deal *Deal) error {
	for _, i := range deal.Pipeline.Stages {
		if i.ID == *deal.CurrentStageId {
			deal.CurrentStageName = &i.Name
			return tx.Omit("Pipeline", "Cars").Save(deal).Error
		}
	}

	return tx.Omit("Pipeline", "Cars").Save(deal).Error
}

func (r *DealRepository) Update(deal *Deal) error {
	return r.db.Omit("Pipeline", "Cars").Save(deal).Error
}

func (r *DealRepository) FindList(page, limit int) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&Deal{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&deals).Error

	return deals, total, err
}

func (r *DealRepository) FindListWithAll(page, limit int) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&Deal{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Pipeline.Stages").Preload("Car").Preload("Client").Offset(offset).Limit(limit).Order("created_at DESC").Find(&deals).Error

	return deals, total, err
}

func (r *DealRepository) FindByOwnerId(ownerID string, page, limit int, search, status string) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	offset := (page - 1) * limit

	query := r.db.Where("owner_id = ?", ownerID)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status LIKE ?", "%"+status+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&deals).Error

	return deals, total, err
}

func (r *DealRepository) FindFullByOwnerId(ownerID string, page, limit int, search, status string) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	offset := (page - 1) * limit

	query := r.db.Where("owner_id = ?", ownerID)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status LIKE ?", "%"+status+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Pipeline.Stages").Preload("Car").Preload("Client").Offset(offset).Limit(limit).Order("created_at DESC").Find(&deals).Error

	return deals, total, err
}

func (r *DealRepository) FindByClientID(clientId, ownerId string) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	query := r.db.Where("client_id = ? AND owner_id = ?", clientId, ownerId)

	if err := query.Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return deals, 0, err
	}

	return deals, total, nil
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
