package text_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/language"
	"go.rtnl.ai/nlp/stem"
	"go.rtnl.ai/nlp/text"
	"go.rtnl.ai/nlp/tokenize"
	"go.rtnl.ai/nlp/tokenlist"
	"go.rtnl.ai/nlp/vector"
)

func TestNew(t *testing.T) {
	t.Run("CheckDefaults", func(t *testing.T) {
		myText, err := text.New("testing text.New()")
		require.NoError(t, err)
		require.NotNil(t, myText)

		require.Nil(t, myText.Vocab())
		require.Equal(t, language.English, myText.Language())
		require.IsType(t, &stem.Porter2Stemmer{}, myText.Stemmer())
		require.IsType(t, &tokenize.RegexTokenizer{}, myText.Tokenizer())
		require.NotNil(t, myText.TypeCounter())
		require.NotNil(t, myText.CountVectorizer())
		require.NotNil(t, myText.CosineSimilarizer())
	})

	t.Run("VocabularyOption", func(t *testing.T) {
		vocab := []string{"one", "two"}
		myText, err := text.New("testing text.New()", text.WithVocabulary(vocab))
		require.NoError(t, err)
		require.NotNil(t, myText)
		require.Equal(t, vocab, myText.Vocab())
	})

	t.Run("LanguageOption", func(t *testing.T) {
		lang := language.English
		myText, err := text.New("testing text.New()", text.WithLanguage(lang))
		require.NoError(t, err)
		require.NotNil(t, myText)
		require.Equal(t, lang, myText.Language())
	})

	t.Run("StemmerOption", func(t *testing.T) {
		stemmer, err := stem.NewPorter2Stemmer(language.English)
		require.NoError(t, err)
		require.NotNil(t, stemmer)

		myText, err := text.New("testing text.New()", text.WithStemmer(stemmer))
		require.NoError(t, err)
		require.NotNil(t, myText)
		require.Equal(t, stemmer, myText.Stemmer())
	})

	t.Run("TokenizerOption", func(t *testing.T) {
		tokenizer := tokenize.NewRegexTokenizer()
		require.NotNil(t, tokenizer)

		myText, err := text.New("testing text.New()", text.WithTokenizer(tokenizer))
		require.NoError(t, err)
		require.NotNil(t, myText)
		require.Equal(t, tokenizer, myText.Tokenizer())
	})
}

func TestTextTypeGetters(t *testing.T) {
	s := "testing 12345"
	myText, err := text.New(s)
	require.NoError(t, err)
	require.NotNil(t, myText)

	textText := myText.Text()
	require.Equal(t, s, textText)

	textString := myText.String()
	require.Equal(t, s, textString)

	textRunes := myText.Runes()
	require.Equal(t, []rune(s), textRunes)

	textBytes := myText.Bytes()
	require.Equal(t, []byte(s), textBytes)

}

func TestTokens(t *testing.T) {
	myText, err := text.New("apple bananna aardvark aardvarks zebra")
	require.NoError(t, err)
	require.NotNil(t, myText)

	expected := tokenlist.New([]string{"apple", "bananna", "aardvark", "aardvarks", "zebra"})
	require.Nil(t, myText.TokensCache())
	tokens, err := myText.Tokens()
	require.NoError(t, err)
	require.Equal(t, expected, tokens)
	require.Equal(t, expected, myText.TokensCache())
}

func TestStems(t *testing.T) {
	myText, err := text.New("apple bananna aardvark aardvarks zebra")
	require.NoError(t, err)
	require.NotNil(t, myText)

	expected := tokenlist.New([]string{"appl", "bananna", "aardvark", "aardvark", "zebra"})
	require.Nil(t, myText.StemsCache())
	stems, err := myText.Stems()
	require.NoError(t, err)
	require.Equal(t, expected, stems)
	require.Equal(t, expected, myText.StemsCache())
}

func TestTypeCount(t *testing.T) {
	myText, err := text.New("apple bananna aardvark aardvarks zebra")
	require.NoError(t, err)
	require.NotNil(t, myText)

	expected := map[string]int{"appl": 1, "bananna": 1, "aardvark": 2, "zebra": 1}
	require.Nil(t, myText.TypeCountCache())
	typeCount, err := myText.TypeCount()
	require.NoError(t, err)
	require.Equal(t, expected, typeCount)
	require.Equal(t, expected, myText.TypeCountCache())
}

func TestVectorizeFrequency(t *testing.T) {
	vocab := []string{"one", "two"}
	myText, err := text.New("one one three", text.WithVocabulary(vocab))
	require.NoError(t, err)
	require.NotNil(t, myText)

	expected := vector.Vector{2, 0}
	vector, err := myText.VectorizeFrequency()
	require.NoError(t, err)
	require.NotNil(t, vector)
	require.Equal(t, expected, vector)
}

func TestVectorizeOneHot(t *testing.T) {
	vocab := []string{"one", "two"}
	myText, err := text.New("one one three", text.WithVocabulary(vocab))
	require.NoError(t, err)
	require.NotNil(t, myText)

	expected := vector.Vector{1, 0}
	vector, err := myText.VectorizeOneHot()
	require.NoError(t, err)
	require.NotNil(t, vector)
	require.Equal(t, expected, vector)
}

func TestCosineSimilarity(t *testing.T) {
	vocab := []string{"one", "two", "three"}
	myText, err := text.New("one three", text.WithVocabulary(vocab))
	require.NoError(t, err)
	require.NotNil(t, myText)

	otherText, err := text.New("two three", text.WithVocabulary(vocab))
	require.NoError(t, err)
	require.NotNil(t, otherText)

	expected := 0.5
	similarity, err := myText.CosineSimilarity(otherText)
	require.NoError(t, err)
	require.InDelta(t, expected, similarity, 1e-12)
}
