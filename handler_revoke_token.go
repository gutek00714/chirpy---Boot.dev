package main

import (
	"net/http"

	"github.com/gutek00714/chirpy---Boot.dev/internal/auth"
)

func (cfg *apiConfig) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	// get the refresh token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token provided")
		return
	}

	// revoke token
	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, "Couldn't revoke token")
		return
	}

	w.WriteHeader(204)
}
