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

	pipelines, total, err := h.pipeService.FindAll(page, limit, scope)

	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	data := &PipelinesResponse{
		Pipelines: pipelines,
		Total: total,
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Pipelines successfully found: %d", total), data)
}

func (h *PipelineHandler) Update(c *gin.Context) {
	var req UpdatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	var pipelineId string
	if err := utils.GetParam(c, "pipelineId", &pipelineId); err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	pipeline, err := h.pipeService.Update(pipelineId, &req, scope)
	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "pipeline successfully updated", pipeline)
}

func (h *PipelineHandler) AddStages(c *gin.Context) {
	var req AddPipelineStages
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	var pipelineId string
	if err := utils.GetParam(c, "pipelineId", &pipelineId); err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	pipeline, err := h.pipeService.AddStages(pipelineId, &req, scope)
	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "pipeline successfully updated", pipeline)
}

func (h *PipelineHandler) UpdateStage(c *gin.Context) {
	var stageId string
	if err := utils.GetStringRequiredQuery(c, "stageId", &stageId); err != nil {
		return
	}

	var pipelineId string
	if err := utils.GetParam(c, "pipelineId", &pipelineId); err != nil {
		return
	}

	var req UpdateStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	_, err = h.pipeService.UpdateStage(&req, pipelineId, stageId, scope)

	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "stage successfully updated")
}

func (h *PipelineHandler) DeleteStage(c *gin.Context) {
	var stageId string
	if err := utils.GetStringRequiredQuery(c, "stageId", &stageId); err != nil {
		return
	}

	var pipelineId string
	if err := utils.GetParam(c, "pipelineId", &pipelineId); err != nil {
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	pipeline, err := h.pipeService.DeleteStage(pipelineId, stageId, scope)
	if err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "stage successfully deleted", pipeline)
}

func (h *PipelineHandler) Delete(c *gin.Context) {
	var pipelineId string
	if err := utils.GetParam(c, "pipelineId", &pipelineId); err != nil {
		return
	}

	scope, err := utils.GetAuthScope(c)
	if err != nil {
		return
	}

	if err := h.pipeService.Delete(pipelineId, scope); err != nil {
		utils.CallErrorResponse(c, err)
		return
	}

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "pipeline successfully deleted")
}