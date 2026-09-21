package car

import (
	"github.com/google/uuid"
	"time"
)

type Car struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`

	Model string `gorm:"not null" json:"model"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
