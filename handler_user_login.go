package main

import (
	"encoding/json"
	"net/http"

	"github.com/gutek00714/chirpy---Boot.dev/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//decode json
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// get the user (email and password) from database
	dbUser, err := cfg.db.FindUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	// compare the provided password with hashed password from database
	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	// if the password matcher return ok
	respondWithJSON(w, 200, databaseUserToUser(dbUser))

}
