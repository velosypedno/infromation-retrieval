package normalizers

import (
	"strings"
)

type LowerCaseNormalizer struct{}

func (n LowerCaseNormalizer) Normalize(tokens []string) []string {
	terms := make([]string, len(tokens))

	for i, token := range tokens {
		terms[i] = strings.ToLower(token)
	}
	return terms
}
