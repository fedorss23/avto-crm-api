package client

import (
	"avto-crm-api/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientService *ClientService
}

func NewClientHandler(clientService *ClientService) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
	}
}

func (h *ClientHandler) FindById(c *gin.Context) {
	id := c.Param("clientId")

	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missed client id", errors.New("missed client id"))
		return 
	}

	ownerId := c.GetString("userId")

	client, err := h.clientService.FindById(id, ownerId)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with getting client by id", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "successfully getting client", client)
}

func (h *ClientHandler) FindListByOwnerId(c *gin.Context) {
	var pageInt int

	page, ok := c.GetQuery("page")
	if !ok {
		pageInt = 1
	} else {
		a, err := strconv.Atoi(page)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "page must be a number", errors.New("page must be a number"))
			return
		}
		pageInt = a
	}

	var limitInt int

	limit, ok := c.GetQuery("limit")
	if !ok {
		limitInt = 10
	} else {
		a, err := strconv.Atoi(limit)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "limit must be a number", errors.New("limit must be a number"))
			return
		}
		limitInt = a
	}

	ownerId := c.GetString("userId")
	if ownerId == "" {
		utils.ErrorResponse(c, 401, "you aren't authorized", errors.New("you need to authorize"))
		return
	}

	clients, total, err := h.clientService.FindListByOwnerId(ownerId, pageInt, limitInt)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "failed with getting clients", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Clients successfully found: %d", total), clients)
}