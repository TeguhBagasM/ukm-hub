package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"organization_id"`
	Name              string         `gorm:"type:varchar(150);not null" json:"name"`
	Slug              string         `gorm:"type:varchar(150);uniqueIndex;not null" json:"slug"`
	Description       string         `gorm:"type:text" json:"description"`
	Location          string         `gorm:"type:varchar(150)" json:"location"`
	StartDate         *time.Time     `json:"start_date"`
	EndDate           *time.Time     `json:"end_date"`
	RegistrationStart *time.Time     `json:"registration_start"`
	RegistrationEnd   *time.Time     `json:"registration_end"`
	Quota             int            `gorm:"default:0" json:"quota"`
	Status            string         `gorm:"type:varchar(20);default:'DRAFT'" json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}
