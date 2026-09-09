package deal

import (
	"avto-crm-api/internal/utils"
	"fmt"
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

func (h *DealHandler) FindAll(c *gin.Context) {
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

	deals, total, err := h.dealService.FindAll(page, limit, isFull)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	data := &DealsWithTotal{
		Deals: deals,
		Total: int(total),
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Успешное получение данных: %d", total), data)
}

func (h *DealHandler) CreateFullDeal(c *gin.Context) {
	var req CreateDealRequest

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs, utils.ValidationErrorCode)
		return
	}

	if err := h.dealService.CreateFullDeal(&req, ownerId); err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with creating deal", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusCreated, "deal successfully created")
}

func (h *DealHandler) Update(c *gin.Context) {
	var req Deal

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs, utils.ValidationErrorCode)
		return
	}

	if err := h.dealService.Update(&req, ownerId); err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deal successfully updated", req)
}

func (h *DealHandler) FindDealsByOwnerId(c *gin.Context) {
	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	deals, total, err := h.dealService.FindDealByOwnerId(ownerId)

	resp := DealsResponse{
		Deals: deals,
		Total: total,
	}

	if err != nil {
		if deals != nil {
			utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("error with count total deals: %s", err.Error()), &resp)
			return
		}
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with getting deals by owner id", err, utils.CodeByError(err))
		return
	}

	data := &DealsWithTotal{
		Deals: deals,
		Total: int(total),
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("deals successfully found: %d", total), data)
}

func (h *DealHandler) FindDealByClientId(c *gin.Context) {
	var clientId string
	if err := utils.GetParam(c, "clientId", &clientId); err != nil {
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	deals, total, err := h.dealService.FindDealByClientId(clientId, ownerId)

	resp := DealsResponse{
		Deals: deals,
		Total: total,
	}

	if err != nil {
		if deals != nil {
			utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("error with count total deals: %s", err.Error()), &resp)
			return
		}
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with getting deals by client id", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deals successfully found", resp)
}

func (h *DealHandler) SetNextStage(c *gin.Context) {
	var dealId string
	if err := utils.GetParam(c, "dealId", &dealId); err != nil {
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	deal, err := h.dealService.SetNextStage(ownerId, dealId)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with change stage", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stage has been changed", deal)
}

func (h *DealHandler) Delete(c *gin.Context) {
	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	var dealId string
	if err := utils.GetParam(c, "dealId", &dealId); err != nil {
		return
	}

	err := h.dealService.Delete(ownerId, dealId)

	if err != nil {
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

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	deal, err := h.dealService.FindById(ownerId, dealId, isFull)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "Error with found deal by id", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deal successfully found", deal)
}

func (h *DealHandler) CancelDeal(c *gin.Context) {
	var dealId string
	if err := utils.GetParam(c, "dealId", &dealId); err != nil {
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	err := h.dealService.ChangeStatus(ownerId, dealId, "inactive")

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with cancel deal", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "deal successfully canceled")
}

func (h *DealHandler) ActiveDeal(c *gin.Context) {
	var dealId string
	if err := utils.GetParam(c, "dealId", &dealId); err != nil {
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	err := h.dealService.ChangeStatus(ownerId, dealId, "active")

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with activate deal", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "deal successfully activated")
}

func (h *DealHandler) ChangeStage(c *gin.Context) {
	var dealId string
	if err := utils.GetParam(c, "dealId", &dealId); err != nil {
		return
	}

	var stageId string
	if err := utils.GetStringRequiredQuery(c, "stageId", &stageId); err != nil {
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	err := h.dealService.ChangeStage(ownerId, dealId, stageId)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), "error with change stage of deal", err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "successfully change stage")
}
