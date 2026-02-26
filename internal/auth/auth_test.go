package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "correcthorsebatterystaple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// hash should not be empty
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}

	// hash should not be the same as the password
	if hash == password {
		t.Fatal("HashPassword returned the plain password")
	}
}

func TestCheckPasswordHash_Correct(t *testing.T) {
	password := "correcthorsebatterystaple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// correct password should match
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed: %v", err)
	}
	if !match {
		t.Fatal("CheckPasswordHash returned false for correct password")
	}
}

func TestCheckPasswordHash_Wrong(t *testing.T) {
	password := "correcthorsebatterystaple"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// wrong password should not match
	match, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed: %v", err)
	}
	if match {
		t.Fatal("CheckPasswordHash returned true for wrong password")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "samepassword"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// same password should produce different hashes (because of random salt)
	if hash1 == hash2 {
		t.Fatal("HashPassword produced identical hashes for the same password — salt may not be working")
	}
}
