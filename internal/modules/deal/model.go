package deal

import (
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/pipeline"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deal struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"not null" json:"name"`
	Status string `gorm:"default:'active'" json:"status"`

	Pipeline   *pipeline.Pipeline `gorm:"foreignKey:DealID;references:ID;constraint:OnDelete:CASCADE" json:"pipeline"`

	CurrentStageId *uuid.UUID `json:"currentStage"`
	CurrentStageName *string `json:"currentStageName"`

	Car *car.Car `gorm:"foreignKey:DealID;references:ID;constraint:OnDelete:CASCADE" json:"car"`

	ManagerID  uuid.UUID `gorm:"not null" json:"managerId"`
	CompanyID uuid.UUID `gorm:"not null" json:"companyId"`

	Client *client.Client `gorm:"foreignKey:DealID;references:ID;constraint:OnDelete:CASCADE" json:"client"`

	DueDate time.Time `gorm:"not null" json:"dueDate"`
	Term int `gorm:"not null" json:"term"`
	Total int `gorm:"not null" json:"total"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

