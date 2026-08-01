package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/monsterarty/bootdev_http_servers_go/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode body", err)
		return
	}

	// Constants and configs
	const maxChirpLength = 140
	const minChirpLength = 0

	if len(params.Body) > maxChirpLength {
		log.Printf("Chirp is too long")
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	if len(params.Body) < minChirpLength {
		log.Printf("Empty Chirp")
		respondWithError(w, http.StatusBadRequest, "Chirp is empty", nil)
		return
	}
	if len(params.UserID) < 1 {
		log.Printf("Empty User_ID")
		respondWithError(w, http.StatusBadRequest, "User_ID is empty", nil)
		return
	}

	ctx := r.Context()
	result, err := cfg.db.CreateChirp(ctx, database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coundnt create Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body:      result.Body,
		UserID:    result.UserID,
	})
}

func profaneFilter(body string, badWords map[string]struct{}) string {
	splitted := strings.Split(body, " ")

	for i, word := range splitted {
		lowered := strings.ToLower(word)
		if _, ok := badWords[lowered]; ok {
			splitted[i] = "****"
		}
	}
	joined := strings.Join(splitted, " ")
	return joined

}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	chirps, err := cfg.db.GetAllChirps(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coundnt fetch Chirps", err)
		return
	}

	sliceChirps := make([]Chirp, 0, len(chirps))

	for _, chirp := range chirps {
		sliceChirps = append(sliceChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, sliceChirps)

}

func (cfg *apiConfig) handlerGetOneChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if len(chirpID) == 0 {
		log.Println("Empty Chirp ID")
		respondWithError(w, http.StatusBadRequest, "Empty Chirp ID", nil)
		return
	}
	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Not an Chipr ID", err)
		return
	}
	ctx := r.Context()
	chirp, err := cfg.db.GetOneChirp(ctx, parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not found Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
