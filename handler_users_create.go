package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	// decode json with email value
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// create user in database giving it email
	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, "Couldn't create user")
		return
	}

	// respond - convert and send back as json
	respondWithJSON(w, 201, databaseUserToUser(dbUser))
}
