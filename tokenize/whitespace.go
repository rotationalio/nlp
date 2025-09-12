package tokenize

import (
	"strings"
)

type WhitespaceTokenizer struct{}

// Ensure [WhitespaceTokenizer] meets the [Tokenizer] interface requirements.
var _ Tokenizer = &WhitespaceTokenizer{}

// Returns a new [WhitespaceTokenizer].
func NewWhitespaceTokenizer() *WhitespaceTokenizer {
	return &WhitespaceTokenizer{}
}

// Tokenize the chunk of text using whitespace. Splits the string using the
// [strings.Fields] function. ALWAYS returns nil for the error.
//
// Example:
//
//	"The quick brown fox jumped over the quicker-- 105.4% quicker, in fact--
//	\n
//	brown fox because it owed the quicker fox $3.14!"
//
// Tokens:
//
//	{"The", "quick", "brown", "fox", "jumped", "over", "the", "quicker--",
//	"105.4%", "quicker,", "in", "fact--", "brown", "fox", "because", "it",
//	"owed", "the", "quicker", "fox", "$3.14!"}
func (t *WhitespaceTokenizer) Tokenize(chunk string) (tokens []string, alwaysNil error) {
	return strings.Fields(chunk), nil
}
