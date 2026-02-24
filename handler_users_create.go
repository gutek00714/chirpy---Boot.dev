package main

import (
	"encoding/json"
	"net/http"

	"github.com/gutek00714/chirpy---Boot.dev/internal/auth"
	"github.com/gutek00714/chirpy---Boot.dev/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// decode json with email value
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// hash the password
	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// create user in database giving it email
	// and hashed password
	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: password,
	})
	if err != nil {
		respondWithError(w, 500, "Couldn't create user")
		return
	}

	// respond - convert and send back as json
	respondWithJSON(w, 201, databaseUserToUser(dbUser))
}
