package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type parameters struct {
		Email string `json:"email"`
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

	dbUser, err := cfg.db.CreateUser(r.Context(), parameter.Email)
	user := User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		log.Printf("Error creating user: %v", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerDeleteAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if os.Getenv("PLATFORM") != "env" {
		respondWithError(w, http.StatusForbidden, "This endpoint is not available in production", nil)
		return
	}

	if err := cfg.db.DeleteAllUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users", err)
		log.Printf("Error deleting users: %v", err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "All users deleted"})
}
