package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()

	// create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		// who issued the token
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),

		// who the token is "about"
		Subject: userID.String(),
	})

	// turn the token into the final JWT string and sign it
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	// return the final string the client will store or send back on requests
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// claims = data that is expected to be inside the token
	claims := &jwt.RegisteredClaims{}

	// split token string into header/claims/signature, decode claims into claims struct,
	// use keyFunc to get the secret key, verify if signature matches
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		//check if token is using HMAC (HS256/384/512)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// must return the same key type that was used when the key was signed
		return []byte(tokenSecret), nil
	})
	// if signature is wrong or token is expired - error and reject it
	if err != nil {
		return uuid.Nil, err
	}

	// user id is stored in Subject -> MakeJWT func
	idStr := claims.Subject
	// convert it from string to uuid
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// reference:
// https://github.com/golang-jwt/jwt/blob/main/example_test.go
