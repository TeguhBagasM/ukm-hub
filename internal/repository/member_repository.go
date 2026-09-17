package repository

import (
	"handler/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(member *entity.Member) error
	FindByID(id uuid.UUID) (*entity.Member, error)
	FindByRegistrationID(regID uuid.UUID) (*entity.Member, error)
	FindByRegistrationIDUnscoped(regID uuid.UUID) (*entity.Member, error)
	ListByOrganization(orgID uuid.UUID, search, divisionID string, page, perPage int) ([]entity.Member, int64, error)
	Update(member *entity.Member) error
	Restore(member *entity.Member) error
	Delete(id uuid.UUID) error
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(member *entity.Member) error {
	return r.db.Create(member).Error
}

func (r *memberRepository) FindByID(id uuid.UUID) (*entity.Member, error) {
	var member entity.Member
	if err := r.db.First(&member, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) FindByRegistrationID(regID uuid.UUID) (*entity.Member, error) {
	var member entity.Member
	if err := r.db.Where("registration_id = ?", regID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) FindByRegistrationIDUnscoped(regID uuid.UUID) (*entity.Member, error) {
	var member entity.Member
	if err := r.db.Unscoped().Where("registration_id = ?", regID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) ListByOrganization(orgID uuid.UUID, search, divisionID string, page, perPage int) ([]entity.Member, int64, error) {
	query := r.db.Model(&entity.Member{}).Where("organization_id = ?", orgID)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR student_id ILIKE ?", like, like, like)
	}
	if divisionID != "" {
		query = query.Where("division_id = ?", divisionID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var members []entity.Member
	if err := query.
		Order("joined_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&members).Error; err != nil {
		return nil, 0, err
	}
	return members, total, nil
}

func (r *memberRepository) Update(member *entity.Member) error {
	return r.db.Save(member).Error
}

func (r *memberRepository) Restore(member *entity.Member) error {
	member.DeletedAt = gorm.DeletedAt{}
	return r.db.Unscoped().Save(member).Error
}

func (r *memberRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Member{}, "id = ?", id).Error
}
