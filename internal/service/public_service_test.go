package service

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"handler/internal/entity"
	"handler/internal/utils"
)

func appErrorStatus(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *utils.AppError, got %T: %v", err, err)
	}
	return appErr.Status
}

func TestValidateAnswerValue(t *testing.T) {
	tests := []struct {
		name    string
		field   entity.FormField
		value   string
		wantErr int
	}{
		{name: "required empty", field: entity.FormField{Label: "Nama", Type: entity.FormFieldTypeText, Required: true}, value: "", wantErr: http.StatusBadRequest},
		{name: "optional empty", field: entity.FormField{Label: "Portfolio", Type: entity.FormFieldTypeURL}, value: "", wantErr: 0},
		{name: "valid email", field: entity.FormField{Label: "Email", Type: entity.FormFieldTypeEmail}, value: "user@kampus.ac.id", wantErr: 0},
		{name: "invalid email", field: entity.FormField{Label: "Email", Type: entity.FormFieldTypeEmail}, value: "not-an-email", wantErr: http.StatusBadRequest},
		{name: "valid number", field: entity.FormField{Label: "Umur", Type: entity.FormFieldTypeNumber}, value: "20", wantErr: 0},
		{name: "decimal number", field: entity.FormField{Label: "IPK", Type: entity.FormFieldTypeNumber}, value: "3.75", wantErr: 0},
		{name: "invalid number", field: entity.FormField{Label: "Umur", Type: entity.FormFieldTypeNumber}, value: "dua puluh", wantErr: http.StatusBadRequest},
		{name: "valid url", field: entity.FormField{Label: "Portfolio", Type: entity.FormFieldTypeURL}, value: "https://example.com", wantErr: 0},
		{name: "invalid url scheme", field: entity.FormField{Label: "Portfolio", Type: entity.FormFieldTypeURL}, value: "ftp://example.com", wantErr: http.StatusBadRequest},
		{name: "valid date", field: entity.FormField{Label: "Tanggal", Type: entity.FormFieldTypeDate}, value: "2026-09-17", wantErr: 0},
		{name: "invalid date", field: entity.FormField{Label: "Tanggal", Type: entity.FormFieldTypeDate}, value: "17-09-2026", wantErr: http.StatusBadRequest},
		{name: "valid phone", field: entity.FormField{Label: "WA", Type: entity.FormFieldTypePhone}, value: "0812-3456-7890", wantErr: 0},
		{name: "phone too short", field: entity.FormField{Label: "WA", Type: entity.FormFieldTypePhone}, value: "12345", wantErr: http.StatusBadRequest},
		{name: "valid select option", field: entity.FormField{Label: "Divisi", Type: entity.FormFieldTypeSelect, Options: `["Programming","Multimedia"]`}, value: "Programming", wantErr: 0},
		{name: "invalid select option", field: entity.FormField{Label: "Divisi", Type: entity.FormFieldTypeSelect, Options: `["Programming","Multimedia"]`}, value: "Finance", wantErr: http.StatusBadRequest},
		{name: "valid checkbox", field: entity.FormField{Label: "Minat", Type: entity.FormFieldTypeCheckbox, Options: `["AI","Web"]`}, value: "AI, Web", wantErr: 0},
		{name: "invalid checkbox option", field: entity.FormField{Label: "Minat", Type: entity.FormFieldTypeCheckbox, Options: `["AI","Web"]`}, value: "AI, Game", wantErr: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAnswerValue(&tt.field, tt.value)
			if tt.wantErr == 0 {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if got := appErrorStatus(t, err); got != tt.wantErr {
				t.Fatalf("expected status %d, got %d", tt.wantErr, got)
			}
		})
	}
}

func TestEnsureRegistrationOpen(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name    string
		event   entity.Event
		count   int
		wantErr int
	}{
		{name: "published and open", event: entity.Event{Status: entity.EventStatusPublished}, wantErr: 0},
		{name: "draft event", event: entity.Event{Status: entity.EventStatusDraft}, wantErr: http.StatusBadRequest},
		{name: "closed event", event: entity.Event{Status: entity.EventStatusClosed}, wantErr: http.StatusBadRequest},
		{name: "registration not started", event: entity.Event{Status: entity.EventStatusPublished, RegistrationStart: &future}, wantErr: http.StatusBadRequest},
		{name: "registration ended", event: entity.Event{Status: entity.EventStatusPublished, RegistrationEnd: &past}, wantErr: http.StatusBadRequest},
		{name: "quota full", event: entity.Event{Status: entity.EventStatusPublished, Quota: 5}, count: 5, wantErr: http.StatusConflict},
		{name: "quota still available", event: entity.Event{Status: entity.EventStatusPublished, Quota: 5}, count: 4, wantErr: 0},
		{name: "unlimited quota", event: entity.Event{Status: entity.EventStatusPublished, Quota: 0}, count: 1000, wantErr: 0},
	}

	svc := &publicService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ensureRegistrationOpen(&tt.event, tt.count)
			if tt.wantErr == 0 {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if got := appErrorStatus(t, err); got != tt.wantErr {
				t.Fatalf("expected status %d, got %d", tt.wantErr, got)
			}
		})
	}
}

func TestOnlyDigits(t *testing.T) {
	if got := onlyDigits("+62 812-3456_7890"); got != "6281234567890" {
		t.Errorf("expected 6281234567890, got %s", got)
	}
	if got := onlyDigits("abc"); got != "" {
		t.Errorf("expected empty string, got %s", got)
	}
}

func TestDecodeOptions(t *testing.T) {
	if got := decodeOptions(`["A","B"]`); len(got) != 2 || got[0] != "A" || got[1] != "B" {
		t.Errorf("unexpected decode result: %v", got)
	}
	if got := decodeOptions(""); len(got) != 0 {
		t.Errorf("expected empty options, got %v", got)
	}
	if got := decodeOptions("not-json"); len(got) != 0 {
		t.Errorf("expected empty options for invalid json, got %v", got)
	}
}
