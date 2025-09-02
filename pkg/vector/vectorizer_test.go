package vector_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/stemming"
	"go.rtnl.ai/nlp/pkg/tokens"
	"go.rtnl.ai/nlp/pkg/vector"
)

func TestNewCountVectorizer(t *testing.T) {
	t.Run("SuccessDefaults", func(t *testing.T) {
		vocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
		vectorizer, err := vector.NewCountVectorizer(vocab)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
	})

	t.Run("SuccessLanguageOption_LanguageEnglish", func(t *testing.T) {
		//setup
		vocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
		lang := enum.LanguageEnglish
		langOpt := vector.CountVectorizerWithLang(lang)

		//test
		vectorizer, err := vector.NewCountVectorizer(vocab, langOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, lang, vectorizer.Language())
	})

	t.Run("SuccessTokenizerOption_RegexTokenizer", func(t *testing.T) {
		//setup
		vocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
		tokenizer := tokens.NewRegexTokenizer()
		tokOpt := vector.CountVectorizerWithTokenizer(tokenizer)

		//test
		vectorizer, err := vector.NewCountVectorizer(vocab, tokOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, tokenizer, vectorizer.Tokenizer())
	})

	t.Run("SuccessStemmerOption_Porter2Stemmer", func(t *testing.T) {
		//setup
		vocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
		stemmer, err := stemming.NewPorter2Stemmer(enum.LanguageEnglish)
		require.NoError(t, err)
		stemOpt := vector.CountVectorizerWithStemmer(stemmer)

		//test
		vectorizer, err := vector.NewCountVectorizer(vocab, stemOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, stemmer, vectorizer.Stemmer())
	})

	t.Run("SuccessTypeCounterOption_TypeCounter", func(t *testing.T) {
		//setup
		vocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
		typecounter, err := tokens.NewTypeCounter()
		require.NoError(t, err)
		tcOpt := vector.CountVectorizerWithTypeCounter(typecounter)

		//test
		vectorizer, err := vector.NewCountVectorizer(vocab, tcOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, typecounter, vectorizer.TypeCounter())
	})

	t.Run("SuccessMethodOption_Frequency", func(t *testing.T) {
		//setup
		vocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
		method := vector.VectorizeFrequency
		methodOpt := vector.CountVectorizerWithMethod(method)

		//test
		vectorizer, err := vector.NewCountVectorizer(vocab, methodOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, method, vectorizer.Method())
	})
}

func TestCountVectorizerVectorize(t *testing.T) {
	defaultVocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
	testcases := []struct {
		Name     string
		Method   vector.VectorizationMethod
		Vocab    []string
		Text     string
		Expected vector.Vector
		Error    error
	}{
		{
			Name:     "AppleBanannaCat_OneHot",
			Method:   vector.VectorizeOneHot,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat",
			Expected: vector.Vector{1, 1, 1, 0, 0, 0},
			Error:    nil,
		},
		{
			Name:     "AppleBanannaCat_Frequency",
			Method:   vector.VectorizeFrequency,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat",
			Expected: vector.Vector{1, 1, 1, 0, 0, 0},
			Error:    nil,
		},
		{
			Name:     "DoubledAppleBanannaCat_OneHot",
			Method:   vector.VectorizeOneHot,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat apple bananna cat",
			Expected: vector.Vector{1, 1, 1, 0, 0, 0},
			Error:    nil,
		},
		{
			Name:     "DoubledAppleBanannaCat_Frequency",
			Method:   vector.VectorizeFrequency,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat apple bananna cat",
			Expected: vector.Vector{2, 2, 2, 0, 0, 0},
			Error:    nil,
		},
		{
			Name:     "WholeVocab_OneHot",
			Method:   vector.VectorizeOneHot,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat xylophone youngster zebra",
			Expected: vector.Vector{1, 1, 1, 1, 1, 1},
			Error:    nil,
		},
		{
			Name:     "WholeVocab_Frequency",
			Method:   vector.VectorizeFrequency,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat xylophone youngster zebra",
			Expected: vector.Vector{1, 1, 1, 1, 1, 1},
			Error:    nil,
		},
		{
			Name:     "DoubledWholeVocab_OneHot",
			Method:   vector.VectorizeOneHot,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat xylophone youngster zebra apple bananna cat xylophone youngster zebra",
			Expected: vector.Vector{1, 1, 1, 1, 1, 1},
			Error:    nil,
		},
		{
			Name:     "DoubledWholeVocab_Frequency",
			Method:   vector.VectorizeFrequency,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat xylophone youngster zebra apple bananna cat xylophone youngster zebra",
			Expected: vector.Vector{2, 2, 2, 2, 2, 2},
			Error:    nil,
		},
	}

	for _, tc := range testcases {
		//setup
		vectorizer, err := vector.NewCountVectorizer(tc.Vocab, vector.CountVectorizerWithMethod(tc.Method))
		require.NoError(t, err)
		require.NotNil(t, vectorizer)

		//test
		actual, err := vectorizer.Vectorize(tc.Text)
		if tc.Error != nil {
			require.Error(t, tc.Error, err)
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, tc.Expected, actual)
	}
}
