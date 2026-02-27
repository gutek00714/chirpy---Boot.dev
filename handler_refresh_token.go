package main

import (
	"net/http"
	"time"

	"github.com/gutek00714/chirpy---Boot.dev/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	// extract the token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token provided")
		return
	}

	// find the user by refresh token
	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Invalid or expired refresh token")
		return
	}

	// create a new access token
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Couldn't create token")
		return
	}

	// respond
	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, response{Token: accessToken})
}
