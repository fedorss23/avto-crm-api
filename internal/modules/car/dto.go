package car	

type CreateCarRequest struct {
	Model string `json:"model" binding:"required"`
}

type CarListWithTotal struct {
	Total int64
	Cars []Car `json:"cars"`
}