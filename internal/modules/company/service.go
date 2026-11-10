package company

import (
	"avto-crm-api/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyService struct {
	companyRepo *CompanyRepository
	db *gorm.DB
}

func NewCompanyService(companyRepo *CompanyRepository, db *gorm.DB) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
		db: db,
	}
}

func (s *CompanyService) Create(req *CreateCompanyRequest, ownerId string) (*Company, error) {
	oid, err := uuid.Parse(ownerId)

	if err != nil {
		return nil, err
	}

	company := &Company{
		Name: req.Name,
		TypeCompany: req.TypeCompany,
		OwnerID: oid,
	}

	if err := s.companyRepo.Create(company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *CompanyService) Update(req *UpdateCompanyRequest, ownerId string, companyId string) (*Company, error) {
	var ncompany *Company

	err := s.db.Transaction(func(tx *gorm.DB) error {
		oid, err := uuid.Parse(ownerId)

		if err != nil {
			return err
		}

		cid, err := uuid.Parse(companyId)

		if err != nil {
			return err
		}

		company, err := s.companyRepo.FindByIdForUpdate(tx, oid, cid)
		if err != nil {
			return err
		}

		updateMap := make(map[string]interface{})

		if req.Name != nil {
			_, err := s.companyRepo.FindByName(oid, *req.Name)

			if err == nil {
				return utils.ErrNameAlreadyExists
			}

			if !errors.Is(err, utils.ErrRecordNotFound) {
				return err
			}

			updateMap["name"] = *req.Name
		}

		if req.TypeCompany != nil {
			updateMap["type_company"] = *req.TypeCompany
		} 

		ncompany = company

		return tx.Model(company).Updates(updateMap).Error
	})

	return ncompany, err
}

func (s *CompanyService) Delete(id string, ownerId string) error {
	return s.companyRepo.Delete(ownerId, id)
}