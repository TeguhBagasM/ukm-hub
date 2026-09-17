package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Form struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	EventID          uuid.UUID      `gorm:"type:uuid;index;not null" json:"event_id"`
	Title            string         `gorm:"type:varchar(150)" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	SubmitWebhookURL string         `gorm:"type:text" json:"submit_webhook_url"`
	IsPublished      bool           `gorm:"default:false" json:"is_published"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (f *Form) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}
