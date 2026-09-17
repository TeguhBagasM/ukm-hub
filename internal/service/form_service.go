package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"handler/internal/dto"
	"handler/internal/entity"
	"handler/internal/repository"
	"handler/internal/utils"

	"github.com/google/uuid"
)

type FormService interface {
	GetByEvent(ctx dto.AuthContext, eventIDStr string) (*dto.FormResponse, error)
	CreateForm(ctx dto.AuthContext, eventIDStr string, req dto.CreateFormRequest) (*dto.FormResponse, error)
	UpdateForm(ctx dto.AuthContext, formIDStr string, req dto.UpdateFormRequest) (*dto.FormResponse, error)
	Publish(ctx dto.AuthContext, formIDStr string) (*dto.FormResponse, error)
	Unpublish(ctx dto.AuthContext, formIDStr string) (*dto.FormResponse, error)
	AddField(ctx dto.AuthContext, formIDStr string, req dto.FormFieldRequest) (*dto.FormFieldResponse, error)
	UpdateField(ctx dto.AuthContext, fieldIDStr string, req dto.FormFieldRequest) (*dto.FormFieldResponse, error)
	DeleteField(ctx dto.AuthContext, fieldIDStr string) error
	DuplicateField(ctx dto.AuthContext, formIDStr, fieldIDStr string) (*dto.FormFieldResponse, error)
	ReorderFields(ctx dto.AuthContext, formIDStr string, items []dto.ReorderFieldRequest) error
}

type formService struct {
	formRepo  repository.FormRepository
	fieldRepo repository.FormFieldRepository
	eventRepo repository.EventRepository
	access    AccessService
}

func NewFormService(formRepo repository.FormRepository, fieldRepo repository.FormFieldRepository, eventRepo repository.EventRepository, access AccessService) FormService {
	return &formService{formRepo: formRepo, fieldRepo: fieldRepo, eventRepo: eventRepo, access: access}
}

// ---------- Form ----------

func (s *formService) GetByEvent(ctx dto.AuthContext, eventIDStr string) (*dto.FormResponse, error) {
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid event id")
	}
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return nil, err
	}
	form, err := s.formRepo.FindByEvent(eventID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "form not found")
	}
	return s.loadFormResponse(form)
}

func (s *formService) CreateForm(ctx dto.AuthContext, eventIDStr string, req dto.CreateFormRequest) (*dto.FormResponse, error) {
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid event id")
	}
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return nil, err
	}
	if existing, _ := s.formRepo.FindByEvent(eventID); existing != nil {
		return nil, utils.NewAppError(http.StatusConflict, "event already has a form")
	}
	if err := validateWebhookURL(req.SubmitWebhookURL); err != nil {
		return nil, err
	}
	form := entity.Form{
		EventID:          eventID,
		Title:            req.Title,
		Description:      req.Description,
		SubmitWebhookURL: strings.TrimSpace(req.SubmitWebhookURL),
		IsPublished:      req.IsPublished,
	}
	if err := s.formRepo.Create(&form); err != nil {
		return nil, err
	}
	return s.loadFormResponse(&form)
}

func (s *formService) UpdateForm(ctx dto.AuthContext, formIDStr string, req dto.UpdateFormRequest) (*dto.FormResponse, error) {
	form, err := s.getFormWithAccess(ctx, formIDStr)
	if err != nil {
		return nil, err
	}
	if err := validateWebhookURL(req.SubmitWebhookURL); err != nil {
		return nil, err
	}
	if req.Title != "" {
		form.Title = req.Title
	}
	if req.Description != "" {
		form.Description = req.Description
	}
	if req.SubmitWebhookURL != "" {
		form.SubmitWebhookURL = strings.TrimSpace(req.SubmitWebhookURL)
	}
	if err := s.formRepo.Update(form); err != nil {
		return nil, err
	}
	return s.loadFormResponse(form)
}

func (s *formService) Publish(ctx dto.AuthContext, formIDStr string) (*dto.FormResponse, error) {
	form, err := s.getFormWithAccess(ctx, formIDStr)
	if err != nil {
		return nil, err
	}
	fields, err := s.fieldRepo.FindByForm(form.ID)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, utils.NewAppError(http.StatusBadRequest, "form must have at least one field before publishing")
	}
	form.IsPublished = true
	if err := s.formRepo.Update(form); err != nil {
		return nil, err
	}
	return s.loadFormResponse(form)
}

func (s *formService) Unpublish(ctx dto.AuthContext, formIDStr string) (*dto.FormResponse, error) {
	form, err := s.getFormWithAccess(ctx, formIDStr)
	if err != nil {
		return nil, err
	}
	form.IsPublished = false
	if err := s.formRepo.Update(form); err != nil {
		return nil, err
	}
	return s.loadFormResponse(form)
}

// ---------- Fields ----------

func (s *formService) AddField(ctx dto.AuthContext, formIDStr string, req dto.FormFieldRequest) (*dto.FormFieldResponse, error) {
	form, err := s.getFormWithAccess(ctx, formIDStr)
	if err != nil {
		return nil, err
	}
	if !entity.IsSupportedFieldType(req.Type) {
		return nil, utils.NewAppError(http.StatusBadRequest, "unsupported field type")
	}
	if err := validateFieldOptions(req.Type, req.Options); err != nil {
		return nil, err
	}
	opts, err := encodeOptions(req.Type, req.Options)
	if err != nil {
		return nil, err
	}
	sortOrder := req.SortOrder
	if sortOrder == 0 {
		max, err := s.fieldRepo.MaxSortOrder(form.ID)
		if err != nil {
			return nil, err
		}
		sortOrder = max + 1
	}
	field := entity.FormField{
		FormID:      form.ID,
		Label:       req.Label,
		Name:        req.Name,
		Type:        req.Type,
		Placeholder: req.Placeholder,
		Description: req.Description,
		Required:    req.Required != nil && *req.Required,
		Options:     opts,
		SortOrder:   sortOrder,
	}
	if err := s.fieldRepo.Create(&field); err != nil {
		return nil, err
	}
	return s.fieldToResponse(&field), nil
}

func (s *formService) UpdateField(ctx dto.AuthContext, fieldIDStr string, req dto.FormFieldRequest) (*dto.FormFieldResponse, error) {
	field, err := s.getFieldWithAccess(ctx, fieldIDStr)
	if err != nil {
		return nil, err
	}
	if req.Type != "" && !entity.IsSupportedFieldType(req.Type) {
		return nil, utils.NewAppError(http.StatusBadRequest, "unsupported field type")
	}
	if err := validateFieldOptions(req.Type, req.Options); err != nil {
		return nil, err
	}
	if req.Label != "" {
		field.Label = req.Label
	}
	if req.Name != "" {
		field.Name = req.Name
	}
	if req.Type != "" {
		field.Type = req.Type
	}
	if req.Placeholder != "" {
		field.Placeholder = req.Placeholder
	}
	if req.Description != "" {
		field.Description = req.Description
	}
	if req.Required != nil {
		field.Required = *req.Required
	}
	if req.Options != nil {
		opts, err := encodeOptions(field.Type, req.Options)
		if err != nil {
			return nil, err
		}
		field.Options = opts
	}
	if req.SortOrder > 0 {
		field.SortOrder = req.SortOrder
	}
	if err := s.fieldRepo.Update(field); err != nil {
		return nil, err
	}
	return s.fieldToResponse(field), nil
}

func (s *formService) DeleteField(ctx dto.AuthContext, fieldIDStr string) error {
	field, err := s.getFieldWithAccess(ctx, fieldIDStr)
	if err != nil {
		return err
	}
	return s.fieldRepo.Delete(field.ID)
}

func (s *formService) DuplicateField(ctx dto.AuthContext, formIDStr, fieldIDStr string) (*dto.FormFieldResponse, error) {
	form, err := s.getFormWithAccess(ctx, formIDStr)
	if err != nil {
		return nil, err
	}
	fieldID, err := uuid.Parse(fieldIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid field id")
	}
	field, err := s.fieldRepo.FindByID(fieldID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "field not found")
	}
	if field.FormID != form.ID {
		return nil, utils.NewAppError(http.StatusForbidden, "field does not belong to this form")
	}
	max, err := s.fieldRepo.MaxSortOrder(form.ID)
	if err != nil {
		return nil, err
	}
	dup := entity.FormField{
		FormID:      form.ID,
		Label:       field.Label,
		Name:        field.Name + "_copy",
		Type:        field.Type,
		Placeholder: field.Placeholder,
		Description: field.Description,
		Required:    field.Required,
		Options:     field.Options,
		SortOrder:   max + 1,
	}
	if err := s.fieldRepo.Create(&dup); err != nil {
		return nil, err
	}
	return s.fieldToResponse(&dup), nil
}

func (s *formService) ReorderFields(ctx dto.AuthContext, formIDStr string, items []dto.ReorderFieldRequest) error {
	form, err := s.getFormWithAccess(ctx, formIDStr)
	if err != nil {
		return err
	}
	fields, err := s.fieldRepo.FindByForm(form.ID)
	if err != nil {
		return err
	}
	byID := make(map[uuid.UUID]*entity.FormField, len(fields))
	for i := range fields {
		byID[fields[i].ID] = &fields[i]
	}
	for _, item := range items {
		id, err := uuid.Parse(item.ID)
		if err != nil {
			return utils.NewAppError(http.StatusBadRequest, "invalid field id")
		}
		field, ok := byID[id]
		if !ok {
			return utils.NewAppError(http.StatusBadRequest, "field does not belong to this form")
		}
		field.SortOrder = item.SortOrder
		if err := s.fieldRepo.Update(field); err != nil {
			return err
		}
	}
	return nil
}

// ---------- Helpers ----------

func (s *formService) getFormWithAccess(ctx dto.AuthContext, formIDStr string) (*entity.Form, error) {
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid form id")
	}
	form, err := s.formRepo.FindByID(formID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "form not found")
	}
	event, err := s.eventRepo.FindByID(form.EventID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *formService) getFieldWithAccess(ctx dto.AuthContext, fieldIDStr string) (*entity.FormField, error) {
	fieldID, err := uuid.Parse(fieldIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid field id")
	}
	field, err := s.fieldRepo.FindByID(fieldID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "field not found")
	}
	form, err := s.formRepo.FindByID(field.FormID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "form not found")
	}
	event, err := s.eventRepo.FindByID(form.EventID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return nil, err
	}
	return field, nil
}

func (s *formService) loadFormResponse(form *entity.Form) (*dto.FormResponse, error) {
	fields, err := s.fieldRepo.FindByForm(form.ID)
	if err != nil {
		return nil, err
	}
	fieldRes := make([]dto.FormFieldResponse, 0, len(fields))
	for i := range fields {
		fieldRes = append(fieldRes, *s.fieldToResponse(&fields[i]))
	}
	return &dto.FormResponse{
		ID:               form.ID.String(),
		EventID:          form.EventID.String(),
		Title:            form.Title,
		Description:      form.Description,
		SubmitWebhookURL: form.SubmitWebhookURL,
		IsPublished:      form.IsPublished,
		Fields:           fieldRes,
		CreatedAt:        form.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        form.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *formService) fieldToResponse(field *entity.FormField) *dto.FormFieldResponse {
	options := []string{}
	if field.Options != "" {
		_ = json.Unmarshal([]byte(field.Options), &options)
	}
	return &dto.FormFieldResponse{
		ID:          field.ID.String(),
		FormID:      field.FormID.String(),
		Label:       field.Label,
		Name:        field.Name,
		Type:        field.Type,
		Placeholder: field.Placeholder,
		Description: field.Description,
		Required:    field.Required,
		Options:     options,
		SortOrder:   field.SortOrder,
	}
}

func validateWebhookURL(raw string) error {
	url := strings.TrimSpace(raw)
	if url == "" {
		return nil
	}
	if !strings.HasPrefix(url, "https://") {
		return utils.NewAppError(http.StatusBadRequest, "submit_webhook_url must use https://")
	}
	return nil
}

func validateFieldOptions(fieldType string, options []string) error {
	if fieldType == entity.FormFieldTypeSelect ||
		fieldType == entity.FormFieldTypeRadio ||
		fieldType == entity.FormFieldTypeCheckbox {
		if len(options) == 0 {
			return utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("field type %s requires at least one option", fieldType))
		}
	}
	return nil
}

func encodeOptions(fieldType string, options []string) (string, error) {
	if fieldType == "" {
		return "", nil
	}
	if fieldType != entity.FormFieldTypeSelect &&
		fieldType != entity.FormFieldTypeRadio &&
		fieldType != entity.FormFieldTypeCheckbox {
		return "", nil
	}
	if len(options) == 0 {
		return "", utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("field type %s requires at least one option", fieldType))
	}
	b, err := json.Marshal(options)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
