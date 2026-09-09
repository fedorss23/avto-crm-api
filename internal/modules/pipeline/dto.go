package pipeline

type PipelinesResponse struct {
	Pipelines []Pipeline `json:"pipelines"`
	Total int64 `json:"total"`
}