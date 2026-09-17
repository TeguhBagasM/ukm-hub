package service

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"time"

	"ukm-hub/internal/dto"
	"ukm-hub/internal/entity"
	"ukm-hub/internal/repository"
	"ukm-hub/internal/utils"

	"github.com/google/uuid"
)

type DashboardService interface {
	GetDashboard(ctx dto.AuthContext, orgIDStr string) (*dto.DashboardResponse, error)
	ExportRegistrations(ctx dto.AuthContext, orgIDStr, eventIDStr string) (string, []byte, error)
	ExportMembers(ctx dto.AuthContext, orgIDStr string) (string, []byte, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
	memberRepo    repository.MemberRepository
	regRepo       repository.RegistrationRepository
	formRepo      repository.FormRepository
	fieldRepo     repository.FormFieldRepository
	divRepo       repository.DivisionRepository
	eventRepo     repository.EventRepository
	access        AccessService
}

func NewDashboardService(
	dashboardRepo repository.DashboardRepository,
	memberRepo repository.MemberRepository,
	regRepo repository.RegistrationRepository,
	formRepo repository.FormRepository,
	fieldRepo repository.FormFieldRepository,
	divRepo repository.DivisionRepository,
	eventRepo repository.EventRepository,
	access AccessService,
) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
		memberRepo:    memberRepo,
		regRepo:       regRepo,
		formRepo:      formRepo,
		fieldRepo:     fieldRepo,
		divRepo:       divRepo,
		eventRepo:     eventRepo,
		access:        access,
	}
}

func (s *dashboardService) GetDashboard(ctx dto.AuthContext, orgIDStr string) (*dto.DashboardResponse, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if err := s.access.MustAccessOrganization(ctx, orgID); err != nil {
		return nil, err
	}

	totalMembers, activeMembers, err := s.dashboardRepo.CountMembers(orgID)
	if err != nil {
		return nil, err
	}
	totalEvents, activeEvents, err := s.dashboardRepo.CountEvents(orgID)
	if err != nil {
		return nil, err
	}
	totalRegs, pendingRegs, err := s.dashboardRepo.CountRegistrations(orgID)
	if err != nil {
		return nil, err
	}
	statusCounts, err := s.dashboardRepo.RegistrationStatusCounts(orgID)
	if err != nil {
		return nil, err
	}
	byDivision, err := s.dashboardRepo.MembersByDivision(orgID)
	if err != nil {
		return nil, err
	}
	recent, err := s.dashboardRepo.RecentRegistrations(orgID, 5)
	if err != nil {
		return nil, err
	}

	res := &dto.DashboardResponse{
		Metrics: dto.DashboardMetrics{
			TotalMembers:        totalMembers,
			ActiveMembers:       activeMembers,
			ActiveEvents:        activeEvents,
			TotalEvents:         totalEvents,
			TotalRegistrations:  totalRegs,
			PendingApplications: pendingRegs,
		},
		RegistrationStatusDistribution: statusCounts,
		MembersByDivision:              make([]dto.MembersByDivisionItem, 0, len(byDivision)),
		RecentRegistrations:            make([]dto.RecentRegistrationItem, 0, len(recent)),
	}
	for _, item := range byDivision {
		divItem := dto.MembersByDivisionItem{DivisionName: item.DivisionName, Count: item.Count}
		if item.DivisionID != nil {
			v := item.DivisionID.String()
			divItem.DivisionID = &v
		}
		res.MembersByDivision = append(res.MembersByDivision, divItem)
	}
	for _, item := range recent {
		res.RecentRegistrations = append(res.RecentRegistrations, dto.RecentRegistrationItem{
			ID:          item.ID.String(),
			EventID:     item.EventID.String(),
			EventTitle:  item.EventTitle,
			Status:      item.Status,
			SubmittedAt: item.SubmittedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}

func (s *dashboardService) ExportRegistrations(ctx dto.AuthContext, orgIDStr, eventIDStr string) (string, []byte, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return "", nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if err := s.access.MustAccessOrganization(ctx, orgID); err != nil {
		return "", nil, err
	}

	var eventID *uuid.UUID
	if eventIDStr != "" {
		parsed, err := uuid.Parse(eventIDStr)
		if err != nil {
			return "", nil, utils.NewAppError(http.StatusBadRequest, "invalid event id")
		}
		event, err := s.eventRepo.FindByID(parsed)
		if err != nil || event.OrganizationID != orgID {
			return "", nil, utils.NewAppError(http.StatusBadRequest, "event does not belong to this organization")
		}
		eventID = &parsed
	}

	rows, err := s.dashboardRepo.ListRegistrationExportRows(orgID, eventID)
	if err != nil {
		return "", nil, err
	}

	regIDs := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		regIDs = append(regIDs, row.ID)
	}
	answers, err := s.regRepo.FindAnswersByRegistrationIDs(regIDs)
	if err != nil {
		return "", nil, err
	}
	answerMap := make(map[uuid.UUID]map[uuid.UUID]string)
	for _, a := range answers {
		if answerMap[a.RegistrationID] == nil {
			answerMap[a.RegistrationID] = make(map[uuid.UUID]string)
		}
		answerMap[a.RegistrationID][a.FormFieldID] = a.Value
	}

	orderedFields, err := s.collectExportFields(rows)
	if err != nil {
		return "", nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	header := []string{"Event"}
	for _, f := range orderedFields {
		header = append(header, f.label)
	}
	header = append(header, "Status", "Submitted At", "Rejection Reason")
	if err := writer.Write(header); err != nil {
		return "", nil, err
	}
	for _, row := range rows {
		record := []string{row.EventTitle}
		for _, f := range orderedFields {
			record = append(record, answerMap[row.ID][f.id])
		}
		record = append(record,
			row.Status,
			row.SubmittedAt.Format("2006-01-02 15:04:05"),
			row.RejectionReason,
		)
		if err := writer.Write(record); err != nil {
			return "", nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", nil, err
	}

	filename := "registrations-" + orgID.String() + "-" + time.Now().Format("20060102-150405") + ".csv"
	return filename, buf.Bytes(), nil
}

type exportField struct {
	id    uuid.UUID
	label string
}

func (s *dashboardService) collectExportFields(rows []repository.RegistrationExportRow) ([]exportField, error) {
	ordered := make([]exportField, 0)
	seen := make(map[uuid.UUID]bool)

	eventCache := make(map[uuid.UUID][]entity.FormField)
	for _, row := range rows {
		if _, ok := eventCache[row.EventID]; !ok {
			fields, err := s.fieldsForEvent(row.EventID)
			if err != nil {
				return nil, err
			}
			eventCache[row.EventID] = fields
		}
		for _, f := range eventCache[row.EventID] {
			if seen[f.ID] {
				continue
			}
			seen[f.ID] = true
			ordered = append(ordered, exportField{id: f.ID, label: f.Label})
		}
	}
	return ordered, nil
}

func (s *dashboardService) fieldsForEvent(eventID uuid.UUID) ([]entity.FormField, error) {
	form, err := s.formRepo.FindByEvent(eventID)
	if err != nil {
		return nil, nil
	}
	fields, err := s.fieldRepo.FindByForm(form.ID)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func (s *dashboardService) ExportMembers(ctx dto.AuthContext, orgIDStr string) (string, []byte, error) {
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return "", nil, utils.NewAppError(http.StatusBadRequest, "invalid organization id")
	}
	if err := s.access.MustAccessOrganization(ctx, orgID); err != nil {
		return "", nil, err
	}

	members, _, err := s.memberRepo.ListByOrganization(orgID, "", "", 1, 1000000)
	if err != nil {
		return "", nil, err
	}
	divisions, err := s.divRepo.FindByOrganization(orgID)
	if err != nil {
		return "", nil, err
	}
	divisionNames := make(map[uuid.UUID]string, len(divisions))
	for i := range divisions {
		divisionNames[divisions[i].ID] = divisions[i].Name
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"Name", "Email", "Phone", "Student ID", "Division", "Status", "Joined At"}); err != nil {
		return "", nil, err
	}
	for i := range members {
		divisionName := ""
		if members[i].DivisionID != nil {
			divisionName = divisionNames[*members[i].DivisionID]
		}
		record := []string{
			members[i].Name,
			members[i].Email,
			members[i].Phone,
			members[i].StudentID,
			divisionName,
			members[i].Status,
			members[i].JoinedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return "", nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", nil, err
	}

	filename := "members-" + orgID.String() + "-" + time.Now().Format("20060102-150405") + ".csv"
	return filename, buf.Bytes(), nil
}
