package company

import "github.com/google/uuid"

type CompanyService struct {
	companyRepo *CompanyRepository
}

func NewCompanyService(companyRepo *CompanyRepository) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
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