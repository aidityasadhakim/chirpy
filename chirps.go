package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aidityasadhakim/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// Handler logic for creating a chirp
	type parameters struct {
		UserID uuid.UUID `json:"user_id"`
		Body   string    `json:"body"`
	}

	parameter := parameters{}

	// Decode the request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&parameter); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding request", err)
		log.Printf("Error decoding request: %v", err)
		return
	}
	// Validate the request
	profane_list := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	for _, word := range strings.Split(parameter.Body, " ") {
		if _, found := profane_list[strings.ToLower(word)]; found {
			parameter.Body = strings.ReplaceAll(parameter.Body, word, "****")
		}
	}

	/*
			f the chirp is valid, you should save it in the database with:
		A new random id: A UUID
		created_at: A non-null timestamp
		updated_at: A non null timestamp
		body: A non-null string
		user_id: This should reference the id of the user who created the chirp, and ON DELETE CASCADE, which will cause a user's chirps to be deleted if the user is deleted.
	*/
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   parameter.Body,
		UserID: parameter.UserID,
	})
	chirpResponse := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		log.Printf("Error creating chirp: %v", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, chirpResponse)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	//  Handler logic for getting all the chirps basically SELECT * FROM chirps

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps", err)
		log.Printf("Error getting chirps: %v", err)
		return
	}

	chirpResponses := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		chirpResponses[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirpResponses)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	// Handler logic for getting a single chirp
	// This should be a GET request to /api/chirps/{id}
	// The id should be passed as a URL parameter
	var id uuid.UUID
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirp", err)
		log.Printf("Error getting chirp: %v", err)
		return
	}

	chirpResponse := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirpResponse)
}
