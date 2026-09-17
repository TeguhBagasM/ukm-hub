package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(150);uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Logo        string         `gorm:"type:text" json:"logo"`
	Email       string         `gorm:"type:varchar(150)" json:"email"`
	Phone       string         `gorm:"type:varchar(30)" json:"phone"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
