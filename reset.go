package main

import (
	"errors"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Invalid platform", errors.New("Invalid platform"))
		return
	}
	ctx := r.Context()
	err := cfg.db.ResetUsers(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt reset users", err)
		return
	}
	cfg.fileserverHits.Store(0)
	respondWithJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Hits reset to 0 and database reset to initial state.",
	})
}
