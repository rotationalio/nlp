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
		{
			Name:     "TwoChars",
			Language: language.English,
			Word:     "an",
			Expected: []string{"an"},
		},
		{
			Name:     "OneChar",
			Language: language.English,
			Word:     "a",
			Expected: []string{"a"},
		},
		{
			Name:     "ThreePunct",
			Language: language.English,
			Word:     ".!?",
			Expected: []string{".", "!", "?"},
		},
		{
			Name:     "TwoPunct_Degnerative_Case",
			Language: language.English,
			Word:     "!?",
			// NOTE: We would normally want `{"!", "?"}`, however the
			// implementation returns `{"!?"}`, which we'll accept as a
			// degenerative case for this algorithm to make things simpler.
			Expected: []string{"!?"},
		},
		{
			Name:     "OnePunctMan",
			Language: language.English,
			Word:     ".",
			Expected: []string{"."},
		},
		{
			Name:     "EmptyString",
			Language: language.English,
			Word:     "",
			Expected: []string{""},
		},
		{
			Name:     "MixedCasing",
			Language: language.English,
			Word:     "JusTiFiCaTion",
			Expected: []string{"Jus", "Ti", "Fi", "Ca", "Tion"},
		},
		{
			Name:     "UpperCasing",
			Language: language.English,
			Word:     "JUSTIFICATION",
			Expected: []string{"JUS", "TI", "FI", "CA", "TION"},
		},
		{
			Name:     "Initialism",
			Language: language.English,
			Word:     "F.B.I.",
			Expected: []string{"F", ".", "B", ".", "I", "."},
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
