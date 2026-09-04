package pipeline

import (
	"avto-crm-api/internal/utils"
	"fmt"
	"net/http"

	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PipelineHandler struct {
	pipeService *PipelineService
}

func NewPipelineHandler(pipeService *PipelineService) *PipelineHandler {
	return &PipelineHandler{
		pipeService: pipeService,
	}
}

func (h *PipelineHandler) FindList(c *gin.Context) {
	var pageInt int
	
	page, exists := c.GetQuery("page")
	if !exists {
		pageInt = 1
	} else {
		a, err := strconv.Atoi(page)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Error with found deal: query-param page must be number", errors.New("Page must be number"))
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
			utils.ErrorResponse(c, http.StatusBadRequest, "Error with found deal: query-param limit must be number", errors.New("Limit must be number"))
			return
		}
		limitInt = a
	}

	pipelines, total, err := h.pipeService.FindAll(pageInt, limitInt)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Error with getting pipelines", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Pipelines successfully found: %d", total), pipelines)
}