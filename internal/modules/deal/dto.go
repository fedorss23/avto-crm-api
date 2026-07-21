package deal

import (
	"github.com/google/uuid"
)

type StageRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateDealRequest struct {
	Name string `json:"name"`

	Stages []StageRequest `json:"stages"`

	Source       string `json:"source" binding:"required"`
	Destination  string `json:"destination" binding:"required"`
	PipelineName string `json:"pipelineName" binding:"required"`

	OwnerID uuid.UUID `json:"ownerId" binding:"required"`
}

type DealsResponse struct {
	Deals []Deal `json:"deals"`
	Total int64  `json:"total"`
	Error error `json:"error"`
}