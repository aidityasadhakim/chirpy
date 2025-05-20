package main

import (
	"strings"
)

func profaneCleaner(s *string) {
	profane_list := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	for _, word := range strings.Split(*s, " ") {
		if _, found := profane_list[strings.ToLower(word)]; found {
			*s = strings.ReplaceAll(*s, word, "****")
		}
	}
}
