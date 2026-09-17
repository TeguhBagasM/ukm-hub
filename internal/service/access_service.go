package service

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/entity"
	"handler/internal/repository"
	"handler/internal/utils"

	"github.com/google/uuid"
)

type AccessService interface {
	CanAccessOrganization(ctx dto.AuthContext, orgID uuid.UUID) (bool, error)
	MustAccessOrganization(ctx dto.AuthContext, orgID uuid.UUID) error
	FindAccessibleOrganizations(ctx dto.AuthContext) ([]entity.Organization, error)
}

type accessService struct {
	orgRepo repository.OrganizationRepository
}

func NewAccessService(orgRepo repository.OrganizationRepository) AccessService {
	return &accessService{orgRepo: orgRepo}
}

func (s *accessService) CanAccessOrganization(ctx dto.AuthContext, orgID uuid.UUID) (bool, error) {
	if ctx.Role == utils.RoleSuperAdmin {
		return true, nil
	}
	if ctx.Role != utils.RoleOrgAdmin {
		return false, nil
	}
	uid, err := uuid.Parse(ctx.UserID)
	if err != nil {
		return false, utils.NewAppError(http.StatusBadRequest, "invalid user id")
	}
	return s.orgRepo.IsAdminOf(uid, orgID)
}

func (s *accessService) MustAccessOrganization(ctx dto.AuthContext, orgID uuid.UUID) error {
	ok, err := s.CanAccessOrganization(ctx, orgID)
	if err != nil {
		return err
	}
	if !ok {
		return utils.NewAppError(http.StatusForbidden, "forbidden: no access to this organization")
	}
	return nil
}

func (s *accessService) FindAccessibleOrganizations(ctx dto.AuthContext) ([]entity.Organization, error) {
	if ctx.Role == utils.RoleSuperAdmin {
		return s.orgRepo.FindAll()
	}
	if ctx.Role != utils.RoleOrgAdmin {
		return nil, utils.NewAppError(http.StatusForbidden, "forbidden: insufficient role")
	}
	uid, err := uuid.Parse(ctx.UserID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid user id")
	}
	return s.orgRepo.FindByAdminUser(uid)
}
