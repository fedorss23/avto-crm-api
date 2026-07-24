package stage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stage struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	PipelineID uuid.UUID `json:"pipelineId"`

	Description *string `json:"description"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
