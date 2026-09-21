package deal

import (
	"time"
)

type StageRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateDealRequest struct {
	Name string `json:"name" binding:"required"`
	Term int `json:"term" binding:"required"`
	Total int `json:"total" binding:"required"`

	Car CarRequest `json:"car" binding:"required"`

	Pipeline PipelineRequest `json:"pipeline" binding:"required"`

	Client ClientRequest `json:"client" binding:"required"`
}

type UpdateDealRequest struct {
	Name *string    `json:"name"`
	Status *string `json:"status"`
	DueDate *string `json:"dueDate"`
	CurrentStageId *string `json:"currentStageId"`
	Term *int `json:"term"`
	Total *int `json:"total"`
}

type UpdateDealInterface struct {
	Name *string    `json:"name"`
	Status *string `json:"status"`
	DueDate *time.Time `json:"dueDate"`
	CurrentStageId *string `json:"currentStageId"`
	CurrentStageName *string `json:"currentStageName"`
	Term *int `json:"term"`
	Total *int `json:"total"`
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

type ClientRequest struct {
	Name string `json:"name" binding:"required"`
	Phone *string `json:"phone"`
	Email *string `json:"email"`
}

type DealsWithTotal struct {
	Deals []Deal `json:"deals"`
	Total int `json:"total"`
	Page int `json:"page"`
	Limit int `json:"limit"`
}

type TotalData struct {
	Total int64 `json:"total"`
	Status string `json:"status"`
}

type DealFilters struct {
	ownerId string
	page int
	limit int
	isFull bool
	status string
	search string
	clientId string
}

type OneDealResponse struct {
	Deal Deal `json:"deal"`
}