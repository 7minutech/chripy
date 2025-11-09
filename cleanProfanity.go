package main

import (
	"strings"
)

var profainWords = map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
var censor string = "****"

func cleanProfanity(resp string) string {
	words := strings.Split(resp, " ")
	for i, word := range words {
		_, ok := profainWords[strings.ToLower(word)]
		if ok {
			words[i] = censor
		}
	}
	return strings.Join(words, " ")
}
