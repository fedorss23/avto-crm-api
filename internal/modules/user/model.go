package user

import (
	"avto-crm-api/internal/modules/company"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`

	Email string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"" json:"-"`

	Name string `gorm:"not null" json:"name"`
	MiddleName *string `json:"middleName"`
	LastName string `gorm:"not null" json:"lastName"`

	Companies []company.Company `gorm:"foreignKey:OwnerID;references:ID;constaint:onDelete:CASCADE" json:"companies"`

	IsLocked bool `gorm:"default:false" json:"isLocked"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}