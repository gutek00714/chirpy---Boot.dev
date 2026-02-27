package main

import (
	"encoding/json"
	"net/http"

	"github.com/gutek00714/chirpy---Boot.dev/internal/auth"
	"github.com/gutek00714/chirpy---Boot.dev/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	// get token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token provided")
		return
	}

	// get email and password
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding parameters")
		return
	}

	// hash the new password
	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// get user information
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Cannot validate")
		return
	}

	// update user in database
	dbUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: password,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, 500, "Couldn't update user")
		return
	}

	respondWithJSON(w, 200, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		ChirpyRed: dbUser.IsChirpyRed,
	})
}
