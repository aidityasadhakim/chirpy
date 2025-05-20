package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aidityasadhakim/chirpy/internal/auth"
	"github.com/aidityasadhakim/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hashedPassword, err := auth.HashPassword(parameter.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		log.Printf("Error hashing password: %v", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          parameter.Email,
		HashedPassword: hashedPassword,
	})
	user := User{
		ID:             dbUser.ID,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	user, err := cfg.db.GetUserByEmail(r.Context(), parameter.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		log.Printf("Invalid email or password: %v", err)
		return
	}
	if err := auth.CheckPasswordHash(parameter.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		log.Printf("Invalid email or password: %v", err)
		return
	}

	userResponse := struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	log.Printf("User logged in successfully: %s", userResponse.Email)
	respondWithJSON(w, http.StatusOK, userResponse)
}
