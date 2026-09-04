package deal

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (r *DealRepository) FindById(id string) (*Deal, error) {
	var deal Deal

	err := r.db.First(&deal, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	return &deal, nil
}

func (r *DealRepository) FindFullById(id string) (*Deal, error) {
	var deal Deal

	err := r.db.Preload("Pipeline.Stages").Preload("Car").Preload("Client").First(&deal, "id = ?", id).Error
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

func (r *DealRepository) FindByOwnerId(ownerID string) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	query := r.db.Where("owner_id = ?", ownerID).Order("created_at DESC")

	if err := query.Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return deals, 0, err
	}

	return deals, total, nil
}

func (r *DealRepository) FindByClientID(clientID string) ([]Deal, int64, error) {
	var deals []Deal
	var total int64

	query := r.db.Where("client_id = ?", clientID)

	if err := query.Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return deals, 0, err
	}

	return deals, total, nil
}

func (r *DealRepository) SetNextStage(ownerID, dealID string) (*Deal, error) {
	var deal Deal

	if err := r.db.Preload("Pipeline.Stages").Where("owner_id = ? AND id = ?", ownerID, dealID).First(&deal).Error; err != nil {
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

	if deal.CurrentStageId == nil {
		if err := r.db.Model(&deal).Update("current_stage", stages[0].ID).Error; err != nil {
			return nil, err
		}

		deal.CurrentStageId = &stages[0].ID
		return &deal, nil
	}

	for i := 0; i < len(stages); i++ {
		if stages[i].ID.String() == deal.CurrentStageId.String() {
			if i == len(stages)-1 {
				return &deal, ErrLastStage
			}
			if err := r.db.Model(&deal).Update("current_stage", stages[i+1].ID).Error; err != nil {
				return nil, err
			}
			deal.CurrentStageId = &stages[i+1].ID
			return &deal, nil
		}
	}

	return nil, errors.New("Server error")
}

func (r *DealRepository) Delete(ownerID, dealID string) error {
	var deal *Deal

	if err := r.db.Where("id = ? AND owner_id = ?", dealID, ownerID).First(&deal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}

		return err
	}

	return r.db.Delete(&Deal{}, "id = ?", dealID).Error
}

func (r *DealRepository) ChangeStatus(ownerId, dealId, status string) error {
	var deal *Deal

	if err := r.db.Where("id = ? AND owner_id = ?", dealId, ownerId).First(&deal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}

		return err
	}

	deal.Status = status
	return r.db.Save(deal).Error
}

func (r *DealRepository) SetStage(ownerId, dealId, stageId string) error {
	var deal *Deal

	if err := r.db.Preload("Pipeline.Stages").Where("id = ? AND owner_id = ?", dealId, ownerId).First(&deal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}

		return err
	}

	parsedId, err := uuid.Parse(stageId)

	if err != nil {
		return err
	}

	for _, i := range deal.Pipeline.Stages {
		if i.ID == parsedId {
			deal.CurrentStageId = &parsedId
			deal.CurrentStageName = &i.Name
			return r.db.Save(deal).Error
		}
	}

	return ErrRecordNotFound
}
