package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FormField struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	FormID      uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_form_fields_form_name,priority:1;not null" json:"form_id"`
	Label       string         `gorm:"type:varchar(150);not null" json:"label"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex:idx_form_fields_form_name,priority:2;not null" json:"name"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"`
	Placeholder string         `gorm:"type:varchar(255)" json:"placeholder"`
	Description string         `gorm:"type:text" json:"description"`
	Required    bool           `gorm:"default:false" json:"required"`
	Options     string         `gorm:"type:text" json:"options"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (f *FormField) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}
