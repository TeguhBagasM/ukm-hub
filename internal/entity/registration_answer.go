package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegistrationAnswer struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RegistrationID uuid.UUID `gorm:"type:uuid;index;not null" json:"registration_id"`
	FormFieldID    uuid.UUID `gorm:"type:uuid;index;not null" json:"form_field_id"`
	Value          string    `gorm:"type:text" json:"value"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (a *RegistrationAnswer) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
