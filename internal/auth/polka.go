package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// same logic as GetBearerToken
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Header doesn't exist")
	}

	if strings.HasPrefix(authHeader, "ApiKey") {
		token := strings.TrimPrefix(authHeader, "ApiKey")
		token = strings.TrimSpace(token)

		if token == "" {
			return "", fmt.Errorf("ApiKey doesn't exist")
		}

		return token, nil
	} else {
		return "", fmt.Errorf("ApiKey doesn't exist")
	}
}
