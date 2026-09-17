package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ukm-hub/internal/dto"
	"ukm-hub/internal/entity"
	"ukm-hub/internal/integration"
	"ukm-hub/internal/repository"
	"ukm-hub/internal/utils"
)

type PublicService interface {
	GetEvent(slug string) (*dto.PublicEventResponse, error)
	GetForm(slug string) (*dto.PublicFormResponse, error)
	Register(slug string, req dto.PublicRegistrationRequest) (*dto.PublicRegistrationResponse, error)
}

type publicService struct {
	eventRepo repository.EventRepository
	formRepo  repository.FormRepository
	fieldRepo repository.FormFieldRepository
	orgRepo   repository.OrganizationRepository
	regRepo   repository.RegistrationRepository
}

func NewPublicService(
	eventRepo repository.EventRepository,
	formRepo repository.FormRepository,
	fieldRepo repository.FormFieldRepository,
	orgRepo repository.OrganizationRepository,
	regRepo repository.RegistrationRepository,
) PublicService {
	return &publicService{
		eventRepo: eventRepo,
		formRepo:  formRepo,
		fieldRepo: fieldRepo,
		orgRepo:   orgRepo,
		regRepo:   regRepo,
	}
}

func (s *publicService) GetEvent(slug string) (*dto.PublicEventResponse, error) {
	event, err := s.eventRepo.FindBySlug(slug)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if event.Status != entity.EventStatusPublished {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found or not published")
	}
	count, err := s.regRepo.CountByEventID(event.ID)
	if err != nil {
		return nil, err
	}
	org, err := s.orgRepo.FindByID(event.OrganizationID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "organization not found")
	}
	res := s.eventToResponse(event, org, count)
	return &res, nil
}

func (s *publicService) GetForm(slug string) (*dto.PublicFormResponse, error) {
	event, err := s.eventRepo.FindBySlug(slug)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if event.Status != entity.EventStatusPublished {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found or not published")
	}
	form, err := s.formRepo.FindByEvent(event.ID)
	if err != nil || !form.IsPublished {
		return nil, utils.NewAppError(http.StatusNotFound, "registration form is not available")
	}
	fields, err := s.fieldRepo.FindByForm(form.ID)
	if err != nil {
		return nil, err
	}
	org, err := s.orgRepo.FindByID(event.OrganizationID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "organization not found")
	}
	count, err := s.regRepo.CountByEventID(event.ID)
	if err != nil {
		return nil, err
	}
	fieldRes := make([]dto.FormFieldResponse, 0, len(fields))
	for i := range fields {
		fieldRes = append(fieldRes, dto.FormFieldResponse{
			ID:          fields[i].ID.String(),
			FormID:      fields[i].FormID.String(),
			Label:       fields[i].Label,
			Name:        fields[i].Name,
			Type:        fields[i].Type,
			Placeholder: fields[i].Placeholder,
			Description: fields[i].Description,
			Required:    fields[i].Required,
			Options:     decodeOptions(fields[i].Options),
			SortOrder:   fields[i].SortOrder,
		})
	}
	return &dto.PublicFormResponse{
		Event: s.eventToResponse(event, org, count),
		Form: dto.PublicFormInfo{
			Title:       form.Title,
			Description: form.Description,
			Fields:      fieldRes,
		},
	}, nil
}

func (s *publicService) Register(slug string, req dto.PublicRegistrationRequest) (*dto.PublicRegistrationResponse, error) {
	event, err := s.eventRepo.FindBySlug(slug)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	org, err := s.orgRepo.FindByID(event.OrganizationID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "organization not found")
	}
	count, err := s.regRepo.CountByEventID(event.ID)
	if err != nil {
		return nil, err
	}

	if err := s.ensureRegistrationOpen(event, count); err != nil {
		return nil, err
	}

	form, err := s.formRepo.FindByEvent(event.ID)
	if err != nil || !form.IsPublished {
		return nil, utils.NewAppError(http.StatusBadRequest, "registration form is not available")
	}
	fields, err := s.fieldRepo.FindByForm(form.ID)
	if err != nil {
		return nil, err
	}

	fieldByID := make(map[string]*entity.FormField, len(fields))
	for i := range fields {
		fieldByID[fields[i].ID.String()] = &fields[i]
	}

	answersByField := make(map[string]string, len(req.Answers))
	for _, a := range req.Answers {
		field, ok := fieldByID[a.FieldID]
		if !ok {
			return nil, utils.NewAppError(http.StatusBadRequest, "invalid form field")
		}
		if _, dup := answersByField[a.FieldID]; dup {
			return nil, utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("duplicate answer for %s", field.Label))
		}
		answersByField[a.FieldID] = strings.TrimSpace(a.Value)
	}

	answers := make([]entity.RegistrationAnswer, 0, len(fields))
	webhookAnswers := make([]dto.WebhookAnswer, 0, len(fields))
	answersMap := make(map[string]string, len(fields))

	for i := range fields {
		field := &fields[i]
		value := answersByField[field.ID.String()]
		if err := validateAnswerValue(field, value); err != nil {
			return nil, err
		}
		if value == "" {
			continue
		}
		answers = append(answers, entity.RegistrationAnswer{
			FormFieldID: field.ID,
			Value:       value,
		})
		webhookAnswers = append(webhookAnswers, dto.WebhookAnswer{
			FieldID: field.ID.String(),
			Label:   field.Label,
			Name:    field.Name,
			Type:    field.Type,
			Value:   value,
		})
		answersMap[field.Name] = value

		if field.Type == entity.FormFieldTypeEmail {
			exists, err := s.regRepo.ExistsAnswer(event.ID, field.ID, value)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, utils.NewAppError(http.StatusConflict, "this email has already registered for this event")
			}
		}
	}

	reg := entity.Registration{
		EventID:     event.ID,
		Status:      entity.RegistrationStatusPending,
		SubmittedAt: time.Now(),
	}
	if err := s.regRepo.CreateWithAnswers(&reg, answers); err != nil {
		return nil, err
	}

	if form.SubmitWebhookURL != "" {
		integration.SendFormWebhook(form.SubmitWebhookURL, dto.RegistrationWebhookPayload{
			Event:          s.eventToResponse(event, org, count+1),
			RegistrationID: reg.ID.String(),
			Status:         reg.Status,
			SubmittedAt:    reg.SubmittedAt.Format(time.RFC3339),
			Answers:        webhookAnswers,
			AnswersMap:     answersMap,
		})
	}

	return &dto.PublicRegistrationResponse{
		ID:          reg.ID.String(),
		Status:      reg.Status,
		SubmittedAt: reg.SubmittedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *publicService) ensureRegistrationOpen(event *entity.Event, count int) error {
	if event.Status != entity.EventStatusPublished {
		return utils.NewAppError(http.StatusBadRequest, "registration is not open")
	}
	now := time.Now()
	if event.RegistrationStart != nil && now.Before(*event.RegistrationStart) {
		return utils.NewAppError(http.StatusBadRequest, "registration has not started yet")
	}
	if event.RegistrationEnd != nil && now.After(*event.RegistrationEnd) {
		return utils.NewAppError(http.StatusBadRequest, "registration is closed")
	}
	if event.Quota > 0 && count >= event.Quota {
		return utils.NewAppError(http.StatusConflict, "quota exceeded")
	}
	return nil
}

func (s *publicService) eventToResponse(event *entity.Event, org *entity.Organization, count int) dto.PublicEventResponse {
	open := event.Status == entity.EventStatusPublished
	now := time.Now()
	if event.RegistrationStart != nil && now.Before(*event.RegistrationStart) {
		open = false
	}
	if event.RegistrationEnd != nil && now.After(*event.RegistrationEnd) {
		open = false
	}
	if event.Quota > 0 && count >= event.Quota {
		open = false
	}
	orgInfo := dto.PublicOrganizationInfo{}
	if org != nil {
		orgInfo = dto.PublicOrganizationInfo{ID: org.ID.String(), Name: org.Name, Slug: org.Slug}
	}
	return dto.PublicEventResponse{
		ID:                 event.ID.String(),
		Organization:       orgInfo,
		Name:               event.Name,
		Slug:               event.Slug,
		Description:        event.Description,
		Location:           event.Location,
		StartDate:          event.StartDate,
		EndDate:            event.EndDate,
		RegistrationStart:  event.RegistrationStart,
		RegistrationEnd:    event.RegistrationEnd,
		Quota:              event.Quota,
		RegistrationCount:  count,
		IsRegistrationOpen: open,
	}
}

func decodeOptions(raw string) []string {
	options := []string{}
	if raw == "" {
		return options
	}
	_ = json.Unmarshal([]byte(raw), &options)
	return options
}

var (
	numberRegex = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
)

func validateAnswerValue(field *entity.FormField, raw string) error {
	value := strings.TrimSpace(raw)
	if value == "" {
		if field.Required {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s is required", field.Label))
		}
		return nil
	}

	switch field.Type {
	case entity.FormFieldTypeEmail:
		if _, err := mail.ParseAddress(value); err != nil {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s must be a valid email", field.Label))
		}
	case entity.FormFieldTypeNumber:
		if !numberRegex.MatchString(value) {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s must be a valid number", field.Label))
		}
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s must be a valid number", field.Label))
		}
	case entity.FormFieldTypeURL:
		u, err := url.ParseRequestURI(value)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s must be a valid URL", field.Label))
		}
	case entity.FormFieldTypeDate:
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s must use YYYY-MM-DD format", field.Label))
		}
	case entity.FormFieldTypePhone:
		digits := onlyDigits(value)
		if len(digits) < 8 || len(digits) > 15 {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s must be a valid phone number", field.Label))
		}
	case entity.FormFieldTypeSelect, entity.FormFieldTypeRadio:
		if !optionExists(field.Options, value) {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s has an invalid option", field.Label))
		}
	case entity.FormFieldTypeCheckbox:
		options := decodeOptions(field.Options)
		for _, part := range strings.Split(value, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if !containsString(options, part) {
				return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("%s has an invalid option", field.Label))
			}
		}
	}
	return nil
}

func optionExists(rawOptions, value string) bool {
	return containsString(decodeOptions(rawOptions), value)
}

func containsString(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func onlyDigits(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
