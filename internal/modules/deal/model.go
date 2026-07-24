package deal

import (
	"avto-crm-api/internal/modules/pipeline"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deal struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"" json:"name"`

	PipelineID *uuid.UUID `json:"pipelineId"`
	Pipeline   *pipeline.Pipeline

	CurrentStage *uuid.UUID `json:"currentStage"`

	CarID *uuid.UUID `json:"carId"`

	OwnerID  uuid.UUID `json:"ownerId"`
	ClientID *uuid.UUID `json:"clientId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
