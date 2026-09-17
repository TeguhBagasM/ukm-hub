package service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"handler/internal/dto"
	"handler/internal/entity"
	"handler/internal/repository"
	"handler/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberService interface {
	List(ctx dto.AuthContext, orgIDStr, search, divisionID string, page, perPage int) (*dto.MemberListResponse, error)
	Get(ctx dto.AuthContext, idStr string) (*dto.MemberResponse, error)
	Convert(ctx dto.AuthContext, registrationIDStr string, req dto.ConvertMemberRequest) (*dto.MemberResponse, error)
	Update(ctx dto.AuthContext, idStr string, req dto.UpdateMemberRequest) (*dto.MemberResponse, error)
	Delete(ctx dto.AuthContext, idStr string) error
}

type memberService struct {
	memberRepo repository.MemberRepository
	regRepo    repository.RegistrationRepository
	eventRepo  repository.EventRepository
	formRepo   repository.FormRepository
	fieldRepo  repository.FormFieldRepository
	divRepo    repository.DivisionRepository
	orgRepo    repository.OrganizationRepository
	access     AccessService
}

func NewMemberService(
	memberRepo repository.MemberRepository,
	regRepo repository.RegistrationRepository,
	eventRepo repository.EventRepository,
	formRepo repository.FormRepository,
	fieldRepo repository.FormFieldRepository,
	divRepo repository.DivisionRepository,
	orgRepo repository.OrganizationRepository,
	access AccessService,
) MemberService {
	return &memberService{
		memberRepo: memberRepo,
		regRepo:    regRepo,
		eventRepo:  eventRepo,
		formRepo:   formRepo,
		fieldRepo:  fieldRepo,
		divRepo:    divRepo,
		orgRepo:    orgRepo,
		access:     access,
	}
}

func (s *memberService) List(ctx dto.AuthContext, orgIDStr, search, divisionID string, page, perPage int) (*dto.MemberListResponse, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if err := s.access.MustAccessOrganization(ctx, orgID); err != nil {
		return nil, err
	}
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	members, total, err := s.memberRepo.ListByOrganization(orgID, search, divisionID, page, perPage)
	if err != nil {
		return nil, err
	}
	items := make([]dto.MemberResponse, 0, len(members))
	for i := range members {
		items = append(items, *s.toResponse(&members[i]))
	}
	return &dto.MemberListResponse{
		Items:      items,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}, nil
}

func (s *memberService) Get(ctx dto.AuthContext, idStr string) (*dto.MemberResponse, error) {
	member, err := s.getWithAccess(ctx, idStr)
	if err != nil {
		return nil, err
	}
	return s.toResponse(member), nil
}

func (s *memberService) Convert(ctx dto.AuthContext, registrationIDStr string, req dto.ConvertMemberRequest) (*dto.MemberResponse, error) {
	regID, err := uuid.Parse(registrationIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid registration id")
	}
	reg, err := s.regRepo.FindByID(regID)
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
	if reg.Status != entity.RegistrationStatusAccepted {
		return nil, utils.NewAppError(http.StatusBadRequest, "only accepted applicants can be converted to members")
	}

	existing, err := s.memberRepo.FindByRegistrationIDUnscoped(regID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil && !existing.DeletedAt.Valid {
		return nil, utils.NewAppError(http.StatusConflict, "this applicant has already been converted to a member")
	}

	derived, err := s.deriveMemberFields(reg, event.OrganizationID)
	if err != nil {
		return nil, err
	}

	name := firstNonEmpty(req.Name, derived.Name)
	email := firstNonEmpty(req.Email, derived.Email)
	phone := firstNonEmpty(req.Phone, derived.Phone)
	studentID := firstNonEmpty(req.StudentID, derived.StudentID)
	status := req.Status
	if status == "" {
		status = entity.MemberStatusActive
	}
	if name == "" {
		return nil, utils.NewAppError(http.StatusBadRequest, "member name is required")
	}

	var divisionID *uuid.UUID
	if req.DivisionID != "" {
		divID, err := uuid.Parse(req.DivisionID)
		if err != nil {
			return nil, utils.NewAppError(http.StatusBadRequest, "invalid division id")
		}
		div, err := s.divRepo.FindByID(divID)
		if err != nil || div.OrganizationID != event.OrganizationID {
			return nil, utils.NewAppError(http.StatusBadRequest, "division does not belong to this organization")
		}
		divisionID = &divID
	} else if derived.DivisionID != nil {
		divisionID = derived.DivisionID
	}

	member := entity.Member{}
	if existing != nil {
		member = *existing
	}
	member.OrganizationID = event.OrganizationID
	member.DivisionID = divisionID
	member.RegistrationID = &regID
	member.Name = name
	member.Email = email
	member.Phone = phone
	member.StudentID = studentID
	member.Status = status
	member.JoinedAt = time.Now()

	if existing != nil {
		if err := s.memberRepo.Restore(&member); err != nil {
			return nil, err
		}
	} else if err := s.memberRepo.Create(&member); err != nil {
		return nil, err
	}
	return s.toResponse(&member), nil
}

func (s *memberService) Update(ctx dto.AuthContext, idStr string, req dto.UpdateMemberRequest) (*dto.MemberResponse, error) {
	member, err := s.getWithAccess(ctx, idStr)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		member.Name = req.Name
	}
	if req.Email != "" {
		member.Email = req.Email
	}
	if req.Phone != "" {
		member.Phone = req.Phone
	}
	if req.StudentID != "" {
		member.StudentID = req.StudentID
	}
	if req.Status != "" {
		member.Status = req.Status
	}
	if req.DivisionID != "" {
		divID, err := uuid.Parse(req.DivisionID)
		if err != nil {
			return nil, utils.NewAppError(http.StatusBadRequest, "invalid division id")
		}
		div, err := s.divRepo.FindByID(divID)
		if err != nil || div.OrganizationID != member.OrganizationID {
			return nil, utils.NewAppError(http.StatusBadRequest, "division does not belong to this organization")
		}
		member.DivisionID = &divID
	}
	if err := s.memberRepo.Update(member); err != nil {
		return nil, err
	}
	return s.toResponse(member), nil
}

func (s *memberService) Delete(ctx dto.AuthContext, idStr string) error {
	member, err := s.getWithAccess(ctx, idStr)
	if err != nil {
		return err
	}
	return s.memberRepo.Delete(member.ID)
}

func (s *memberService) getWithAccess(ctx dto.AuthContext, idStr string) (*entity.Member, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid member id")
	}
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "member not found")
	}
	if err := s.access.MustAccessOrganization(ctx, member.OrganizationID); err != nil {
		return nil, err
	}
	return member, nil
}

type derivedMember struct {
	Name       string
	Email      string
	Phone      string
	StudentID  string
	DivisionID *uuid.UUID
}

func (s *memberService) deriveMemberFields(reg *entity.Registration, orgID uuid.UUID) (derivedMember, error) {
	var result derivedMember

	form, err := s.formRepo.FindByEvent(reg.EventID)
	if err != nil {
		return result, nil
	}
	fields, err := s.fieldRepo.FindByForm(form.ID)
	if err != nil {
		return result, err
	}
	fieldByID := make(map[uuid.UUID]*entity.FormField, len(fields))
	for i := range fields {
		fieldByID[fields[i].ID] = &fields[i]
	}
	answers, err := s.regRepo.FindAnswersByRegistrationIDs([]uuid.UUID{reg.ID})
	if err != nil {
		return result, err
	}

	divisions, err := s.divRepo.FindByOrganization(orgID)
	if err != nil {
		return result, err
	}

	for i := range answers {
		field := fieldByID[answers[i].FormFieldID]
		if field == nil || answers[i].Value == "" {
			continue
		}
		key := strings.ToLower(field.Name + " " + field.Label)
		value := strings.TrimSpace(answers[i].Value)

		switch {
		case field.Type == entity.FormFieldTypeEmail || strings.Contains(key, "email"):
			setIfEmpty(&result.Email, value)
		case strings.Contains(key, "nama"):
			setIfEmpty(&result.Name, value)
		case field.Type == entity.FormFieldTypePhone || containsAny(key, "phone", "telepon", "whatsapp", "hp", "wa"):
			setIfEmpty(&result.Phone, value)
		case containsAny(key, "nim", "nrp", "student"):
			setIfEmpty(&result.StudentID, value)
		}

		if result.DivisionID == nil && (strings.Contains(key, "divisi") || strings.Contains(key, "division") ||
			field.Type == entity.FormFieldTypeSelect || field.Type == entity.FormFieldTypeRadio) {
			for j := range divisions {
				if strings.EqualFold(divisions[j].Name, value) {
					id := divisions[j].ID
					result.DivisionID = &id
					break
				}
			}
		}
	}
	return result, nil
}

func (s *memberService) toResponse(member *entity.Member) *dto.MemberResponse {
	res := &dto.MemberResponse{
		ID:             member.ID.String(),
		OrganizationID: member.OrganizationID.String(),
		Name:           member.Name,
		Email:          member.Email,
		Phone:          member.Phone,
		StudentID:      member.StudentID,
		Status:         member.Status,
		JoinedAt:       member.JoinedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:      member.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      member.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if member.DivisionID != nil {
		v := member.DivisionID.String()
		res.DivisionID = &v
	}
	if member.RegistrationID != nil {
		v := member.RegistrationID.String()
		res.RegistrationID = &v
	}
	return res
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func setIfEmpty(target *string, value string) {
	if *target == "" {
		*target = value
	}
}

func containsAny(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}
