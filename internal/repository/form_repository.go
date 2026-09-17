package repository

import (
	"ukm-hub/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FormRepository interface {
	Create(form *entity.Form) error
	FindByID(id uuid.UUID) (*entity.Form, error)
	FindByEvent(eventID uuid.UUID) (*entity.Form, error)
	Update(form *entity.Form) error
	Delete(id uuid.UUID) error
}

type formRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) FormRepository {
	return &formRepository{db: db}
}

func (r *formRepository) Create(form *entity.Form) error {
	return r.db.Create(form).Error
}

func (r *formRepository) FindByID(id uuid.UUID) (*entity.Form, error) {
	var form entity.Form
	if err := r.db.First(&form, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *formRepository) FindByEvent(eventID uuid.UUID) (*entity.Form, error) {
	var form entity.Form
	if err := r.db.Where("event_id = ?", eventID).First(&form).Error; err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *formRepository) Update(form *entity.Form) error {
	return r.db.Save(form).Error
}

func (r *formRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Form{}, "id = ?", id).Error
}

type FormFieldRepository interface {
	Create(field *entity.FormField) error
	FindByID(id uuid.UUID) (*entity.FormField, error)
	FindByForm(formID uuid.UUID) ([]entity.FormField, error)
	MaxSortOrder(formID uuid.UUID) (int, error)
	Update(field *entity.FormField) error
	Delete(id uuid.UUID) error
}

type formFieldRepository struct {
	db *gorm.DB
}

func NewFormFieldRepository(db *gorm.DB) FormFieldRepository {
	return &formFieldRepository{db: db}
}

func (r *formFieldRepository) Create(field *entity.FormField) error {
	return r.db.Create(field).Error
}

func (r *formFieldRepository) FindByID(id uuid.UUID) (*entity.FormField, error) {
	var field entity.FormField
	if err := r.db.First(&field, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

func (r *formFieldRepository) FindByForm(formID uuid.UUID) ([]entity.FormField, error) {
	var fields []entity.FormField
	if err := r.db.Where("form_id = ?", formID).
		Order("sort_order ASC, created_at ASC").
		Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

func (r *formFieldRepository) MaxSortOrder(formID uuid.UUID) (int, error) {
	var max int
	if err := r.db.Model(&entity.FormField{}).
		Where("form_id = ?", formID).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error; err != nil {
		return 0, err
	}
	return max, nil
}

func (r *formFieldRepository) Update(field *entity.FormField) error {
	return r.db.Save(field).Error
}

func (r *formFieldRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.FormField{}, "id = ?", id).Error
}
