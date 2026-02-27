package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gutek00714/chirpy---Boot.dev/internal/auth"
	"github.com/gutek00714/chirpy---Boot.dev/internal/database"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`

		//optional
		// Expires_in_seconds int `json:"expires_in_seconds"`
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

	// check if expiration time is specified by the client
	// defaultExpiration := time.Hour
	// if params.Expires_in_seconds > 0 {
	// 	requestedExpiration := time.Duration(params.Expires_in_seconds) * time.Second

	// 	if requestedExpiration < defaultExpiration {
	// 		defaultExpiration = requestedExpiration
	// 	}
	// }

	// create JWT
	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Couldn't create JWT")
		return
	}

	// generate refresh token
	refreshToken := auth.MakeRefreshToken()
	now := time.Now().UTC()
	expires := now.Add(60 * 24 * time.Hour) // 60 days

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: expires,
	})
	if err != nil {
		respondWithError(w, 500, "Couldn't create refresh token")
		return
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// if the password matches return ok
	respondWithJSON(w, 200, response{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
			ChirpyRed: dbUser.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})

}
