package similarity_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/errors"
	"go.rtnl.ai/nlp/pkg/similarity"
	"go.rtnl.ai/nlp/pkg/tokens"
	"go.rtnl.ai/nlp/pkg/vector"
)

func TestNewCosineSimilarizer(t *testing.T) {
	t.Run("SuccessDefaults", func(t *testing.T) {
		sim, err := similarity.NewCosineSimilarizer([]string{"this", "is", "a", "test"})
		require.NoError(t, err)
		require.NotNil(t, sim)
	})

	t.Run("SuccessLanguageOption_LanguageEnglish", func(t *testing.T) {
		// setup
		lang := enum.LanguageEnglish
		optLang := similarity.CosineSimilarizerWithLanguage(lang)
		vocab := []string{"this", "is", "a", "test"}

		// test
		sim, err := similarity.NewCosineSimilarizer(vocab, optLang)
		require.NoError(t, err)
		require.NotNil(t, sim)
		require.Equal(t, lang, sim.Language())
	})

	t.Run("SuccessTokenizerOption_RegexTokenizer", func(t *testing.T) {
		// setup
		tok := tokens.NewRegexTokenizer()
		optTok := similarity.CosineSimilarizerWithTokenizer(tok)
		vocab := []string{"this", "is", "a", "test"}

		// test
		sim, err := similarity.NewCosineSimilarizer(vocab, optTok)
		require.NoError(t, err)
		require.NotNil(t, sim)
		require.Equal(t, tok, sim.Tokenizer())
	})

	t.Run("SuccessVectorizerOption_CountVectorizer", func(t *testing.T) {
		// setup
		vocab := []string{"this", "is", "a", "test"}
		vec, err := vector.NewCountVectorizer(vocab)
		require.NoError(t, err)
		optVec := similarity.CosineSimilarizerWithVectorizer(vec)

		// test
		sim, err := similarity.NewCosineSimilarizer(vocab, optVec)
		require.NoError(t, err)
		require.NotNil(t, sim)
		require.Equal(t, vec, sim.Vectorizer())
	})
}

func TestCosineSimilarity(t *testing.T) {
	testcases := []struct {
		Name     string
		First    string
		Second   string
		Expected float64
		Error    error
	}{
		{
			Name:     "SuccessExactMatch",
			First:    "apple bananna cat",
			Second:   "apple bananna cat",
			Expected: 1.0,
			Error:    nil,
		},
		{
			Name:     "SuccessZeroMatch",
			First:    "xylophone youngster zebra",
			Second:   "apple bananna cat",
			Expected: 0.0,
			Error:    nil,
		},
		{
			Name:     "SuccessTwoThirdsMatch",
			First:    "apple youngster zebra",
			Second:   "apple bananna zebra",
			Expected: (2.0 / 3.0),
			Error:    nil,
		},
		{
			Name:     "ErrorUndefinedValue",
			First:    "returns error undefined value",
			Second:   "returns error undefined value",
			Expected: 0.0,
			Error:    errors.ErrUndefinedValue,
		},
	}

	// setup
	vocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
	similarizer, err := similarity.NewCosineSimilarizer(vocab)
	require.NoError(t, err)
	require.NotNil(t, similarizer)

	// tests
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			sim, err := similarizer.Similarity(tc.First, tc.Second)
			if tc.Error != nil {
				require.Error(t, err, tc.Error)
			} else {
				require.NoError(t, err)
			}
			require.InDeltaf(t, tc.Expected, sim, 1e-12, "expected %f got %f a difference of %e", tc.Expected, sim, math.Abs(tc.Expected-sim))
		})
	}
}
