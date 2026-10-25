package deal

import (
	"avto-crm-api/internal/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DealHandler struct {
	dealService *DealService
}

func NewDealHandler(dealService *DealService) *DealHandler {
	return &DealHandler{
		dealService: dealService,
	}
}

func (h *DealHandler) CreateFullDeal(c *gin.Context) {
	var req CreateDealRequest

	scope, err := h.getDealScopeWithQueryManager(c)
	if err != nil {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	deal, err := h.dealService.CreateFullDeal(&req, scope)
	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with creating deal", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "deal successfully created", &OneDealResponse{
		Deal: *deal,
	})
}

func (h *DealHandler) Update(c *gin.Context) {
	var req UpdateDealRequest

	scope, err := h.getDealScope(c)
	if err != nil {
		return
	}

	var dealId string
	if err := utils.GetParam(c, "dealId", &dealId); err != nil {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	deal, err := h.dealService.Update(&req, dealId, scope)
	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deal successfully updated", &OneDealResponse{
		Deal: *deal,
	})
}

func (h *DealHandler) GetTotalByStatus(c *gin.Context) {
	scope, err := h.getDealScopeWithQueryManager(c)
	if err != nil {
		return
	}

	status, exists := c.GetQuery("status")
	if exists {
		if status != "inactive" && status != "active" && status != "completed" && status != "" {
			utils.ValidationErrorResponse(
				c,
				map[string]string{"status": "status must be: inactive, active or completed"},
			)
			return
		}
	}

	total, err := h.dealService.GetTotalByStatus(scope, status)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with getting deals by owner id", err, utils.CodeByError(err))
		return
	}

	data := &TotalData{
		Total:  total,
		Status: status,
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("deals successfully found: %d", total), data)
}

func (h *DealHandler) FindList(c *gin.Context) {
	scope, err := h.getDealScopeWithQueryManager(c)
	if err != nil {
		return
	}

	var page int
	if err := utils.GetNumberQuery(c, "page", &page, 1, utils.PageErrorCode); err != nil {
		return
	}

	var limit int
	if err := utils.GetNumberQuery(c, "limit", &limit, 10, utils.LimitErrorCode); err != nil {
		return
	}

	var isFull bool
	if err := utils.GetBoolQuery(c, "isFull", &isFull, false, utils.IsFullErrorCode); err != nil {
		return
	}

	status, exists := c.GetQuery("status")
	if exists {
		if status != "inactive" && status != "active" && status != "completed" && status != "" {
			utils.ValidationErrorResponse(
				c,
				map[string]string{"status": "status must be: inactive, active or completed"},
			)
			return
		}
	}

	search, _ := c.GetQuery("search")
	clientId, _ := c.GetQuery("clientId")

	deals, total, err := h.dealService.FindDealByOwnerId(&DealFilters{
		scope: scope,
		clientId:  clientId,
		page:      page,
		limit:     limit,
		status:    status,
		search:    search,
		isFull:    isFull,
	})

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with getting deals by owner id", err, utils.CodeByError(err))
		return
	}

	data := &DealsWithTotal{
		Deals: deals,
		Total: int(total),
		Page:  page,
		Limit: limit,
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("deals successfully found: %d", total), data)
}

func (h *DealHandler) Delete(c *gin.Context) {
	scope, err := h.getDealScope(c)
	if err != nil {
		return
	}

	var dealId string
	if err := utils.GetParam(c, "dealId", &dealId); err != nil {
		return
	}

	if err := h.dealService.Delete(dealId, scope); err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "Error with delete deal", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "deal successfully deleted")
}

func (h *DealHandler) FindById(c *gin.Context) {
	var dealId string
	if err := utils.GetParam(c, "dealId", &dealId); err != nil {
		return
	}

	var isFull bool
	if err := utils.GetBoolQuery(c, "isFull", &isFull, false, utils.IsFullErrorCode); err != nil {
		return
	}

	scope, err := h.getDealScope(c)
	if err != nil {
		return
	}

	deal, err := h.dealService.FindById(dealId, scope, isFull)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "Error with found deal by id", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deal successfully found", &OneDealResponse{
		Deal: *deal,
	})
}

func (h *DealHandler) getDealScope(c *gin.Context) (*DealScope, error) {
	var scope DealScope

	var userId string
	userType, err := utils.GetUserId(c, &userId)
	if err != nil {
		return nil, err
	}

	scope.CompanyId = c.GetString("companyId")

	if userType == "manager" {
		scope.ManagerId = &userId
	}

	return &scope, nil
}

func (h *DealHandler) getDealScopeWithQueryManager(c *gin.Context) (*DealScope, error) {
	var scope DealScope

	var userId string
	userType, err := utils.GetUserId(c, &userId)
	if err != nil {
		return nil, err
	}

	log.Println(userId)

	scope.CompanyId = c.GetString("companyId")
	
	if userType == "manager" {
		scope.ManagerId = &userId
	} else {
		manager, ok := c.GetQuery("managerId")

		if ok && manager != "" {
			scope.ManagerId = &manager
		} 
	}

	return &scope, nil
}