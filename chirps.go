package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldnt decode body", err)
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp is too long")
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	if len(params.Body) <= 0 {
		log.Printf("Empty chirp")
		respondWithError(w, 400, "Chirp is empty", nil)
		return
	}

	respondWithJSON(w, 200, returnVals{Valid: true})

}
