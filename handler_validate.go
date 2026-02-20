package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

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
		body := getCleanedBody(params.Body, badWords)
		respondWithJSON(w, 200, validResponse{CleanedBody: body})
		// respondWithJSON(w, 200, validResponse{Valid: true})
	}
}

type validResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

// check it there are restricted words
func getCleanedBody(body string, badWords map[string]struct{}) string {
	// split the body into separate words
	split_body := strings.Split(body, " ")

	// loop through every word to check if it's in the badWords map
	for i, word := range split_body {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			split_body[i] = "****"
		}
	}
	body = strings.Join(split_body, " ")
	return body
}
