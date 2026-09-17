package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Registration struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	EventID         uuid.UUID `gorm:"type:uuid;index;not null" json:"event_id"`
	Status          string    `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	RejectionReason string    `gorm:"type:text" json:"rejection_reason"`
	SubmittedAt     time.Time `json:"submitted_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (r *Registration) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
