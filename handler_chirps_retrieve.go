package main

import "net/http"

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	// get all the chirps from the database
	// []database.Chirp
	dbChirps, err := cfg.db.RetrieveChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Couldn't retrieve chirps")
		return
	}

	// convert []database.Chirp into []Chirp (struct in handler_chirps_create.go)
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	// respond
	respondWithJSON(w, 200, chirps)
}
