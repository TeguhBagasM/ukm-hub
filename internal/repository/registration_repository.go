package repository

import (
	"handler/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegistrationRepository interface {
	Create(reg *entity.Registration) error
	CreateWithAnswers(reg *entity.Registration, answers []entity.RegistrationAnswer) error
	CountByEventID(eventID uuid.UUID) (int, error)
	ExistsAnswer(eventID, formFieldID uuid.UUID, value string) (bool, error)
	ListByEvent(eventID uuid.UUID, status, search string, page, perPage int) ([]entity.Registration, int64, error)
	FindByID(id uuid.UUID) (*entity.Registration, error)
	FindAnswersByRegistrationIDs(ids []uuid.UUID) ([]entity.RegistrationAnswer, error)
	Update(reg *entity.Registration) error
}

type registrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{db: db}
}

func (r *registrationRepository) Create(reg *entity.Registration) error {
	return r.db.Create(reg).Error
}

func (r *registrationRepository) CreateWithAnswers(reg *entity.Registration, answers []entity.RegistrationAnswer) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reg).Error; err != nil {
			return err
		}
		for i := range answers {
			answers[i].RegistrationID = reg.ID
			if err := tx.Create(&answers[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *registrationRepository) CountByEventID(eventID uuid.UUID) (int, error) {
	var count int64
	if err := r.db.Model(&entity.Registration{}).
		Where("event_id = ?", eventID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *registrationRepository) ExistsAnswer(eventID, formFieldID uuid.UUID, value string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.RegistrationAnswer{}).
		Joins("JOIN registrations ON registrations.id = registration_answers.registration_id").
		Where("registrations.event_id = ? AND registration_answers.form_field_id = ? AND LOWER(registration_answers.value) = LOWER(?)", eventID, formFieldID, value).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *registrationRepository) ListByEvent(eventID uuid.UUID, status, search string, page, perPage int) ([]entity.Registration, int64, error) {
	query := r.db.Model(&entity.Registration{}).Where("event_id = ?", eventID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		sub := r.db.Model(&entity.RegistrationAnswer{}).
			Select("registration_id").
			Where("value ILIKE ?", "%"+search+"%")
		query = query.Where("id IN (?)", sub)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var regs []entity.Registration
	if err := query.
		Order("submitted_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&regs).Error; err != nil {
		return nil, 0, err
	}
	return regs, total, nil
}

func (r *registrationRepository) FindByID(id uuid.UUID) (*entity.Registration, error) {
	var reg entity.Registration
	if err := r.db.First(&reg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *registrationRepository) FindAnswersByRegistrationIDs(ids []uuid.UUID) ([]entity.RegistrationAnswer, error) {
	if len(ids) == 0 {
		return []entity.RegistrationAnswer{}, nil
	}
	var answers []entity.RegistrationAnswer
	if err := r.db.Where("registration_id IN ?", ids).
		Order("created_at ASC").
		Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *registrationRepository) Update(reg *entity.Registration) error {
	return r.db.Save(reg).Error
}
