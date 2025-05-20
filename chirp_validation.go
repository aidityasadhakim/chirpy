package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type parameters struct {
		Body string `json:"body"`
	}

	type validRespVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// Validate the chirp
	req := json.NewDecoder(r.Body)
	params := parameters{}
	if err := req.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		log.Printf("Error decoding request: %v", err)
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp too long: %s", params.Body)
		respondWithError(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	profaneCleaner(&params.Body)

	respondWithJSON(w, http.StatusOK, validRespVals{
		CleanedBody: params.Body,
	})
}
