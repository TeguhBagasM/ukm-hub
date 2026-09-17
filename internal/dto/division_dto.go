package dto

type CreateDivisionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=150"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateDivisionRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=150"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type DivisionResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
