package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gutek00714/chirpy---Boot.dev/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUserChirpyRed(w http.ResponseWriter, r *http.Request) {
	// check if apikey matches
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaSecret {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			User_id string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding parameters")
		return
	}

	//check if event is user.upgraded
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	// parse user id into uuid
	user_id, err := uuid.Parse(params.Data.User_id)
	if err != nil {
		respondWithError(w, 500, "Couldn't parse user id")
		return
	}

	// upgrade user
	err = cfg.db.UpgradeUserChirpyRed(r.Context(), user_id)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	w.WriteHeader(204)
}
