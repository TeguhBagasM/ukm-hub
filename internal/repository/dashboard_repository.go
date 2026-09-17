package repository

import (
	"time"

	"ukm-hub/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecentRegistration struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	EventTitle  string
	Status      string
	SubmittedAt time.Time
}

type DivisionMemberCount struct {
	DivisionID   *uuid.UUID
	DivisionName string
	Count        int64
}

type RegistrationExportRow struct {
	ID              uuid.UUID
	EventID         uuid.UUID
	EventTitle      string
	Status          string
	RejectionReason string
	SubmittedAt     time.Time
}

type DashboardRepository interface {
	CountMembers(orgID uuid.UUID) (total, active int64, err error)
	CountEvents(orgID uuid.UUID) (total, active int64, err error)
	CountRegistrations(orgID uuid.UUID) (total, pending int64, err error)
	RegistrationStatusCounts(orgID uuid.UUID) (map[string]int64, error)
	MembersByDivision(orgID uuid.UUID) ([]DivisionMemberCount, error)
	RecentRegistrations(orgID uuid.UUID, limit int) ([]RecentRegistration, error)
	ListRegistrationExportRows(orgID uuid.UUID, eventID *uuid.UUID) ([]RegistrationExportRow, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) CountMembers(orgID uuid.UUID) (int64, int64, error) {
	var total, active int64
	if err := r.db.Model(&entity.Member{}).
		Where("organization_id = ?", orgID).
		Count(&total).Error; err != nil {
		return 0, 0, err
	}
	if err := r.db.Model(&entity.Member{}).
		Where("organization_id = ? AND status = ?", orgID, entity.MemberStatusActive).
		Count(&active).Error; err != nil {
		return 0, 0, err
	}
	return total, active, nil
}

func (r *dashboardRepository) CountEvents(orgID uuid.UUID) (int64, int64, error) {
	var total, active int64
	if err := r.db.Model(&entity.Event{}).
		Where("organization_id = ?", orgID).
		Count(&total).Error; err != nil {
		return 0, 0, err
	}
	if err := r.db.Model(&entity.Event{}).
		Where("organization_id = ? AND status = ?", orgID, entity.EventStatusPublished).
		Count(&active).Error; err != nil {
		return 0, 0, err
	}
	return total, active, nil
}

func (r *dashboardRepository) CountRegistrations(orgID uuid.UUID) (int64, int64, error) {
	base := func() *gorm.DB {
		return r.db.Model(&entity.Registration{}).
			Joins("JOIN events ON events.id = registrations.event_id AND events.deleted_at IS NULL").
			Where("events.organization_id = ?", orgID)
	}
	var total, pending int64
	if err := base().Count(&total).Error; err != nil {
		return 0, 0, err
	}
	if err := base().Where("registrations.status = ?", entity.RegistrationStatusPending).Count(&pending).Error; err != nil {
		return 0, 0, err
	}
	return total, pending, nil
}

func (r *dashboardRepository) RegistrationStatusCounts(orgID uuid.UUID) (map[string]int64, error) {
	type statusCount struct {
		Status string
		Count  int64
	}
	var rows []statusCount
	if err := r.db.Model(&entity.Registration{}).
		Select("registrations.status AS status, COUNT(*) AS count").
		Joins("JOIN events ON events.id = registrations.event_id AND events.deleted_at IS NULL").
		Where("events.organization_id = ?", orgID).
		Group("registrations.status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := map[string]int64{
		entity.RegistrationStatusPending:  0,
		entity.RegistrationStatusAccepted: 0,
		entity.RegistrationStatusRejected: 0,
	}
	for _, row := range rows {
		result[row.Status] = row.Count
	}
	return result, nil
}

func (r *dashboardRepository) MembersByDivision(orgID uuid.UUID) ([]DivisionMemberCount, error) {
	var rows []DivisionMemberCount
	if err := r.db.Model(&entity.Member{}).
		Select("members.division_id AS division_id, COALESCE(divisions.name, 'Tanpa Divisi') AS division_name, COUNT(*) AS count").
		Joins("LEFT JOIN divisions ON divisions.id = members.division_id AND divisions.deleted_at IS NULL").
		Where("members.organization_id = ? AND members.deleted_at IS NULL", orgID).
		Group("members.division_id, divisions.name").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *dashboardRepository) RecentRegistrations(orgID uuid.UUID, limit int) ([]RecentRegistration, error) {
	var rows []RecentRegistration
	if err := r.db.Model(&entity.Registration{}).
		Select("registrations.id AS id, registrations.event_id AS event_id, events.name AS event_title, registrations.status AS status, registrations.submitted_at AS submitted_at").
		Joins("JOIN events ON events.id = registrations.event_id AND events.deleted_at IS NULL").
		Where("events.organization_id = ?", orgID).
		Order("registrations.submitted_at DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *dashboardRepository) ListRegistrationExportRows(orgID uuid.UUID, eventID *uuid.UUID) ([]RegistrationExportRow, error) {
	var rows []RegistrationExportRow
	query := r.db.Model(&entity.Registration{}).
		Select("registrations.id AS id, registrations.event_id AS event_id, registrations.status AS status, registrations.rejection_reason AS rejection_reason, registrations.submitted_at AS submitted_at, events.name AS event_title").
		Joins("JOIN events ON events.id = registrations.event_id AND events.deleted_at IS NULL").
		Where("events.organization_id = ?", orgID)
	if eventID != nil {
		query = query.Where("registrations.event_id = ?", *eventID)
	}
	if err := query.Order("registrations.submitted_at ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
