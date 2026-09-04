package client

import (
	"github.com/google/uuid"
	"time"

)

type Client struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	DealID *uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"dealId"`
	OwnerID uuid.UUID `json:"ownerId"`

	Email *string `json:"email"`
	Phone *string `json:"phone"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}