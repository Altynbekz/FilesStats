package main

import (
	"strings"
	"unicode"
)

func tokenize(text string) []string {

	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	wholeText := []string{}
	for _, word := range words {
		word := strings.Map(unicode.ToLower, word)
		wholeText = append(wholeText, word)
	}

	return wholeText
}
