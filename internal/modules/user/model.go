package user

import (
	"avto-crm-api/internal/modules/deal"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`

	Email string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"" json:"-"`

	Name string `gorm:"not null" json:"name"`
	MiddleName string `gorm:"not null" json:"middleName"`
	LastName string `json:"lastName,omitempty"`

	OwnerDeals []deal.Deal `gorm:"foreignKey:OwnerId" json:"ownerDeals"`
	ClientDeals []deal.Deal `gorm:"foreignKey:ClientId" json:"clientDeals"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}