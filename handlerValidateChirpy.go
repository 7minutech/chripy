package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	const maxChirpLength int = 140

	var params parameters

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	resp := params.Body

	cleanResp := cleanProfanity(resp)

	respondWithJSON(w, http.StatusOK, returnVal{CleanedBody: cleanResp})

}

var profainWords = []string{"kerfuffle", "sharbert", "fornax"}
var censor string = "****"

func cleanProfanity(resp string) string {
	words := strings.Split(resp, " ")
	for i, word := range words {
		if inProfainWords(word) {
			words[i] = censor
		}
	}
	return strings.Join(words, " ")
}

func inProfainWords(word string) bool {
	for _, profanity := range profainWords {
		if strings.ToLower(word) == profanity {
			return true
		}
	}
	return false
}
