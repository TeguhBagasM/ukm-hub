package dto

type DashboardMetrics struct {
	TotalMembers        int64 `json:"total_members"`
	ActiveMembers       int64 `json:"active_members"`
	ActiveEvents        int64 `json:"active_events"`
	TotalEvents         int64 `json:"total_events"`
	TotalRegistrations  int64 `json:"total_registrations"`
	PendingApplications int64 `json:"pending_applications"`
}

type MembersByDivisionItem struct {
	DivisionID   *string `json:"division_id"`
	DivisionName string  `json:"division_name"`
	Count        int64   `json:"count"`
}

type RecentRegistrationItem struct {
	ID          string `json:"id"`
	EventID     string `json:"event_id"`
	EventTitle  string `json:"event_title"`
	Status      string `json:"status"`
	SubmittedAt string `json:"submitted_at"`
}

type DashboardResponse struct {
	Metrics                        DashboardMetrics         `json:"metrics"`
	RegistrationStatusDistribution map[string]int64         `json:"registration_status_distribution"`
	MembersByDivision              []MembersByDivisionItem  `json:"members_by_division"`
	RecentRegistrations            []RecentRegistrationItem `json:"recent_registrations"`
}
