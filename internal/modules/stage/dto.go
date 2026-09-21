package stage

type UpdateStageRequest struct {
	Name *string `json:"name"`
	Number *string `json:"number"`
	Destination *string `json:"destination"`
}

type UpdateStageResponse struct {
	Stage Stage `json:"stage"`
}