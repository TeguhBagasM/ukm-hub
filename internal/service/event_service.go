package service

import (
	"net/http"
	"time"

	"handler/internal/dto"
	"handler/internal/entity"
	"handler/internal/repository"
	"handler/internal/utils"

	"github.com/google/uuid"
)

type EventService interface {
	ListByOrganization(ctx dto.AuthContext, orgIDStr string) ([]dto.EventResponse, error)
	Get(ctx dto.AuthContext, idStr string) (*dto.EventResponse, error)
	Create(ctx dto.AuthContext, orgIDStr string, req dto.CreateEventRequest) (*dto.EventResponse, error)
	Update(ctx dto.AuthContext, idStr string, req dto.UpdateEventRequest) (*dto.EventResponse, error)
	Delete(ctx dto.AuthContext, idStr string) error
	Publish(ctx dto.AuthContext, idStr string) (*dto.EventResponse, error)
	Close(ctx dto.AuthContext, idStr string) (*dto.EventResponse, error)
	Archive(ctx dto.AuthContext, idStr string) (*dto.EventResponse, error)
}

type eventService struct {
	eventRepo repository.EventRepository
	regRepo   repository.RegistrationRepository
	access    AccessService
}

func NewEventService(eventRepo repository.EventRepository, regRepo repository.RegistrationRepository, access AccessService) EventService {
	return &eventService{eventRepo: eventRepo, regRepo: regRepo, access: access}
}

func (s *eventService) ListByOrganization(ctx dto.AuthContext, orgIDStr string) ([]dto.EventResponse, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if err := s.access.MustAccessOrganization(ctx, orgID); err != nil {
		return nil, err
	}
	events, err := s.eventRepo.FindByOrganization(orgID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.EventResponse, 0, len(events))
	for i := range events {
		res = append(res, *s.toResponse(&events[i], 0))
	}
	return res, nil
}

func (s *eventService) Get(ctx dto.AuthContext, idStr string) (*dto.EventResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid event id")
	}
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return nil, err
	}
	count, err := s.regRepo.CountByEventID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(event, count), nil
}

func (s *eventService) Create(ctx dto.AuthContext, orgIDStr string, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if err := s.access.MustAccessOrganization(ctx, orgID); err != nil {
		return nil, err
	}
	if existing, _ := s.eventRepo.FindBySlug(req.Slug); existing != nil {
		return nil, utils.NewAppError(http.StatusConflict, "slug already used")
	}
	if err := validateEventDates(req.RegistrationStart, req.RegistrationEnd, req.StartDate, req.EndDate); err != nil {
		return nil, err
	}
	status := req.Status
	if status == "" {
		status = entity.EventStatusDraft
	}
	event := entity.Event{
		OrganizationID:    orgID,
		Name:              req.Name,
		Slug:              req.Slug,
		Description:       req.Description,
		Location:          req.Location,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		RegistrationStart: req.RegistrationStart,
		RegistrationEnd:   req.RegistrationEnd,
		Quota:             req.Quota,
		Status:            status,
	}
	if err := s.eventRepo.Create(&event); err != nil {
		return nil, err
	}
	return s.toResponse(&event, 0), nil
}

func (s *eventService) Update(ctx dto.AuthContext, idStr string, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid event id")
	}
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return nil, err
	}
	if req.Name != "" {
		event.Name = req.Name
	}
	if req.Slug != "" {
		if other, _ := s.eventRepo.FindBySlug(req.Slug); other != nil && other.ID != id {
			return nil, utils.NewAppError(http.StatusConflict, "slug already used")
		}
		event.Slug = req.Slug
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Location != "" {
		event.Location = req.Location
	}
	if req.StartDate != nil {
		event.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		event.EndDate = req.EndDate
	}
	if req.RegistrationStart != nil {
		event.RegistrationStart = req.RegistrationStart
	}
	if req.RegistrationEnd != nil {
		event.RegistrationEnd = req.RegistrationEnd
	}
	if req.Quota != nil {
		event.Quota = *req.Quota
	}
	if req.Status != "" {
		event.Status = req.Status
	}
	if err := validateEventDates(event.RegistrationStart, event.RegistrationEnd, event.StartDate, event.EndDate); err != nil {
		return nil, err
	}
	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}
	count, err := s.regRepo.CountByEventID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(event, count), nil
}

func (s *eventService) Delete(ctx dto.AuthContext, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.NewAppError(http.StatusBadRequest, "invalid event id")
	}
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return err
	}
	return s.eventRepo.Delete(id)
}

func (s *eventService) Publish(ctx dto.AuthContext, idStr string) (*dto.EventResponse, error) {
	event, err := s.getEventForOrgAdmin(ctx, idStr)
	if err != nil {
		return nil, err
	}
	if event.Status == entity.EventStatusPublished {
		return nil, utils.NewAppError(http.StatusConflict, "event is already published")
	}
	if event.RegistrationEnd == nil || event.RegistrationEnd.Before(time.Now()) {
		return nil, utils.NewAppError(http.StatusBadRequest, "registration period must be set and still in the future")
	}
	event.Status = entity.EventStatusPublished
	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}
	return s.toResponse(event, 0), nil
}

func (s *eventService) Close(ctx dto.AuthContext, idStr string) (*dto.EventResponse, error) {
	event, err := s.getEventForOrgAdmin(ctx, idStr)
	if err != nil {
		return nil, err
	}
	if event.Status != entity.EventStatusPublished {
		return nil, utils.NewAppError(http.StatusBadRequest, "only published events can be closed")
	}
	event.Status = entity.EventStatusClosed
	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}
	return s.toResponse(event, 0), nil
}

func (s *eventService) Archive(ctx dto.AuthContext, idStr string) (*dto.EventResponse, error) {
	event, err := s.getEventForOrgAdmin(ctx, idStr)
	if err != nil {
		return nil, err
	}
	if event.Status == entity.EventStatusArchived {
		return nil, utils.NewAppError(http.StatusConflict, "event is already archived")
	}
	event.Status = entity.EventStatusArchived
	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}
	return s.toResponse(event, 0), nil
}

func (s *eventService) getEventForOrgAdmin(ctx dto.AuthContext, idStr string) (*entity.Event, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid event id")
	}
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return nil, err
	}
	return event, nil
}

func validateEventDates(regStart, regEnd, start, end *time.Time) error {
	if regStart != nil && regEnd != nil && regEnd.Before(*regStart) {
		return utils.NewAppError(http.StatusBadRequest, "registration_end must be after registration_start")
	}
	if start != nil && end != nil && end.Before(*start) {
		return utils.NewAppError(http.StatusBadRequest, "end_date must be after start_date")
	}
	return nil
}

func (s *eventService) toResponse(event *entity.Event, count int) *dto.EventResponse {
	return &dto.EventResponse{
		ID:                event.ID.String(),
		OrganizationID:    event.OrganizationID.String(),
		Name:              event.Name,
		Slug:              event.Slug,
		Description:       event.Description,
		Location:          event.Location,
		StartDate:         event.StartDate,
		EndDate:           event.EndDate,
		RegistrationStart: event.RegistrationStart,
		RegistrationEnd:   event.RegistrationEnd,
		Quota:             event.Quota,
		Status:            event.Status,
		RegistrationCount: count,
		CreatedAt:         event.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         event.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
