package repository

import (
	"handler/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *entity.Event) error
	FindByID(id uuid.UUID) (*entity.Event, error)
	FindBySlug(slug string) (*entity.Event, error)
	FindByOrganization(orgID uuid.UUID) ([]entity.Event, error)
	Update(event *entity.Event) error
	Delete(id uuid.UUID) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *entity.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindByID(id uuid.UUID) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindBySlug(slug string) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("slug = ?", slug).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindByOrganization(orgID uuid.UUID) ([]entity.Event, error) {
	var events []entity.Event
	if err := r.db.Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) Update(event *entity.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Event{}, "id = ?", id).Error
}
