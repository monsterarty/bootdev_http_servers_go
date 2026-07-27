package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	if len(params.Body) > maxChirpLength {
		log.Printf("Chirp is too long")
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	if len(params.Body) < minChirpLength {
		log.Printf("Empty chirp")
		respondWithError(w, http.StatusInternalServerError, "Chirp is empty", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: profaneFilter(params.Body, badWords),
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
