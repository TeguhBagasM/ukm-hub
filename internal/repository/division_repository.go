package repository

import (
	"ukm-hub/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DivisionRepository interface {
	Create(div *entity.Division) error
	FindByID(id uuid.UUID) (*entity.Division, error)
	FindByOrganization(orgID uuid.UUID) ([]entity.Division, error)
	Update(div *entity.Division) error
	Delete(id uuid.UUID) error
}

type divisionRepository struct {
	db *gorm.DB
}

func NewDivisionRepository(db *gorm.DB) DivisionRepository {
	return &divisionRepository{db: db}
}

func (r *divisionRepository) Create(div *entity.Division) error {
	return r.db.Create(div).Error
}

func (r *divisionRepository) FindByID(id uuid.UUID) (*entity.Division, error) {
	var div entity.Division
	if err := r.db.First(&div, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &div, nil
}

func (r *divisionRepository) FindByOrganization(orgID uuid.UUID) ([]entity.Division, error) {
	var divs []entity.Division
	if err := r.db.Where("organization_id = ?", orgID).
		Order("created_at ASC").
		Find(&divs).Error; err != nil {
		return nil, err
	}
	return divs, nil
}

func (r *divisionRepository) Update(div *entity.Division) error {
	return r.db.Save(div).Error
}

func (r *divisionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Division{}, "id = ?", id).Error
}
