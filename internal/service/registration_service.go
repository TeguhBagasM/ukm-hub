package service

import (
	"net/http"
	"strconv"

	"ukm-hub/internal/dto"
	"ukm-hub/internal/entity"
	"ukm-hub/internal/repository"
	"ukm-hub/internal/utils"

	"github.com/google/uuid"
)

type RegistrationService interface {
	ListByEvent(ctx dto.AuthContext, eventIDStr string, status, search string, page, perPage int) (*dto.RegistrationListResponse, error)
	Get(ctx dto.AuthContext, idStr string) (*dto.RegistrationResponse, error)
	Accept(ctx dto.AuthContext, idStr string) (*dto.RegistrationResponse, error)
	Reject(ctx dto.AuthContext, idStr string, reason string) (*dto.RegistrationResponse, error)
}

type registrationService struct {
	regRepo   repository.RegistrationRepository
	eventRepo repository.EventRepository
	formRepo  repository.FormRepository
	fieldRepo repository.FormFieldRepository
	access    AccessService
}

func NewRegistrationService(
	regRepo repository.RegistrationRepository,
	eventRepo repository.EventRepository,
	formRepo repository.FormRepository,
	fieldRepo repository.FormFieldRepository,
	access AccessService,
) RegistrationService {
	return &registrationService{
		regRepo:   regRepo,
		eventRepo: eventRepo,
		formRepo:  formRepo,
		fieldRepo: fieldRepo,
		access:    access,
	}
}

func (s *registrationService) ListByEvent(ctx dto.AuthContext, eventIDStr string, status, search string, page, perPage int) (*dto.RegistrationListResponse, error) {
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
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	regs, total, err := s.regRepo.ListByEvent(eventID, status, search, page, perPage)
	if err != nil {
		return nil, err
	}
	fieldByID, err := s.fieldMap(eventID)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(regs))
	for i := range regs {
		ids = append(ids, regs[i].ID)
	}
	answers, err := s.regRepo.FindAnswersByRegistrationIDs(ids)
	if err != nil {
		return nil, err
	}
	answersByReg := make(map[uuid.UUID][]dto.RegistrationAnswerResponse)
	for i := range answers {
		a := answers[i]
		answersByReg[a.RegistrationID] = append(answersByReg[a.RegistrationID], toAnswerResponse(a, fieldByID[a.FormFieldID]))
	}

	items := make([]dto.RegistrationResponse, 0, len(regs))
	for i := range regs {
		items = append(items, *s.toResponse(&regs[i], answersByReg[regs[i].ID]))
	}

	return &dto.RegistrationListResponse{
		Items:      items,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}, nil
}

func (s *registrationService) Get(ctx dto.AuthContext, idStr string) (*dto.RegistrationResponse, error) {
	reg, err := s.getWithAccess(ctx, idStr)
	if err != nil {
		return nil, err
	}
	fieldByID, err := s.fieldMap(reg.EventID)
	if err != nil {
		return nil, err
	}
	answers, err := s.regRepo.FindAnswersByRegistrationIDs([]uuid.UUID{reg.ID})
	if err != nil {
		return nil, err
	}
	answerRes := make([]dto.RegistrationAnswerResponse, 0, len(answers))
	for i := range answers {
		answerRes = append(answerRes, toAnswerResponse(answers[i], fieldByID[answers[i].FormFieldID]))
	}
	return s.toResponse(reg, answerRes), nil
}

func (s *registrationService) Accept(ctx dto.AuthContext, idStr string) (*dto.RegistrationResponse, error) {
	reg, err := s.getWithAccess(ctx, idStr)
	if err != nil {
		return nil, err
	}
	if reg.Status == entity.RegistrationStatusAccepted {
		return nil, utils.NewAppError(http.StatusConflict, "registration is already accepted")
	}
	reg.Status = entity.RegistrationStatusAccepted
	reg.RejectionReason = ""
	if err := s.regRepo.Update(reg); err != nil {
		return nil, err
	}
	return s.Get(ctx, idStr)
}

func (s *registrationService) Reject(ctx dto.AuthContext, idStr string, reason string) (*dto.RegistrationResponse, error) {
	reg, err := s.getWithAccess(ctx, idStr)
	if err != nil {
		return nil, err
	}
	if reg.Status == entity.RegistrationStatusRejected {
		return nil, utils.NewAppError(http.StatusConflict, "registration is already rejected")
	}
	reg.Status = entity.RegistrationStatusRejected
	reg.RejectionReason = reason
	if err := s.regRepo.Update(reg); err != nil {
		return nil, err
	}
	return s.Get(ctx, idStr)
}

func (s *registrationService) getWithAccess(ctx dto.AuthContext, idStr string) (*entity.Registration, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid registration id")
	}
	reg, err := s.regRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "registration not found")
	}
	event, err := s.eventRepo.FindByID(reg.EventID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "event not found")
	}
	if err := s.access.MustAccessOrganization(ctx, event.OrganizationID); err != nil {
		return nil, err
	}
	return reg, nil
}

func (s *registrationService) fieldMap(eventID uuid.UUID) (map[uuid.UUID]*entity.FormField, error) {
	result := make(map[uuid.UUID]*entity.FormField)
	form, err := s.formRepo.FindByEvent(eventID)
	if err != nil {
		return result, nil
	}
	fields, err := s.fieldRepo.FindByForm(form.ID)
	if err != nil {
		return nil, err
	}
	for i := range fields {
		result[fields[i].ID] = &fields[i]
	}
	return result, nil
}

func (s *registrationService) toResponse(reg *entity.Registration, answers []dto.RegistrationAnswerResponse) *dto.RegistrationResponse {
	if answers == nil {
		answers = []dto.RegistrationAnswerResponse{}
	}
	return &dto.RegistrationResponse{
		ID:              reg.ID.String(),
		EventID:         reg.EventID.String(),
		Status:          reg.Status,
		RejectionReason: reg.RejectionReason,
		SubmittedAt:     reg.SubmittedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:       reg.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       reg.UpdatedAt.Format("2006-01-02 15:04:05"),
		Answers:         answers,
	}
}

func toAnswerResponse(answer entity.RegistrationAnswer, field *entity.FormField) dto.RegistrationAnswerResponse {
	res := dto.RegistrationAnswerResponse{
		FieldID: answer.FormFieldID.String(),
		Value:   answer.Value,
	}
	if field != nil {
		res.Label = field.Label
		res.Name = field.Name
		res.Type = field.Type
	} else {
		res.Name = answer.FormFieldID.String()
	}
	return res
}

func ParsePagination(pageStr, perPageStr string) (int, int) {
	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)
	return page, perPage
}
