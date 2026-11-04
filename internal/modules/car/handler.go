package car

import (
	"avto-crm-api/internal/utils"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	carService *CarService
}

func NewCarHandler(carService *CarService) *CarHandler {
	return &CarHandler{
		carService: carService,
	}
}

func (h *CarHandler) FindAll(c *gin.Context) {
	var page int
	if err := utils.GetNumberQuery(c, "page", &page, 1, utils.PageErrorCode); err != nil {
		return
	}

	var limit int
	if err := utils.GetNumberQuery(c, "limit", &limit, 10, utils.LimitErrorCode); err != nil {
		return
	}

	scope, err := utils.GetAuthScopeWithManager(c)
	if err != nil {
		return
	}

	cars, total, err := h.carService.FindList(page, limit, scope)

	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Получено карточек: %d", total), &CarListWithTotal{
		Total: total,
		Cars: cars,
	})
}

func (h *CarHandler) Update(c *gin.Context) {
	var req UpdateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	var carId string
	if err := utils.GetParam(c, "carId", &carId); err != nil {
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	car, err := h.carService.Update(&req, carId, scope)
	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Машина успешно обновлена", UpdateCarResponse{
		Car: *car,
	})
}

func (h *CarHandler) Delete(c *gin.Context) {
	var carId string
	if err := utils.GetParam(c, "carId", &carId); err != nil {
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	if err := h.carService.Delete(carId, scope); err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "car successfully deleted")
}