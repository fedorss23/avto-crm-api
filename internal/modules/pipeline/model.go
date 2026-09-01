package pipeline

import (
	"avto-crm-api/internal/modules/stage"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pipeline struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"" json:"name"`

	DealID *uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"dealId"`

	Source string `gorm:"not null" json:"source"`
	Destination string `gorm:"not null" json:"destination"`

	Stages []stage.Stage `gorm:"foreignKey:PipelineID" json:"stages"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
