package company

import (
	"avto-crm-api/internal/utils"
	"errors"

	"gorm.io/gorm"
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

func (r *CompanyRepository) FindById(ownerId string, companyId string) (*Company, error) {
	var company Company
	if err := r.db.Where("owner_id = ? AND id = ?", ownerId, companyId).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &company, nil
}