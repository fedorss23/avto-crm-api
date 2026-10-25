package manager

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Manager struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`

	Login    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"" json:"-"`

	Name       string  `gorm:"not null" json:"name"`
	MiddleName string  `gorm:"not null" json:"middleName"`
	LastName   *string `json:"lastName"`

	Phone *string `json:"phone"`
	Role  string  `gorm:"default:standart" json:"role"`

	CompanyID uuid.UUID `gorm:"not null" json:"companyId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
