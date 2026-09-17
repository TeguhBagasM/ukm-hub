package utils

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("admin123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "admin123" {
		t.Fatal("password hash must not equal the plain password")
	}
	if !CheckPasswordHash("admin123", hash) {
		t.Fatal("CheckPasswordHash rejected the correct password")
	}
	if CheckPasswordHash("wrong-password", hash) {
		t.Fatal("CheckPasswordHash accepted an incorrect password")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	userID := uuid.New()
	token, err := GenerateToken(userID, RoleOrgAdmin)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}
	if claims.UserID != userID.String() {
		t.Errorf("expected user id %s, got %s", userID, claims.UserID)
	}
	if claims.Role != RoleOrgAdmin {
		t.Errorf("expected role %s, got %s", RoleOrgAdmin, claims.Role)
	}
}

func TestValidateTokenRejectsTamperedToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateToken(uuid.New(), RoleSuperAdmin)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	tampered := tamperPayload(t, token)

	if _, err := ValidateToken(tampered); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
	if _, err := ValidateToken("not-a-jwt"); err == nil {
		t.Fatal("expected malformed token to be rejected")
	}
}

func tamperPayload(t *testing.T, token string) string {
	t.Helper()
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 JWT segments, got %d", len(parts))
	}
	payload := []byte(parts[1])
	if payload[0] == 'A' {
		payload[0] = 'B'
	} else {
		payload[0] = 'A'
	}
	return parts[0] + "." + string(payload) + "." + parts[2]
}

func TestValidateTokenRejectsWrongSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "issuer-secret")
	token, err := GenerateToken(uuid.New(), RoleOrgAdmin)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	t.Setenv("JWT_SECRET", "attacker-secret")
	if _, err := ValidateToken(token); err == nil {
		t.Fatal("expected token signed with a different secret to be rejected")
	} else if !strings.Contains(err.Error(), "signature") {
		t.Logf("token rejected as expected with error: %v", err)
	}
}
