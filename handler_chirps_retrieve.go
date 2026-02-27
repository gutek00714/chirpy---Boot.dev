package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/gutek00714/chirpy---Boot.dev/internal/database"
)

// retrieve all chirps ordered by created_at
func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	// optional author_id
	s := r.URL.Query().Get("author_id")

	var dbChirps []database.Chirp
	var err error
	var authorID uuid.UUID

	if s != "" {
		authorID, err = uuid.Parse(s)
		if err != nil {
			respondWithError(w, 400, "Invalid author_id")
			return
		}

		dbChirps, err = cfg.db.GetUsersChirps(r.Context(), authorID)
	} else {
		// get all the chirps from the database
		// []database.Chirp
		dbChirps, err = cfg.db.RetrieveChirps(r.Context())
	}
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

	// sort by created_at
	sortOrder := r.URL.Query().Get("sort")
	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	// respond
	respondWithJSON(w, 200, chirps)
}

// retrieve one chirp by id
func (cfg *apiConfig) handlerOneChirpRetrieve(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpID_UUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 404, "Invalid ID")
		return
	}

	// get chirp from database
	foundChirp, err := cfg.db.RetrieveOneChirp(r.Context(), chirpID_UUID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	// respond
	respondWithJSON(w, 200, Chirp{
		ID:        foundChirp.ID,
		CreatedAt: foundChirp.CreatedAt,
		UpdatedAt: foundChirp.UpdatedAt,
		Body:      foundChirp.Body,
		UserID:    foundChirp.UserID,
	})
}
