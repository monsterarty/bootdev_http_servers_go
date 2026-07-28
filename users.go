package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type user struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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
	ctx := r.Context()
	result, err := cfg.db.CreateUser(ctx, params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Email:     result.Email,
	})
}
