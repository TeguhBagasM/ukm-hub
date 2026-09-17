package dto

type CreateFormRequest struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	SubmitWebhookURL string `json:"submit_webhook_url"`
	IsPublished      bool   `json:"is_published"`
}

type UpdateFormRequest struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	SubmitWebhookURL string `json:"submit_webhook_url"`
}

type FormFieldRequest struct {
	Label       string   `json:"label" binding:"required,min=1,max=150"`
	Name        string   `json:"name" binding:"required,min=1,max=100"`
	Type        string   `json:"type" binding:"required"`
	Placeholder string   `json:"placeholder"`
	Description string   `json:"description"`
	Required    *bool    `json:"required"`
	Options     []string `json:"options"`
	SortOrder   int      `json:"sort_order"`
}

type ReorderFieldRequest struct {
	ID        string `json:"id" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

type FormFieldResponse struct {
	ID          string   `json:"id"`
	FormID      string   `json:"form_id"`
	Label       string   `json:"label"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Placeholder string   `json:"placeholder"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Options     []string `json:"options"`
	SortOrder   int      `json:"sort_order"`
}

type FormResponse struct {
	ID               string              `json:"id"`
	EventID          string              `json:"event_id"`
	Title            string              `json:"title"`
	Description      string              `json:"description"`
	SubmitWebhookURL string              `json:"submit_webhook_url"`
	IsPublished      bool                `json:"is_published"`
	Fields           []FormFieldResponse `json:"fields"`
	CreatedAt        string              `json:"created_at"`
	UpdatedAt        string              `json:"updated_at"`
}
