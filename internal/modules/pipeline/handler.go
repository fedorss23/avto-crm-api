package pipeline

import (
	"avto-crm-api/internal/utils"
	"fmt"
	"net/http"

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
	var page int
	if err := utils.GetNumberQuery(c, "page", &page, 1, pageErrorCode); err != nil {
		return
	}

	var limit int
	if err := utils.GetNumberQuery(c, "limit", &limit, 10, limitErrorCode); err != nil {
		return
	}

	pipelines, total, err := h.pipeService.FindAll(page, limit)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err, CodeByError(err))
		return
	}

	data := &PipelinesResponse{
		Pipelines: pipelines,
		Total: total,
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Pipelines successfully found: %d", total), data)
}