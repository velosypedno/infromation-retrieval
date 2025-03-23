package internal

import (
	"strings"
	"unicode"
)

type DefaultTokenizer struct{}

func (t DefaultTokenizer) Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '\'' && r != '-'
	})
}
