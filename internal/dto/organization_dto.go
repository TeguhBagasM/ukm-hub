package dto

type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=150"`
	Slug        string `json:"slug" binding:"required,min=3,max=150"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateOrganizationRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=150"`
	Slug        string `json:"slug" binding:"omitempty,min=3,max=150"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type OrganizationResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
