package manager

import (
	"avto-crm-api/internal/modules/deal"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Manager struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`

	Login    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"" json:"-"`

	Name       string  `gorm:"not null" json:"name"`
	MiddleName *string  `json:"middleName"`
	LastName   string `gorm:"not null" json:"lastName"`

	Phone *string `json:"phone"`
	Role  string  `gorm:"default:standart" json:"role"`

	CompanyID uuid.UUID `gorm:"not null" json:"companyId"`
	OwnerID uuid.UUID `gorm:"not null" json:"dealId"`

	Deals []deal.Deal `gorm:"foreignKey:ManagerID;references:ID;constraint:OnDelete:CASCADE" json:"deals"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
