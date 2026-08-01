package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/monsterarty/bootdev_http_servers_go/internal/auth"
	"github.com/monsterarty/bootdev_http_servers_go/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldnt decode body", err)
		return
	}
	if len(params.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty", err)
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt hash password", err)
	}
	ctx := r.Context()
	result, err := cfg.db.CreateUser(ctx, database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Email:     result.Email,
	})
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldnt decode body", err)
		return
	}
	if len(params.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty", err)
		return
	}
	if len(params.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "Password cannot be empty", err)
		return
	}
	ctx := r.Context()
	result, err := cfg.db.LoginUser(ctx, params.Email)
	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusInternalServerError, "Internal error", nil)
		return
	}

	passCheck, err := auth.CheckPasswordHash(params.Password, result.HashedPassword)
	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusInternalServerError, "Internal error", nil)
		return
	}
	if passCheck == false {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Email:     result.Email,
	})
}
