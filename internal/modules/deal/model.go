package deal

import (
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/pipeline"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deal struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"not null" json:"name"`
	Status string `gorm:"default:'active'" json:"status"`

	Pipeline   *pipeline.Pipeline `gorm:"not null" json:"pipeline"`
	PipelineId *uuid.UUID `json:"pipelineId"`

	CurrentStageId *uuid.UUID `json:"currentStage"`
	CurrentStageName *string `json:"currentStageName"`

	Car *car.Car `gorm:"not null" json:"car"`
	CarId *uuid.UUID `json:"carId"`

	OwnerID  uuid.UUID `gorm:"not null" json:"ownerId"`

	Client *client.Client `gorm:"not null" json:"client"`
	ClientId *uuid.UUID `json:"clientId"`

	DueDate time.Time `gorm:"not null" json:"dueDate"`
	Term int `gorm:"not null" json:"term"`
	Total int `gorm:"not null" json:"total"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
