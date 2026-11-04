package car

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Car struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`

	Model string `gorm:"not null" json:"model"`

	DealID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"dealId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
