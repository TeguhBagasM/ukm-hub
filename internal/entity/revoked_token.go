package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RevokedToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	JTI       string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"jti"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *RevokedToken) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
