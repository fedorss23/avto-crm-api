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

func (h *CarHandler) Create(c *gin.Context) {
	var req *Car

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs, utils.ValidationErrorCode)
		return
	}

	if err := h.carService.carRepo.Create(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err, CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Машина успешно создана", req)
}

func (h *CarHandler) FindAll(c *gin.Context) {
	var page int
	if err := utils.GetNumberQuery(c, "page", &page, 1, pageErrorCode); err != nil {
		return
	}

	var limit int
	if err := utils.GetNumberQuery(c, "limit", &limit, 10, limitErrorCode); err != nil {
		return
	}

	cars, total, err := h.carService.FindList(page, limit)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err, CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Получено карточек: %d", total), &CarListWithTotal{
		Total: total,
		Cars: cars,
	})
}

func (h *CarHandler) Update(c *gin.Context) {
	var req *Car

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs, utils.ValidationErrorCode)
		return
	}

	if err := h.carService.Update(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err, CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Машина успешно обновлена", req)
}