package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;index;not null" json:"organization_id"`
	DivisionID     *uuid.UUID     `gorm:"type:uuid;index" json:"division_id"`
	RegistrationID *uuid.UUID     `gorm:"type:uuid;uniqueIndex" json:"registration_id"`
	Name           string         `gorm:"type:varchar(150);not null" json:"name"`
	Email          string         `gorm:"type:varchar(150)" json:"email"`
	Phone          string         `gorm:"type:varchar(30)" json:"phone"`
	StudentID      string         `gorm:"type:varchar(50)" json:"student_id"`
	Status         string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	JoinedAt       time.Time      `json:"joined_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Member) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
