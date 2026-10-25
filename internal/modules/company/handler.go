package company

import (
	"avto-crm-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	companyService *CompanyService
}

func NewCompanyHandler(companyService *CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	company, err := h.companyService.Create(&req, ownerId)

	if err != nil {
		utils.CallErrrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "company successfully created", &CreateCompanyResponse{
		Company: *company,
	})
}