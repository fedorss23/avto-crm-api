package pipeline

import "avto-crm-api/internal/modules/stage"

type PipelinesResponse struct {
	Pipelines []Pipeline `json:"pipelines"`
	Total int64 `json:"total"`
}

type UpdatePipelineRequest struct {
	Name *string `json:"name"`
	Source *string `json:"source"`
	Destination *string `json:"destination"`
}

type AddPipelineStages struct {
	Stages []AddStage `json:"stages"`
}

type AddStage struct {
	Name string `json:"name"`
	Description *string `json:"description"`
}

type UpdatePipelineResponse struct {
	Pipeline Pipeline `json:"pipeline"`
}

type UpdateStageRequest struct {
	Name *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateStageResponse struct {
	Stage stage.Stage `json:"stage"`
}