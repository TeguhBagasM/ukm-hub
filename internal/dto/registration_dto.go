package dto

type RejectRegistrationRequest struct {
	Reason string `json:"reason" binding:"required,min=3,max=500"`
}

type RegistrationAnswerResponse struct {
	FieldID string `json:"field_id"`
	Label   string `json:"label"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

type RegistrationResponse struct {
	ID              string                       `json:"id"`
	EventID         string                       `json:"event_id"`
	Status          string                       `json:"status"`
	RejectionReason string                       `json:"rejection_reason"`
	SubmittedAt     string                       `json:"submitted_at"`
	CreatedAt       string                       `json:"created_at"`
	UpdatedAt       string                       `json:"updated_at"`
	Answers         []RegistrationAnswerResponse `json:"answers"`
}

type RegistrationListResponse struct {
	Items      []RegistrationResponse `json:"items"`
	Page       int                    `json:"page"`
	PerPage    int                    `json:"per_page"`
	Total      int64                  `json:"total"`
	TotalPages int                    `json:"total_pages"`
}
