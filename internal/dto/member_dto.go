package dto

type ConvertMemberRequest struct {
	DivisionID string `json:"division_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	StudentID  string `json:"student_id"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateMemberRequest struct {
	DivisionID string `json:"division_id"`
	Name       string `json:"name" binding:"omitempty,min=2,max=150"`
	Email      string `json:"email" binding:"omitempty,email"`
	Phone      string `json:"phone"`
	StudentID  string `json:"student_id"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type MemberResponse struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"organization_id"`
	DivisionID     *string `json:"division_id"`
	RegistrationID *string `json:"registration_id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	StudentID      string  `json:"student_id"`
	Status         string  `json:"status"`
	JoinedAt       string  `json:"joined_at"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type MemberListResponse struct {
	Items      []MemberResponse `json:"items"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	Total      int64            `json:"total"`
	TotalPages int              `json:"total_pages"`
}
