package stage

import (
	"avto-crm-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StageHandler struct {
	stageService *StageService
}

func NewStageHandler(stageService *StageService) *StageHandler {
	return &StageHandler{
		stageService: stageService,
	}
}

func (s *StageHandler) Update(c *gin.Context) {
	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	var stageId string
	if err := utils.GetParam(c, "stageId", &stageId); err != nil {
		return
	}

	var req *UpdateStageRequest

	if err := c.ShouldBindJSON(req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	stage, err := s.stageService.Update(req, stageId, ownerId)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "stage successfully updated", &UpdateStageResponse{
		Stage: *stage,
	})
}
