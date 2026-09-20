package client

import (
	"avto-crm-api/internal/utils"
	"fmt"
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
	var clientId string
	if err := utils.GetParam(c, "clientId", &clientId); err != nil {
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	client, err := h.clientService.FindById(clientId, ownerId)

	if err != nil {
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), err.Error(), err, CodeByError(err))
	}

	utils.SuccessResponse(c, http.StatusOK, "successfully getting client", client)
}

func (h *ClientHandler) FindListByOwnerId(c *gin.Context) {
	var page int
	if err := utils.GetNumberQuery(c, "page", &page, 1, pageErrorCode); err != nil {
		return
	}

	var limit int
	if err := utils.GetNumberQuery(c, "limit", &limit, 1, limitErrorCode); err != nil {
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	clients, total, err := h.clientService.FindListByOwnerId(ownerId, page, limit)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err, CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Clients successfully found: %d", total), &ClientWithTotal{
		Clients: clients,
		Total: total,
	})
}