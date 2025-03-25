package main

import (
	"strings"
	"unicode"
)

type Tokenizer interface {
	Tokenize(text string) []string
}

type DefaultTokenizer struct{}

func (t DefaultTokenizer) Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '\'' && r != '-'
	})
}
