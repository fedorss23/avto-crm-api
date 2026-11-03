package pipeline

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