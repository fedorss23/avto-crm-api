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

	Source      string `gorm:"not null" json:"source"`
	Destination string `gorm:"not null" json:"destination"`

	DealID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"dealId"`

	Stages []stage.Stage `gorm:"foreignKey:PipelineID;references:ID;constraint:OnDelete:CASCADE" json:"stages"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
