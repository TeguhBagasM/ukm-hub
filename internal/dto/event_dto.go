package dto

import "time"

type CreateEventRequest struct {
	Name              string     `json:"name" binding:"required,min=3,max=150"`
	Slug              string     `json:"slug" binding:"required,min=3,max=150"`
	Description       string     `json:"description"`
	Location          string     `json:"location"`
	StartDate         *time.Time `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
	RegistrationStart *time.Time `json:"registration_start"`
	RegistrationEnd   *time.Time `json:"registration_end"`
	Quota             int        `json:"quota"`
	Status            string     `json:"status" binding:"omitempty,oneof=DRAFT PUBLISHED CLOSED ARCHIVED"`
}

type UpdateEventRequest struct {
	Name              string     `json:"name" binding:"omitempty,min=3,max=150"`
	Slug              string     `json:"slug" binding:"omitempty,min=3,max=150"`
	Description       string     `json:"description"`
	Location          string     `json:"location"`
	StartDate         *time.Time `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
	RegistrationStart *time.Time `json:"registration_start"`
	RegistrationEnd   *time.Time `json:"registration_end"`
	Quota             *int       `json:"quota"`
	Status            string     `json:"status" binding:"omitempty,oneof=DRAFT PUBLISHED CLOSED ARCHIVED"`
}

type EventResponse struct {
	ID                string     `json:"id"`
	OrganizationID    string     `json:"organization_id"`
	Name              string     `json:"name"`
	Slug              string     `json:"slug"`
	Description       string     `json:"description"`
	Location          string     `json:"location"`
	StartDate         *time.Time `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
	RegistrationStart *time.Time `json:"registration_start"`
	RegistrationEnd   *time.Time `json:"registration_end"`
	Quota             int        `json:"quota"`
	Status            string     `json:"status"`
	RegistrationCount int        `json:"registration_count"`
	CreatedAt         string     `json:"created_at"`
	UpdatedAt         string     `json:"updated_at"`
}
