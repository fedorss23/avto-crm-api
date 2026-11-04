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

type ClientScope struct {
	CompanyId string
	ManagerId *string
}

func NewClientHandler(clientService *ClientService) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
	}
}

func (h *ClientHandler) FindList(c *gin.Context) {
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

	clients, total, err := h.clientService.FindList(page, limit, scope)

	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Clients successfully found: %d", total), &ClientWithTotal{
		Clients: clients,
		Total: total,
	})
}

func (h *ClientHandler) Update(c *gin.Context) {
	var req UpdateClientRequest
	if err := c.ShouldBindJSON(req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	var clientId string
	if err := utils.GetParam(c, "clientId", &clientId); err != nil {
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	client, err := h.clientService.Update(&req, clientId, scope)
	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "client successfully updated", OneClientResponse{
		Client: *client,
	})
}

func (h *ClientHandler) Delete(c *gin.Context) {
	var clientId string
	if err := utils.GetParam(c, "clientId", &clientId); err != nil {
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	if err := h.clientService.Delete(clientId, scope); err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "client successfully deleted")
}