package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gutek00714/chirpy---Boot.dev/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	// get cleaned body from the validateChirp below and respond
	cleaned_body, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	// call the database to create a new chirp record
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: params.User_id,
	})

	// handle possible error from database return
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	// map the database result into chirp struct and respond
	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	// check length
	if len(body) > 140 {
		return "", errors.New("Chirp is too long")
	} else {

		// clean bad words
		cleaned_body := getCleanedBody(body, badWords)
		return cleaned_body, nil
	}
}
