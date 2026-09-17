package dto

import "time"

type RegistrationAnswerRequest struct {
	FieldID string `json:"field_id" binding:"required"`
	Value   string `json:"value"`
}

type PublicRegistrationRequest struct {
	Answers []RegistrationAnswerRequest `json:"answers" binding:"required,min=1,dive"`
}

type PublicOrganizationInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PublicEventResponse struct {
	ID                 string                 `json:"id"`
	Organization       PublicOrganizationInfo `json:"organization"`
	Name               string                 `json:"name"`
	Slug               string                 `json:"slug"`
	Description        string                 `json:"description"`
	Location           string                 `json:"location"`
	StartDate          *time.Time             `json:"start_date"`
	EndDate            *time.Time             `json:"end_date"`
	RegistrationStart  *time.Time             `json:"registration_start"`
	RegistrationEnd    *time.Time             `json:"registration_end"`
	Quota              int                    `json:"quota"`
	RegistrationCount  int                    `json:"registration_count"`
	IsRegistrationOpen bool                   `json:"is_registration_open"`
}

type PublicFormInfo struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Fields      []FormFieldResponse `json:"fields"`
}

type PublicFormResponse struct {
	Event PublicEventResponse `json:"event"`
	Form  PublicFormInfo      `json:"form"`
}

type PublicRegistrationResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	SubmittedAt string `json:"submitted_at"`
}

type WebhookAnswer struct {
	FieldID string `json:"field_id"`
	Label   string `json:"label"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

type RegistrationWebhookPayload struct {
	Event          PublicEventResponse `json:"event"`
	RegistrationID string              `json:"registration_id"`
	Status         string              `json:"status"`
	SubmittedAt    string              `json:"submitted_at"`
	Answers        []WebhookAnswer     `json:"answers"`
	AnswersMap     map[string]string   `json:"answers_map"`
}
