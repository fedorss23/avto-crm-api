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
	service *DealService
}

func NewDealHandler(service *DealService) *DealHandler {
	return &DealHandler{
		service: service,
	}
}

func (h *DealHandler) FindAll(c *gin.Context) {
	page, exists := c.GetQuery("page")
	if !exists {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка при получении сделок: пропущен query-параметер page", errors.New("Error with query param: page"))
		return
	}

	limit, exists := c.GetQuery("limit")
	if !exists {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка при получении сделок: пропущен query-параметер limit", errors.New("Error with query param: limit"))
		return
	}

	intPage, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка при получении сделок: параметр page должен быть int", errors.New("Query param page must be a int"))
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка при получении сделок: параметр page должен быть int", errors.New("Query param page must be a int"))
		return
	}

	deals, total, err := h.service.FindAll(intPage, intLimit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Не удалось получить все сделки", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Успешное получение данных: %d", total), deals)
}

func (h *DealHandler) CreateFullDeal(c *gin.Context) {
	var req CreateDealRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "body is invalid", err)
		return
	}

	if err := h.service.CreateFullDeal(&req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with creating deal", err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusCreated, "deal successfully created")
}

func (h *DealHandler) Update(c *gin.Context) {
	var req Deal

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "body is invalid", err)
		return
	}

	if err := h.service.Update(&req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with updating deal", err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusCreated, "deal successfully updated")
}

func (h *DealHandler) FindDealByOwnerId(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "you need to provide a query param: userId", errors.New("missing userId"))
		return
	}

	deals, total, err := h.service.FindDealByOwnerId(id)

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
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with getting deals by owner id", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deals successfully found", resp)
}

func (h *DealHandler) FindDealByClientId(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "you need to provide a query param: userId", errors.New("missing userId"))
		return
	}

	deals, total, err := h.service.FindDealByClientId(id)

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
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with getting deals by client id", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "deals successfully found", resp)
}
