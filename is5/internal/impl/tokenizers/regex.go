package tokenizers

import (
	"regexp"
)

type RegexTokenizer struct{}

var wordRegex = regexp.MustCompile(`(?i)\b[a-zA-Z]+(?:[-'][a-zA-Z]+)*\b`)

func (t RegexTokenizer) Tokenize(text string) []string {
	return wordRegex.FindAllString(text, -1)
}
