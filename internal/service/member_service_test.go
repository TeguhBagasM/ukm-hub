package service

import (
	"net/http"
	"testing"
	"time"

	"ukm-hub/internal/dto"
	"ukm-hub/internal/entity"
	"ukm-hub/internal/repository"
	"ukm-hub/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeAccess struct {
	err error
}

func (f fakeAccess) CanAccessOrganization(dto.AuthContext, uuid.UUID) (bool, error) {
	return f.err == nil, f.err
}

func (f fakeAccess) MustAccessOrganization(dto.AuthContext, uuid.UUID) error {
	return f.err
}

func (f fakeAccess) FindAccessibleOrganizations(dto.AuthContext) ([]entity.Organization, error) {
	return nil, f.err
}

type fakeMemberRepo struct {
	repository.MemberRepository
	existing *entity.Member
	created  *entity.Member
	restored *entity.Member
}

func (f *fakeMemberRepo) FindByRegistrationIDUnscoped(id uuid.UUID) (*entity.Member, error) {
	if f.existing != nil {
		return f.existing, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeMemberRepo) Create(member *entity.Member) error {
	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	f.created = member
	return nil
}

func (f *fakeMemberRepo) Restore(member *entity.Member) error {
	f.restored = member
	return nil
}

type fakeRegRepo struct {
	repository.RegistrationRepository
	registration *entity.Registration
}

func (f *fakeRegRepo) FindByID(id uuid.UUID) (*entity.Registration, error) {
	if f.registration == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.registration, nil
}

func (f *fakeRegRepo) FindAnswersByRegistrationIDs(ids []uuid.UUID) ([]entity.RegistrationAnswer, error) {
	return nil, nil
}

type fakeEventRepo struct {
	repository.EventRepository
	event *entity.Event
}

func (f *fakeEventRepo) FindByID(id uuid.UUID) (*entity.Event, error) {
	if f.event == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.event, nil
}

type fakeFormRepo struct {
	repository.FormRepository
}

func (f *fakeFormRepo) FindByEvent(eventID uuid.UUID) (*entity.Form, error) {
	return nil, gorm.ErrRecordNotFound
}

type fakeFieldRepo struct {
	repository.FormFieldRepository
}

func (f *fakeFieldRepo) FindByForm(formID uuid.UUID) ([]entity.FormField, error) {
	return nil, nil
}

type fakeDivisionRepo struct {
	repository.DivisionRepository
}

func (f *fakeDivisionRepo) FindByOrganization(orgID uuid.UUID) ([]entity.Division, error) {
	return nil, nil
}

func (f *fakeDivisionRepo) FindByID(id uuid.UUID) (*entity.Division, error) {
	return nil, gorm.ErrRecordNotFound
}

func newTestMemberService(memberRepo *fakeMemberRepo, regRepo *fakeRegRepo, eventRepo *fakeEventRepo, access AccessService) MemberService {
	return NewMemberService(memberRepo, regRepo, eventRepo, &fakeFormRepo{}, &fakeFieldRepo{}, &fakeDivisionRepo{}, nil, access)
}

func authCtx() dto.AuthContext {
	return dto.AuthContext{UserID: uuid.New().String(), Role: utils.RoleOrgAdmin}
}

func TestConvertRejectsNonAcceptedApplicant(t *testing.T) {
	orgID, eventID, regID := uuid.New(), uuid.New(), uuid.New()
	memberRepo := &fakeMemberRepo{}
	svc := newTestMemberService(memberRepo, &fakeRegRepo{
		registration: &entity.Registration{ID: regID, EventID: eventID, Status: entity.RegistrationStatusPending},
	}, &fakeEventRepo{event: &entity.Event{ID: eventID, OrganizationID: orgID}}, fakeAccess{})

	_, err := svc.Convert(authCtx(), regID.String(), dto.ConvertMemberRequest{Name: "Ani"})
	if got := appErrorStatus(t, err); got != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, got)
	}
	if memberRepo.created != nil {
		t.Fatal("member must not be created for a non-accepted applicant")
	}
}

func TestConvertPreventsDuplicateConversion(t *testing.T) {
	orgID, eventID, regID := uuid.New(), uuid.New(), uuid.New()
	existing := &entity.Member{ID: uuid.New(), OrganizationID: orgID, RegistrationID: &regID, Name: "Ani"}
	memberRepo := &fakeMemberRepo{existing: existing}
	svc := newTestMemberService(memberRepo, &fakeRegRepo{
		registration: &entity.Registration{ID: regID, EventID: eventID, Status: entity.RegistrationStatusAccepted},
	}, &fakeEventRepo{event: &entity.Event{ID: eventID, OrganizationID: orgID}}, fakeAccess{})

	_, err := svc.Convert(authCtx(), regID.String(), dto.ConvertMemberRequest{Name: "Ani"})
	if got := appErrorStatus(t, err); got != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, got)
	}
	if memberRepo.created != nil {
		t.Fatal("duplicate conversion must not create a new member")
	}
}

func TestConvertRestoresSoftDeletedMember(t *testing.T) {
	orgID, eventID, regID := uuid.New(), uuid.New(), uuid.New()
	existing := &entity.Member{
		ID:             uuid.New(),
		OrganizationID: orgID,
		RegistrationID: &regID,
		Name:           "Ani",
		DeletedAt:      gorm.DeletedAt{Time: time.Now(), Valid: true},
	}
	memberRepo := &fakeMemberRepo{existing: existing}
	svc := newTestMemberService(memberRepo, &fakeRegRepo{
		registration: &entity.Registration{ID: regID, EventID: eventID, Status: entity.RegistrationStatusAccepted},
	}, &fakeEventRepo{event: &entity.Event{ID: eventID, OrganizationID: orgID}}, fakeAccess{})

	res, err := svc.Convert(authCtx(), regID.String(), dto.ConvertMemberRequest{Name: "Ani Restored"})
	if err != nil {
		t.Fatalf("expected conversion to restore the soft-deleted member, got error: %v", err)
	}
	if memberRepo.restored == nil {
		t.Fatal("expected the soft-deleted member to be restored")
	}
	if memberRepo.created != nil {
		t.Fatal("restoring must not create a new member row")
	}
	if res.ID != existing.ID.String() {
		t.Fatalf("expected restored member id %s, got %s", existing.ID, res.ID)
	}
	if res.Name != "Ani Restored" {
		t.Fatalf("expected updated name, got %s", res.Name)
	}
}

func TestConvertCreatesMember(t *testing.T) {
	orgID, eventID, regID := uuid.New(), uuid.New(), uuid.New()
	memberRepo := &fakeMemberRepo{}
	svc := newTestMemberService(memberRepo, &fakeRegRepo{
		registration: &entity.Registration{ID: regID, EventID: eventID, Status: entity.RegistrationStatusAccepted},
	}, &fakeEventRepo{event: &entity.Event{ID: eventID, OrganizationID: orgID}}, fakeAccess{})

	res, err := svc.Convert(authCtx(), regID.String(), dto.ConvertMemberRequest{Name: "Ani", Email: "ani@kampus.ac.id"})
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}
	if memberRepo.created == nil {
		t.Fatal("expected a new member to be created")
	}
	if res.Status != entity.MemberStatusActive {
		t.Fatalf("expected default status %s, got %s", entity.MemberStatusActive, res.Status)
	}
	if res.RegistrationID == nil || *res.RegistrationID != regID.String() {
		t.Fatalf("expected registration id %s to be linked", regID)
	}
}

func TestConvertRequiresMemberName(t *testing.T) {
	orgID, eventID, regID := uuid.New(), uuid.New(), uuid.New()
	memberRepo := &fakeMemberRepo{}
	svc := newTestMemberService(memberRepo, &fakeRegRepo{
		registration: &entity.Registration{ID: regID, EventID: eventID, Status: entity.RegistrationStatusAccepted},
	}, &fakeEventRepo{event: &entity.Event{ID: eventID, OrganizationID: orgID}}, fakeAccess{})

	_, err := svc.Convert(authCtx(), regID.String(), dto.ConvertMemberRequest{})
	if got := appErrorStatus(t, err); got != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, got)
	}
}

func TestConvertRejectsInaccessibleOrganization(t *testing.T) {
	orgID, eventID, regID := uuid.New(), uuid.New(), uuid.New()
	memberRepo := &fakeMemberRepo{}
	svc := newTestMemberService(memberRepo, &fakeRegRepo{
		registration: &entity.Registration{ID: regID, EventID: eventID, Status: entity.RegistrationStatusAccepted},
	}, &fakeEventRepo{event: &entity.Event{ID: eventID, OrganizationID: orgID}},
		fakeAccess{err: utils.NewAppError(http.StatusForbidden, "forbidden")})

	_, err := svc.Convert(authCtx(), regID.String(), dto.ConvertMemberRequest{Name: "Ani"})
	if got := appErrorStatus(t, err); got != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, got)
	}
	if memberRepo.created != nil || memberRepo.restored != nil {
		t.Fatal("no member should be created when access is denied")
	}
}
