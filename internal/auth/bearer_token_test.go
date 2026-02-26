package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken_ReturnsToken(t *testing.T) {
	header := http.Header{
		"Authorization": []string{"Bearer my-token"},
	}

	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("GetBearerToken failed: %v", err)
	}
	if token != "my-token" {
		t.Fatalf("got %v, want %v", token, "my-token")
	}
}

func TestGetBearerToken_NoHeader(t *testing.T) {
	header := http.Header{}
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatal("expected an error for missing header, but got nil")
	}
}

func TestGetBearerToken_NoAuthorizationHeader(t *testing.T) {
	header := http.Header{
		"Test": []string{"Bearer my-token"},
	}
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatal("expected an error for missing header, but got nil")
	}
}

func TestGetBearerToken_NoToken(t *testing.T) {
	header := http.Header{
		"Authorization": []string{"Bearer"},
	}
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatal("expected an error for missing token, but got nil")
	}
}

func TestGetBearerToken_NoBearer(t *testing.T) {
	header := http.Header{
		"Authorization": []string{"my-token"},
	}
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatal("expected an error for missing bearer, but got nil")
	}
}
