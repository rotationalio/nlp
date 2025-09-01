package similarity_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/enum"
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

	t.Run("SuccessLanguageOption", func(t *testing.T) {
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

	t.Run("SuccessTokenizerOption", func(t *testing.T) {
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

	t.Run("SuccessVectorizerOption", func(t *testing.T) {
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
	// TODO test CosineSimilarizer.Similarity() with several different text chunks
}
