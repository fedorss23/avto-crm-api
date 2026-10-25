package company

import (
	"avto-crm-api/internal/modules/manager"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"not null" json:"name"`

	TypeCompany string `gorm:"not null" json:"type"`

	OwnerID uuid.UUID `gorm:"not null" json:"ownerId"`
	
	Managers []manager.Manager `gorm:"foreignKey:CompanyID" json:"managers"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
