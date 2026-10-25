package company

type CreateCompanyRequest struct {
	Name string `json:"name" binding:"required"`
	TypeCompany string `json:"typeCompany" binding:"required"`
}

type CreateCompanyResponse struct {
	Company Company `json:"company"`
}