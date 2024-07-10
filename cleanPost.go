package main

import (
	"strings"
)

var badWords = [3]string{"kerfuffle", "sharbert", "fornax"}

func cleanPost(post string) string {
	words := strings.Split(post, " ")

	for idx := 0; idx < len(words); idx++ {
		for _, badWord := range badWords {
			if strings.ToLower(words[idx]) == badWord {
				words[idx] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}
