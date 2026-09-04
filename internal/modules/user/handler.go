package user

import (
	"avto-crm-api/internal/utils"
	"errors"
	"fmt"
	"net/http"

	"strconv"

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
	userId := c.Param("userId")

	if userId == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "error with param user id", errors.New("missing param userId"))
		return
	}

	err := h.userService.Delete(userId)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error with deleting user", err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusCreated, "user successfully deleted")
}

func (h *UserHandler) FindList(c *gin.Context) {
	var pageInt int
	
	page, exists := c.GetQuery("page")
	if !exists {
		pageInt = 1
	} else {
		a, err := strconv.Atoi(page)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Error with found users: query-param page must be number", errors.New("Page must be number"))
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
			utils.ErrorResponse(c, http.StatusBadRequest, "Error with found users: query-param limit must be number", errors.New("Limit must be number"))
			return
		}
		limitInt = a
	}

	users, total, err := h.userService.FindList(pageInt, limitInt)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Error with found users", err)
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("successfully found: %d", total), users)
}