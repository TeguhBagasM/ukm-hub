package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationAdmin struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_org_admin_user_org,priority:1;not null" json:"user_id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_org_admin_user_org,priority:2;not null" json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (oa *OrganizationAdmin) BeforeCreate(tx *gorm.DB) (err error) {
	if oa.ID == uuid.Nil {
		oa.ID = uuid.New()
	}
	return
}
