package car

import (
	"avto-crm-api/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
		utils.ValidationErrorResponse(c, errs)
		return
	}

	if err := h.carService.carRepo.Create(req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка на сервере", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Машина успешно создана", req)
}

func (h *CarHandler) FindAll(c *gin.Context) {
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

	cars, total, err := h.carService.FindList(intPage, intLimit)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении машин", err)
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
		utils.ValidationErrorResponse(c, errs)
		return
	}

	if err := h.carService.Update(req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении машины: ошибка со стороны сервера", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Машина успешно обновлена", req)
}