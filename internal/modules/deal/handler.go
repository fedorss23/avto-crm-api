package deal

import (
	"avto-crm-api/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	var pageInt int
	
	page, exists := c.GetQuery("page")
	if !exists {
		pageInt = 1
	} else {
		a, err := strconv.Atoi(page)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Error with found deal: query-param page must be number", errors.New("Page must be number"))
			return
		}
		pageInt = a
	}

	var limitInt int

	limit, exists := c.GetQuery("limit")
	if !exists {
		limitInt = 10
	} else {
		a, err := strconv.Atoi(limit)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Error with found deal: query-param limit must be number", errors.New("Limit must be number"))
			return
		}
		limitInt = a
	}

	isFull, exists := c.GetQuery("isFull")

	var isFullBool bool

	if !exists {
		isFullBool = false
	} else {
		a, err := strconv.ParseBool(isFull)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Параметр isFull должен быть true или false", errors.New("Error with query param isFull"))
			return
		}
		isFullBool = a
	}

	var deals []Deal
	var total int64
	var ferr error

	deals, total, ferr = h.dealService.FindAll(pageInt, limitInt, isFullBool)

	if ferr != nil {
		utils.ErrorResponse(c, ErrorToHTTPStatus(ferr), "Error with getting deals", ferr)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Успешное получение данных: %d", total), deals)
}

func (h *DealHandler) CreateFullDeal(c *gin.Context) {
	var req CreateDealRequest

	ownerId := c.GetString("userId")

	if ownerId == "" {
		utils.ErrorResponse(c, http.StatusForbidden, "you are not authorized", errors.New("not authorized"))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	if err := h.dealService.CreateFullDeal(&req, ownerId); err != nil {
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), "error with creating deal", err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusCreated, "deal successfully created")
}

func (h *DealHandler) Update(c *gin.Context) {
	var req Deal

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	if err := h.dealService.Update(&req); err != nil {
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), "error with updating deal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "deal successfully updated", req)
}

func (h *DealHandler) FindDealByOwnerId(c *gin.Context) {
	id := c.GetString("userId")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "you need to provide a query param: userId", errors.New("missing userId"))
		return
	}

	deals, total, err := h.dealService.FindDealByOwnerId(id)

	

	if err != nil {
		if deals != nil {
			utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("error with count total deals: %s", err.Error()), deals)
			return
		}
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), "error with getting deals by owner id", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("deals successfully found: %d", total), deals)
}

func (h *DealHandler) FindDealByClientId(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "you need to provide a query param: userId", errors.New("missing userId"))
		return
	}

	deals, total, err := h.dealService.FindDealByClientId(id)

	resp := DealsResponse{
		Deals: deals,
		Total: total,
	}

	if err != nil {
		if deals != nil {
			resp.Error = err
			utils.SuccessResponse(c, http.StatusOK, "error with count total deals", &resp)
			return
		}
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), "error with getting deals by client id", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deals successfully found", resp)
}

func (h *DealHandler) SetNextStage(c *gin.Context) {
	dealId := c.Param("dealId")

	if dealId == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Param dealId must be valid", errors.New("Missed dealId"))
		return
	}

	userId := c.GetString("userId")

	if userId == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "error with authorization", errors.New("you aren't authorized"))
	}

	deal, err := h.dealService.SetNextPage(userId, dealId)

	if err != nil {
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), "error with change stage", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stage has been changed", deal)
}

func (h *DealHandler) Delete(c *gin.Context) {
	ownerId := c.GetString("userId")

	if ownerId == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "error with authorization", errors.New("you aren't authorized"))
		return
	}

	dealId := c.Param("dealId")

	if dealId == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Param dealId must be valid", errors.New("Missed dealId"))
		return
	}

	err := h.dealService.Delete(ownerId, dealId)

	if err != nil {
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), "Error with delete deal", err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "deal successfully deleted")
}

func (h *DealHandler) FindById(c *gin.Context) {
	id := c.Param("dealId")

	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Param dealId must be valid", errors.New("Missed dealId"))
		return
	}

	isFull, exists := c.GetQuery("isFull")

	var isFullBool bool

	if !exists {
		isFullBool = false
	} else {
		a, err := strconv.ParseBool(isFull)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Param isFull must be true or false", errors.New("Error with query param isFull"))
			return
		}
		isFullBool = a
	}

	deal, err := h.dealService.FindById(id, isFullBool)

	if err != nil {
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), "Error with found deal by id", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deal successfully found", deal)
}

func (h *DealHandler) CancelDeal(c *gin.Context) {
	dealId := c.Param("dealId")

	if dealId == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Param dealId must be valid", errors.New("Missed dealId"))
		return
	}

	ownerId := c.GetString("userId")

	if ownerId == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "error with authorization", errors.New("you aren't authorized"))
		return
	}

	err := h.dealService.ChangeStatus(ownerId, dealId, "inactive")

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with cancel deal", errors.New("error with cancel deal"))
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "deal successfully canceled")
}

func (h *DealHandler) ActiveDeal(c *gin.Context) {
	dealId := c.Param("dealId")

	if dealId == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Param dealId must be valid", errors.New("Missed dealId"))
		return
	}

	ownerId := c.GetString("userId")

	if ownerId == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "error with authorization", errors.New("you aren't authorized"))
		return
	}

	err := h.dealService.ChangeStatus(ownerId, dealId, "active")

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with activate deal", err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "deal successfully activated")
}

func (h *DealHandler) ChangeStage(c *gin.Context) {
	dealId := c.Param("dealId")

	if dealId == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Param dealId must be valid", errors.New("Missed dealId"))
		return
	}

	stageId, ok := c.GetQuery("stageId")

	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "Query param stageId must be valid", errors.New("Missed stageId"))
		return
	}

	ownerId := c.GetString("userId")

	if ownerId == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "error with authorization", errors.New("you aren't authorized"))
		return
	}

	err := h.dealService.ChangeStage(ownerId, dealId, stageId)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with change stage of deal", err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "successfully change stage")
}