package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	// check if header exists
	if authHeader == "" {
		return "", fmt.Errorf("Header doesn't exist")
	}

	// // check if the header string starts with "Bearer" and trim it
	// if strings.HasPrefix(authHeader, "Bearer ") {
	// 	token := strings.Replace(authHeader, "Bearer ", "", 1)
	// 	return token, nil
	// } else {
	// 	return "", fmt.Errorf("Bearer token doesn't exists")
	// }

	// MORE SECURE VERSION
	// check if the header string starts with "Bearer" and trim it
	if strings.HasPrefix(authHeader, "Bearer") {
		token := strings.TrimPrefix(authHeader, "Bearer")
		token = strings.TrimSpace(token)

		// check if token is empty
		if token == "" {
			return "", fmt.Errorf("Bearer token doesn't exists")
		}

		return token, nil
	} else {
		return "", fmt.Errorf("Bearer token doesn't exists")
	}
}
