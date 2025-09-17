package tokenize_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/language"
	"go.rtnl.ai/nlp/tokenize"
)

func TestSSPSyllableTokenizer(t *testing.T) {
	testcases := []struct {
		Name     string
		Language language.Language
		Word     string
		Expected []string
	}{
		{
			Name:     "NLTK_justification", // a test in NLTK
			Language: language.English,
			Word:     "justification",
			Expected: []string{"jus", "ti", "fi", "ca", "tion"},
		},
		{
			Name:     "NLTK_10k_Nines", // a test in NLTK
			Language: language.English,
			Word:     strings.Repeat("9", 10_000),
			Expected: []string{strings.Repeat("9", 10_000)},
		},
		{
			Name:     "Punctuation",
			Language: language.English,
			Word:     "ice-nine",
			Expected: []string{"i", "ce", "-", "ni", "ne"},
		},
		{
			Name:     "Alphanumeric",
			Language: language.English,
			Word:     "ice9",
			Expected: []string{"i", "ce9"},
		},
		{
			Name:     "Alpha_Space_Numeric",
			Language: language.English,
			Word:     "ice 9",
			Expected: []string{"i", "ce", " ", "9"},
		},
		{
			Name:     "Two_Words",
			Language: language.English,
			Word:     "two words",
			Expected: []string{"two", " ", "words"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			tokenizer, err := tokenize.NewSSPSyllableTokenizer(tc.Language)
			require.NoError(t, err)
			require.NotNil(t, tokenizer)

			tokens, err := tokenizer.Tokenize(tc.Word)
			require.Nil(t, err, "the error should ALWAYS be nil from this function")
			require.Equal(t, tc.Expected, tokens)
		})
	}
}
