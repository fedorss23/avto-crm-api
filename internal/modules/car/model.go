package car

import (
	"time"
	"github.com/google/uuid"
)

type Car struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`

	DealID uuid.UUID `gorm:"" json:"dealId"`
	
	Model string `gorm:"not null" json:"model"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}