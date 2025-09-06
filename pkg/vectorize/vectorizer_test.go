package vectorize_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/language"
	"go.rtnl.ai/nlp/pkg/stem"
	"go.rtnl.ai/nlp/pkg/tokenize"
	"go.rtnl.ai/nlp/pkg/vector"
	"go.rtnl.ai/nlp/pkg/vectorize"
)

func TestNewCountVectorizer(t *testing.T) {
	t.Run("SuccessDefaults", func(t *testing.T) {
		vectorizer, err := vectorize.NewCountVectorizer()
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
	})

	t.Run("SuccessVocabOption", func(t *testing.T) {
		//setup
		vocab := []string{"one", "two", "three"}
		vocabOpt := vectorize.CountVectorizerWithVocab(vocab)

		//test
		vectorizer, err := vectorize.NewCountVectorizer(vocabOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, vocab, vectorizer.Vocab())
	})

	t.Run("SuccessLanguageOption_LanguageEnglish", func(t *testing.T) {
		//setup
		lang := language.English
		langOpt := vectorize.CountVectorizerWithLang(lang)

		//test
		vectorizer, err := vectorize.NewCountVectorizer(langOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, lang, vectorizer.Language())
	})

	t.Run("SuccessTokenizerOption_RegexTokenizer", func(t *testing.T) {
		//setup
		tokenizer := tokenize.NewRegexTokenizer()
		tokOpt := vectorize.CountVectorizerWithTokenizer(tokenizer)

		//test
		vectorizer, err := vectorize.NewCountVectorizer(tokOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, tokenizer, vectorizer.Tokenizer())
	})

	t.Run("SuccessStemmerOption_Porter2Stemmer", func(t *testing.T) {
		//setup
		stemmer, err := stem.NewPorter2Stemmer(language.English)
		require.NoError(t, err)
		stemOpt := vectorize.CountVectorizerWithStemmer(stemmer)

		//test
		vectorizer, err := vectorize.NewCountVectorizer(stemOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, stemmer, vectorizer.Stemmer())
	})

	t.Run("SuccessTypeCounterOption_TypeCounter", func(t *testing.T) {
		//setup
		typecounter, err := tokenize.NewTypeCounter()
		require.NoError(t, err)
		tcOpt := vectorize.CountVectorizerWithTypeCounter(typecounter)

		//test
		vectorizer, err := vectorize.NewCountVectorizer(tcOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, typecounter, vectorizer.TypeCounter())
	})

	t.Run("SuccessMethodOption_Frequency", func(t *testing.T) {
		//setup
		method := vectorize.VectorizeFrequency
		methodOpt := vectorize.CountVectorizerWithMethod(method)

		//test
		vectorizer, err := vectorize.NewCountVectorizer(methodOpt)
		require.NoError(t, err)
		require.NotNil(t, vectorizer)
		require.Equal(t, method, vectorizer.Method())
	})
}

func TestCountVectorizerVectorize(t *testing.T) {
	defaultVocab := []string{"apple", "bananna", "cat", "xylophone", "youngster", "zebra"}
	testcases := []struct {
		Name     string
		Method   vectorize.VectorizationMethod
		Vocab    []string
		Text     string
		Expected vector.Vector
		Error    error
	}{
		{
			Name:     "AppleBanannaCat_OneHot",
			Method:   vectorize.VectorizeOneHot,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat",
			Expected: vector.Vector{1, 1, 1, 0, 0, 0},
			Error:    nil,
		},
		{
			Name:     "AppleBanannaCat_Frequency",
			Method:   vectorize.VectorizeFrequency,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat",
			Expected: vector.Vector{1, 1, 1, 0, 0, 0},
			Error:    nil,
		},
		{
			Name:     "DoubledAppleBanannaCat_OneHot",
			Method:   vectorize.VectorizeOneHot,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat apple bananna cat",
			Expected: vector.Vector{1, 1, 1, 0, 0, 0},
			Error:    nil,
		},
		{
			Name:     "DoubledAppleBanannaCat_Frequency",
			Method:   vectorize.VectorizeFrequency,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat apple bananna cat",
			Expected: vector.Vector{2, 2, 2, 0, 0, 0},
			Error:    nil,
		},
		{
			Name:     "WholeVocab_OneHot",
			Method:   vectorize.VectorizeOneHot,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat xylophone youngster zebra",
			Expected: vector.Vector{1, 1, 1, 1, 1, 1},
			Error:    nil,
		},
		{
			Name:     "WholeVocab_Frequency",
			Method:   vectorize.VectorizeFrequency,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat xylophone youngster zebra",
			Expected: vector.Vector{1, 1, 1, 1, 1, 1},
			Error:    nil,
		},
		{
			Name:     "DoubledWholeVocab_OneHot",
			Method:   vectorize.VectorizeOneHot,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat xylophone youngster zebra apple bananna cat xylophone youngster zebra",
			Expected: vector.Vector{1, 1, 1, 1, 1, 1},
			Error:    nil,
		},
		{
			Name:     "DoubledWholeVocab_Frequency",
			Method:   vectorize.VectorizeFrequency,
			Vocab:    defaultVocab,
			Text:     "apple bananna cat xylophone youngster zebra apple bananna cat xylophone youngster zebra",
			Expected: vector.Vector{2, 2, 2, 2, 2, 2},
			Error:    nil,
		},
	}

	for _, tc := range testcases {

		//setup
		vectorizer, err := vectorize.NewCountVectorizer(
			vectorize.CountVectorizerWithMethod(tc.Method),
			vectorize.CountVectorizerWithVocab(tc.Vocab),
		)
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
