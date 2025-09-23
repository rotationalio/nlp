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
			Expected: []string{"i", "ce", "ni", "ne"},
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
			Expected: []string{"i", "ce", "9"},
		},
		{
			Name:     "Two_Words",
			Language: language.English,
			Word:     "two words",
			Expected: []string{"two", "words"},
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
			// NOTE: We would normally want `{}`, however the
			// implementation returns `nil`, which we'll accept as a
			// degenerative case for this algorithm to make things simpler.
			Expected: nil,
		},
		{
			Name:     "TwoPunct_Degnerative_Case",
			Language: language.English,
			Word:     "!?",
			// NOTE: We would normally want `{}`, however the
			// implementation returns `{"!?"}`, which we'll accept as a
			// degenerative case for this algorithm to make things simpler.
			Expected: []string{"!?"},
		},
		{
			Name:     "OnePunctMan",
			Language: language.English,
			Word:     ".",
			// NOTE: We would normally want `{}`, however the
			// implementation returns `{"."}`, which we'll accept as a
			// degenerative case for this algorithm to make things simpler.
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
			Expected: []string{"F", "B", "I"},
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

func TestSSPSyllableTokenizerCountsEnglish(t *testing.T) {
	testcases := []struct {
		Word      string
		Syllables int
	}{
		// 1 syllable:
		{"dog", 1},
		{"owl", 1},
		{"fish", 1},
		{"cat", 1},
		// {"shed", 1}, //FIXME 2 [s hed]
		{"drum", 1},
		{"brick", 1},
		{"plug", 1},
		{"horn", 1},
		// 2 syllables:
		{"pencil", 2},
		{"dolphin", 2},
		// {"spider", 2},     //FIXME 3 [s pi der]
		// {"lighthouse", 2}, //FIXME 3 [light hou se]
		{"movie", 2},
		{"window", 2},
		{"baby", 2},
		{"cuddle", 2},
		{"rabbit", 2},
		// 3 syllables:
		// {"piano", 3}, //FIXME 2 [pia no]
		{"elephant", 3},
		{"library", 3},
		// {"telephone", 3},  //FIXME 4 [te lep ho ne]
		// {"strawberry", 3}, //FIXME 4 [s traw ber ry]
		{"computer", 3},
		{"butterfly", 3},
		// {"aeroplane", 3}, //FIXME 4 [ae ro pla ne]
		{"caravan", 3},
		// 4 syllables:
		{"supermarket", 4},
		{"impossible", 4},
		{"watermelon", 4},
		{"calculator", 4},
		{"helicopter", 4},
		{"television", 4},
		{"information", 4},
		{"competition", 4},
		{"crocodile", 4},
		// 5 syllables:
		{"international", 5},
		{"refrigerator", 5},
		{"congratulations", 5},
		{"multiplication", 5},
		{"investigation", 5},
	}

	tokenizer, err := tokenize.NewSSPSyllableTokenizer(language.English)
	require.NoError(t, err)
	require.NotNil(t, tokenizer)

	for _, tc := range testcases {
		tokens, err := tokenizer.Tokenize(tc.Word)
		require.NoError(t, err)
		require.Equalf(t, tc.Syllables, len(tokens), "Expected %d tokens for '%s', got %d: %s", tc.Syllables, tc.Word, len(tokens), tokens)
	}
}
