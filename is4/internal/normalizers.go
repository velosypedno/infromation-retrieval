package internal

import (
	"strings"

	"github.com/reiver/go-porterstemmer"
)

type LowerCaseNormalizer struct{}

func (n LowerCaseNormalizer) Normalize(tokens []string) []string {
	terms := make([]string, len(tokens))

	for i, token := range tokens {
		terms[i] = strings.ToLower(token)
	}
	return terms
}

type StemmingNormalizer struct{}

func (n StemmingNormalizer) Normalize(tokens []string) []string {
	terms := make([]string, len(tokens))
	for i, token := range tokens {
		terms[i] = porterstemmer.StemString(token)
	}
	return terms
}
