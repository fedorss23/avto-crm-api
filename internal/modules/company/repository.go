package company

import (
	"avto-crm-api/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{
		db: db,
	}
}

func (r *CompanyRepository) Create(company *Company) error {
	return r.db.Create(&company).Error
}

func (r *CompanyRepository) FindById(ownerId uuid.UUID, companyId uuid.UUID) (*Company, error) {
	var company Company
	if err := r.db.Where("owner_id = ? AND id = ?", ownerId, companyId).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) FindFullById(ownerId, companyId string) (*Company, error) {
	var company Company
	if err := r.db.Preload("Managers").Where("owner_id = ? AND id = ?", ownerId, companyId).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) FindByIdForUpdate(tx *gorm.DB, ownerId uuid.UUID, companyId uuid.UUID) (*Company, error) {
	var company Company

	err := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("owner_id = ? AND id = ?", ownerId, companyId).First(&company).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) FindByName(ownerId uuid.UUID, name string) (*Company, error) {
	var company Company
	if err := r.db.Where("owner_id = ? AND name = ?", ownerId, name).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) FindByNameWithTx(tx *gorm.DB, ownerId uuid.UUID, name string) (*Company, error) {
	var company Company
	if err := tx.Where("owner_id = ? AND name = ?", ownerId, name).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) FindList(ownerId uuid.UUID) ([]Company, error) {
	var companies []Company

	if err := r.db.Where("owner_id = ?", ownerId).Find(&companies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return companies, nil
}

func (r *CompanyRepository) Delete(ownerId, companyId string) error {
	result := r.db.Where("id = ? AND owner_id = ?", companyId, ownerId).Unscoped().Delete(&Company{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}