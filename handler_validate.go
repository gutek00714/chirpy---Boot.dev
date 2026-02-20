package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateHelperFunction(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	// check if body is 140 or less characters
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		respondWithJSON(w, 200, validResponse{Valid: true})
	}
}

type validResponse struct {
	Valid bool `json:"valid"`
}
