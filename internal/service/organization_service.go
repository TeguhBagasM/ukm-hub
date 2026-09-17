package service

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/entity"
	"handler/internal/repository"
	"handler/internal/utils"

	"github.com/google/uuid"
)

type OrganizationService interface {
	List(ctx dto.AuthContext) ([]dto.OrganizationResponse, error)
	Get(ctx dto.AuthContext, idStr string) (*dto.OrganizationResponse, error)
	Create(ctx dto.AuthContext, req dto.CreateOrganizationRequest) (*dto.OrganizationResponse, error)
	Update(ctx dto.AuthContext, idStr string, req dto.UpdateOrganizationRequest) (*dto.OrganizationResponse, error)
	Delete(ctx dto.AuthContext, idStr string) error
}

type organizationService struct {
	orgRepo repository.OrganizationRepository
	access  AccessService
}

func NewOrganizationService(orgRepo repository.OrganizationRepository, access AccessService) OrganizationService {
	return &organizationService{orgRepo: orgRepo, access: access}
}

func (s *organizationService) List(ctx dto.AuthContext) ([]dto.OrganizationResponse, error) {
	orgs, err := s.access.FindAccessibleOrganizations(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.OrganizationResponse, 0, len(orgs))
	for i := range orgs {
		res = append(res, *s.toResponse(&orgs[i]))
	}
	return res, nil
}

func (s *organizationService) Get(ctx dto.AuthContext, idStr string) (*dto.OrganizationResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	org, err := s.orgRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "organization not found")
	}
	if err := s.access.MustAccessOrganization(ctx, id); err != nil {
		return nil, err
	}
	return s.toResponse(org), nil
}

func (s *organizationService) Create(ctx dto.AuthContext, req dto.CreateOrganizationRequest) (*dto.OrganizationResponse, error) {
	if ctx.Role != utils.RoleSuperAdmin {
		return nil, utils.NewAppError(http.StatusForbidden, "forbidden: only super admin can create organizations")
	}
	if existing, _ := s.orgRepo.FindBySlug(req.Slug); existing != nil {
		return nil, utils.NewAppError(http.StatusConflict, "slug already used")
	}
	status := req.Status
	if status == "" {
		status = entity.OrganizationStatusActive
	}
	org := entity.Organization{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Logo:        req.Logo,
		Email:       req.Email,
		Phone:       req.Phone,
		Status:      status,
	}
	if err := s.orgRepo.Create(&org); err != nil {
		return nil, err
	}
	return s.toResponse(&org), nil
}

func (s *organizationService) Update(ctx dto.AuthContext, idStr string, req dto.UpdateOrganizationRequest) (*dto.OrganizationResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	org, err := s.orgRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "organization not found")
	}
	if err := s.access.MustAccessOrganization(ctx, id); err != nil {
		return nil, err
	}

	isSuper := ctx.Role == utils.RoleSuperAdmin

	if req.Name != "" {
		org.Name = req.Name
	}
	if req.Description != "" {
		org.Description = req.Description
	}
	if req.Logo != "" {
		org.Logo = req.Logo
	}
	if req.Email != "" {
		org.Email = req.Email
	}
	if req.Phone != "" {
		org.Phone = req.Phone
	}
	if req.Slug != "" {
		if !isSuper {
			return nil, utils.NewAppError(http.StatusForbidden, "forbidden: only super admin can change slug")
		}
		if other, _ := s.orgRepo.FindBySlug(req.Slug); other != nil && other.ID != id {
			return nil, utils.NewAppError(http.StatusConflict, "slug already used")
		}
		org.Slug = req.Slug
	}
	if req.Status != "" {
		if !isSuper {
			return nil, utils.NewAppError(http.StatusForbidden, "forbidden: only super admin can change status")
		}
		org.Status = req.Status
	}

	if err := s.orgRepo.Update(org); err != nil {
		return nil, err
	}
	return s.toResponse(org), nil
}

func (s *organizationService) Delete(ctx dto.AuthContext, idStr string) error {
	if ctx.Role != utils.RoleSuperAdmin {
		return utils.NewAppError(http.StatusForbidden, "forbidden: only super admin can delete organizations")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if _, err := s.orgRepo.FindByID(id); err != nil {
		return utils.NewAppError(http.StatusNotFound, "organization not found")
	}
	return s.orgRepo.Delete(id)
}

func (s *organizationService) toResponse(org *entity.Organization) *dto.OrganizationResponse {
	return &dto.OrganizationResponse{
		ID:          org.ID.String(),
		Name:        org.Name,
		Slug:        org.Slug,
		Description: org.Description,
		Logo:        org.Logo,
		Email:       org.Email,
		Phone:       org.Phone,
		Status:      org.Status,
		CreatedAt:   org.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   org.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
