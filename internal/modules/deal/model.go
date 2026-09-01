package deal

import (
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/pipeline"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deal struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"" json:"name"`
	Status string `gorm:"default:'active'" json:"status"`

	Pipeline   *pipeline.Pipeline `json:"pipeline"`

	CurrentStage *uuid.UUID `json:"currentStage"`

	Car *car.Car `json:"car"`

	OwnerID  uuid.UUID `json:"ownerId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
