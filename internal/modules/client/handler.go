package client

import (
	"avto-crm-api/internal/utils"
	"errors"
	"net/http"

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
	id := c.Query("id")

	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missed client id", errors.New("missed client id"))
		return 
	}

	client, err := h.clientService.FindById(id)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "error with getting client by id", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "successfully getting client", client)
}