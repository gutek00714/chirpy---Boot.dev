package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gutek00714/chirpy---Boot.dev/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// get token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token provided")
		return
	}

	// validate token
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Cannot validate")
		return
	}

	// get chirp id
	chirp := r.PathValue("chirpID")

	// convert chirp into uuid
	chirp_string, err := uuid.Parse(chirp)
	if err != nil {
		respondWithError(w, 500, "Couldn't convert chirp id")
		return
	}

	// fetch chirpid from database
	chirp_database, err := cfg.db.RetrieveOneChirp(r.Context(), chirp_string)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	// compare caller id with authors id
	if userID != chirp_database.UserID {
		respondWithError(w, 403, "Wrong user")
		return
	}

	// delete the chirp
	err = cfg.db.DeleteChirp(r.Context(), chirp_string)
	if err != nil {
		respondWithError(w, 500, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(204)
}
