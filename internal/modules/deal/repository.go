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
	return tx.Omit("Pipeline", "Car").Create(deal).Error
}

func (r *DealRepository) FindById(tx *gorm.DB, id string) (*Deal, error) {
	var deal Deal

	err := tx.First(&deal, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}	

	return &deal, nil
}

func (r *DealRepository) FindFullById(tx *gorm.DB, id string) (*Deal, error) {
	var deal Deal

	err := tx.Preload("Pipeline.Stages").Preload("Car").First(&deal, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}	

	return &deal, nil
}

func (r *DealRepository) Update(tx *gorm.DB, deal *Deal) error {
	return tx.Omit("Pipeline", "Cars").Save(deal).Error
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

func (r *DealRepository) FindListWithAll(tx *gorm.DB, page, limit int) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	offset := (page - 1) * limit

	if err := tx.Model(&Deal{}).Count(&total).Error; err != nil {
		return  nil, 0, err
	}

	err := tx.Preload("Pipeline.Stages").Preload("Car").Offset(offset).Limit(limit).Order("created_at DESC").Find(&deals).Error

	return deals, total, err
}

func (r *DealRepository) FindByOwnerId(tx *gorm.DB, ownerID string) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	query := tx.Where("owner_id = ?", ownerID)

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

	query := tx.Where("client_id = ?", clientID)

	if err := query.Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return deals, 0, err
	}

	return deals, total, nil
}

func (r *DealRepository) SetNextStage(tx *gorm.DB, ownerID, dealID string) (*Deal, error) {
	var deal Deal

	if err := tx.Preload("Pipeline.Stages").Where("owner_id = ? AND id = ?", ownerID, dealID).First(&deal).Error; err != nil {
		return nil, err
	}

	if deal.Pipeline == nil {
		return nil, ErrNotPipeline
	}

	if deal.Pipeline.Stages == nil {
		return nil, ErrNotStages
	}

	stages := deal.Pipeline.Stages

	if len(stages) == 0 {
		return nil, ErrEmptyPipeline
	}

	if deal.CurrentStage == nil {
		if err := tx.Model(&deal).Update("current_stage", stages[0].ID).Error; err != nil {
			return nil, err
		}

		deal.CurrentStage = &stages[0].ID
		return &deal, nil
	}

	for i := 0; i < len(stages); i++ {
		if stages[i].ID.String() == deal.CurrentStage.String() {
			if i == len(stages) - 1 {
				return &deal, ErrLastStage
			}
			if err := tx.Model(&deal).Update("current_stage", stages[i+1].ID).Error; err != nil {
				return nil, err
			}
			deal.CurrentStage = &stages[i+1].ID
			return &deal, nil
		}
	}

	return nil, errors.New("Server error")
}

func (r *DealRepository) Delete(tx *gorm.DB, ownerID, dealID string) error {
	var deal Deal

	if err := tx.Where("id = ? AND owner_id = ?", dealID, ownerID).First(&deal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}

		return err
	}

	return tx.Delete(&Deal{}, "id = ?", dealID).Error
}