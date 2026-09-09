package user

import (
	"avto-crm-api/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Delete(c *gin.Context) {
	var userId string
	if err := utils.GetStringRequiredQuery(c, "userId", &userId); err != nil {
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	err := h.userService.Delete(userId, ownerId)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusCreated, "user successfully deleted")
}

func (h *UserHandler) FindList(c *gin.Context) {
	var page int
	if err := utils.GetNumberQuery(c, "page", &page, 1, ""); err != nil {
		return
	}

	var limit int
	if err := utils.GetNumberQuery(c, "limit", &limit, 10, ""); err != nil {
		return
	}

	users, total, err := h.userService.FindList(page, limit)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("successfully found: %d", total), users)
}