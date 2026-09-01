package deal

type StageRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateDealRequest struct {
	Name string `json:"name" binding:"required"`

	Car CarRequest `json:"car" binding:"required"`

	Pipeline PipelineRequest `json:"pipeline" binding:"required"`
}

type DealsResponse struct {
	Deals []Deal `json:"deals"`
	Total int64  `json:"total"`
	Error error `json:"error"`
}
type CarRequest struct {
	Model string `json:"model"`
}

type PipelineRequest struct {
	Name string `json:"name" binding:"required"`
	Source string `json:"source" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	Stages []StageRequest `json:"stages" binding:"required"`
}