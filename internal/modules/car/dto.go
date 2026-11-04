package car	

type CreateCarRequest struct {
	Model string `json:"model" binding:"required"`
}

type CarListWithTotal struct {
	Total int64
	Cars []Car `json:"cars"`
}

type UpdateCarRequest struct {
	Model *string `json:"model"`
}

type UpdateCarResponse struct {
	Car Car `json:"car"`
}