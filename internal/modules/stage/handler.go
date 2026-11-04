package stage

import (
)

type StageHandler struct {
	stageService *StageService
}

func NewStageHandler(stageService *StageService) *StageHandler {
	return &StageHandler{
		stageService: stageService,
	}
}