package service

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/entity"
	"handler/internal/repository"
	"handler/internal/utils"

	"github.com/google/uuid"
)

type DivisionService interface {
	ListByOrganization(ctx dto.AuthContext, orgIDStr string) ([]dto.DivisionResponse, error)
	Get(ctx dto.AuthContext, idStr string) (*dto.DivisionResponse, error)
	Create(ctx dto.AuthContext, orgIDStr string, req dto.CreateDivisionRequest) (*dto.DivisionResponse, error)
	Update(ctx dto.AuthContext, idStr string, req dto.UpdateDivisionRequest) (*dto.DivisionResponse, error)
	Delete(ctx dto.AuthContext, idStr string) error
}

type divisionService struct {
	divRepo repository.DivisionRepository
	access  AccessService
}

func NewDivisionService(divRepo repository.DivisionRepository, access AccessService) DivisionService {
	return &divisionService{divRepo: divRepo, access: access}
}

func (s *divisionService) ListByOrganization(ctx dto.AuthContext, orgIDStr string) ([]dto.DivisionResponse, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if err := s.access.MustAccessOrganization(ctx, orgID); err != nil {
		return nil, err
	}
	divs, err := s.divRepo.FindByOrganization(orgID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.DivisionResponse, 0, len(divs))
	for i := range divs {
		res = append(res, *s.toResponse(&divs[i]))
	}
	return res, nil
}

func (s *divisionService) Get(ctx dto.AuthContext, idStr string) (*dto.DivisionResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid division id")
	}
	div, err := s.divRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "division not found")
	}
	if err := s.access.MustAccessOrganization(ctx, div.OrganizationID); err != nil {
		return nil, err
	}
	return s.toResponse(div), nil
}

func (s *divisionService) Create(ctx dto.AuthContext, orgIDStr string, req dto.CreateDivisionRequest) (*dto.DivisionResponse, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if err := s.access.MustAccessOrganization(ctx, orgID); err != nil {
		return nil, err
	}
	status := req.Status
	if status == "" {
		status = entity.DivisionStatusActive
	}
	div := entity.Division{
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		Status:         status,
	}
	if err := s.divRepo.Create(&div); err != nil {
		return nil, err
	}
	return s.toResponse(&div), nil
}

func (s *divisionService) Update(ctx dto.AuthContext, idStr string, req dto.UpdateDivisionRequest) (*dto.DivisionResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid division id")
	}
	div, err := s.divRepo.FindByID(id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "division not found")
	}
	if err := s.access.MustAccessOrganization(ctx, div.OrganizationID); err != nil {
		return nil, err
	}
	if req.Name != "" {
		div.Name = req.Name
	}
	if req.Description != "" {
		div.Description = req.Description
	}
	if req.Status != "" {
		div.Status = req.Status
	}
	if err := s.divRepo.Update(div); err != nil {
		return nil, err
	}
	return s.toResponse(div), nil
}

func (s *divisionService) Delete(ctx dto.AuthContext, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.NewAppError(http.StatusBadRequest, "invalid division id")
	}
	div, err := s.divRepo.FindByID(id)
	if err != nil {
		return utils.NewAppError(http.StatusNotFound, "division not found")
	}
	if err := s.access.MustAccessOrganization(ctx, div.OrganizationID); err != nil {
		return err
	}
	return s.divRepo.Delete(id)
}

func (s *divisionService) toResponse(div *entity.Division) *dto.DivisionResponse {
	return &dto.DivisionResponse{
		ID:             div.ID.String(),
		OrganizationID: div.OrganizationID.String(),
		Name:           div.Name,
		Description:    div.Description,
		Status:         div.Status,
		CreatedAt:      div.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      div.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
