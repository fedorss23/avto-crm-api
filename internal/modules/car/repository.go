package car

import (
	"avto-crm-api/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CarRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

func (r *CarRepository) Create(tx *gorm.DB, car *Car) error {
	return tx.Create(car).Error
}

func (r *CarRepository) FindList(page, limit int, scope *utils.AuthScope) ([]Car, int64, error) {
	var cars []Car
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&Car{}).
		Joins("JOIN deals ON deals.id = cars.deal_id").
		Where("deals.company_id = ?", scope.CompanyId)

	if scope.ManagerId != nil {
		query = query.Where("deals.manager_id = ?", scope.ManagerId)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&cars).Error

	return cars, total, err
}

func (r *CarRepository) FindForUpdate(tx *gorm.DB, id string, scope *utils.AuthScope) (*Car, error) {
	var car Car
	query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Joins("JOIN deals ON deals.id = cars.deal_id").
		Where("cars.id = ?", id).
		Where("deals.company_id = ?", scope.CompanyId)

	if scope.ManagerId != nil {
		query = query.Where("deals.manager_id = ?", *scope.ManagerId)
	}

	err := query.First(&car).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &car, nil
}

func (r *CarRepository) DeleteWithDeal(tx *gorm.DB, id uuid.UUID) error {
	query := tx.Where("id = ?", id)

	result := query.Delete(&Car{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CarRepository) Delete(id string, scope *utils.AuthScope) error {
	query := r.db.Model(&Car{}).
		Joins("JOIN deals ON deals.id = cars.deal_id").
		Where("cars.id = ?", id).
		Where("deals.company_id = ?", scope.CompanyId)

	if scope.ManagerId != nil {
		query = query.Where("deals.manager_id = ?", *scope.ManagerId)
	}

	result := query.Unscoped().Delete(&Car{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}