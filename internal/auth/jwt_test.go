package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT_ReturnsToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("MakeJWT returned empty token")
	}
}

func TestValidateJWT_Valid(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// create a token and validate it
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// validate should return the same user ID
	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if gotID != userID {
		t.Fatalf("ValidateJWT returned wrong ID: got %v, want %v", gotID, userID)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// validating with wrong secret should fail
	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("ValidateJWT should have failed with wrong secret, but didn't")
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// create a token that expires immediately (negative duration = already expired)
	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// validating an expired token should fail
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("ValidateJWT should have failed with expired token, but didn't")
	}
}
