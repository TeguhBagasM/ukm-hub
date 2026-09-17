package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Division struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;index;not null" json:"organization_id"`
	Name           string         `gorm:"type:varchar(150);not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Status         string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Division) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return
}
