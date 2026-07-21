package deal

import (
	"errors"
	"gorm.io/gorm"
)

type DealRepository struct {}

func NewDealRepository(db *gorm.DB) *DealRepository {
	return &DealRepository{}
}

func (r *DealRepository) Create(tx *gorm.DB, deal *Deal) error {
	return tx.Create(deal).Error
}


func (r *DealRepository) FindById(tx *gorm.DB, id string) (*Deal, error) {
	var deal Deal

	err := tx.First(&deal, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}	

	return &deal, nil
}

func (r *DealRepository) Update(tx *gorm.DB, deal *Deal) error {
	return tx.Save(deal).Error
} 

func (r *DealRepository) Delete(tx *gorm.DB, id string) error {
	return tx.Delete(&Deal{}, "id = ?", id).Error
}

func (r *DealRepository) FindList(tx *gorm.DB, page, limit int) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	offset := (page-1)*limit

	if err := tx.Model(&Deal{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := tx.Offset(offset).Limit(limit).Order("created_at DESC").Find(&deals).Error

	return deals, total, err
}

func (r *DealRepository) FindByOwnerId(tx *gorm.DB, ownerID string) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	query := tx.Where("ownerId = ?", ownerID)

	if err := query.Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return deals, 0, err
	}

	return deals, total, nil
}

func (r *DealRepository) FindByClientID(tx *gorm.DB, clientID string) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	query := tx.Where("clientId = ?", clientID)

	if err := query.Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return deals, 0, err
	}

	return deals, total, nil
}