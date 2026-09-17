package repository

import (
	"ukm-hub/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(org *entity.Organization) error
	FindByID(id uuid.UUID) (*entity.Organization, error)
	FindBySlug(slug string) (*entity.Organization, error)
	FindAll() ([]entity.Organization, error)
	FindByAdminUser(userID uuid.UUID) ([]entity.Organization, error)
	IsAdminOf(userID, orgID uuid.UUID) (bool, error)
	Update(org *entity.Organization) error
	Delete(id uuid.UUID) error
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(org *entity.Organization) error {
	return r.db.Create(org).Error
}

func (r *organizationRepository) FindByID(id uuid.UUID) (*entity.Organization, error) {
	var org entity.Organization
	if err := r.db.First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) FindBySlug(slug string) (*entity.Organization, error) {
	var org entity.Organization
	if err := r.db.Where("slug = ?", slug).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) FindAll() ([]entity.Organization, error) {
	var orgs []entity.Organization
	if err := r.db.Order("created_at ASC").Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *organizationRepository) FindByAdminUser(userID uuid.UUID) ([]entity.Organization, error) {
	var orgs []entity.Organization
	if err := r.db.
		Joins("JOIN organization_admins oa ON oa.organization_id = organizations.id").
		Where("oa.user_id = ?", userID).
		Order("organizations.created_at ASC").
		Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *organizationRepository) IsAdminOf(userID, orgID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.OrganizationAdmin{}).
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *organizationRepository) Update(org *entity.Organization) error {
	return r.db.Save(org).Error
}

func (r *organizationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Organization{}, "id = ?", id).Error
}
