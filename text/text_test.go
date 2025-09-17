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
		require.NotNil(t, myText.WhitespaceTokenizer())
		require.NotNil(t, myText.SentenceSegmenter())
		require.NotNil(t, myText.SSPSyllableTokenizer())
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

func TestWordsAndCount(t *testing.T) {
	myText, err := text.New("apple bananna aardvark aardvarks zebra")
	require.NoError(t, err)
	require.NotNil(t, myText)

	expected := tokenlist.New([]string{"apple", "bananna", "aardvark", "aardvarks", "zebra"})
	require.Nil(t, myText.WordsCache())
	words := myText.Words()
	require.Equal(t, expected, words)
	require.Equal(t, expected, myText.WordsCache())
	require.Equal(t, len(expected), myText.WordCount())
}

func TestSentencesAndCount(t *testing.T) {
	myText, err := text.New("The quick brown fox, Mr. Fox, jumped over the quicker-- 105.4% quicker, in fact-- \n brown fox, the Hon. Judge Fox, because it owed the quicker fox $3.14! Isn't that amazing!?\n I think so!\t Crazy times, indeed. Ellipses... Interrobang!? This last sentence has no punctuation at the end")
	require.NoError(t, err)
	require.NotNil(t, myText)

	expected := tokenlist.New([]string{
		"The quick brown fox, Mr. Fox, jumped over the quicker-- 105.4% quicker, in fact-- brown fox, the Hon. Judge Fox, because it owed the quicker fox $3.14!",
		"Isn't that amazing!?",
		"I think so!",
		"Crazy times, indeed.",
		"Ellipses...",
		"Interrobang!?",
		"This last sentence has no punctuation at the end",
	})
	require.Nil(t, myText.SentencesCache())
	actual := myText.Sentences()
	require.Equal(t, expected, actual)
	require.Equal(t, expected, myText.SentencesCache())
	require.Equal(t, len(expected), myText.SentenceCount())
}

func TestSyllablesAndCount(t *testing.T) {
	myText, err := text.New("justification ice-nine ice9 ice 9 two words")
	require.NoError(t, err)
	require.NotNil(t, myText)

	expected := [][]string{
		{"jus", "ti", "fi", "ca", "tion"},
		{"i", "ce", "-", "ni", "ne"},
		{"i", "ce9"},
		{"i", "ce"},
		{"9"},
		{"two"},
		{"words"},
	}
	require.Nil(t, myText.SyllablesCache())
	actual := myText.Syllables()
	require.Equal(t, expected, actual)
	require.Equal(t, expected, myText.SyllablesCache())
	require.Equal(t, 17, myText.SyllableCount())
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

func TestFleschKincaidErrorsOnly(t *testing.T) {
	myText, err := text.New("The cat sat on the mat.")
	require.NoError(t, err)
	require.NotNil(t, myText)

	// We only need to make sure that no errors or panics happen; the
	// correctness tests are in the [readability_test] package.
	_, err = myText.FleschKincaidReadingEase()
	require.NoError(t, err)
	_, err = myText.FleschKincaidGradeLevel()
	require.NoError(t, err)
}

// Tests that the docstring for [text.Text] work properly; if this ever fails
// please fix it and then copy the lines that do not have the 'require' checks
// into that functions docstring.
func TestTextDocs(t *testing.T) {
	// Create a [text.Text] with the default settings
	myText, err := text.New("apple aardvarks zebra bananna aardvark")
	require.NoError(t, err)
	require.NotNil(t, myText)

	// Get all of the word tokens
	myTokens, err := myText.Tokens() // TokenList
	require.NoError(t, err)
	require.NotNil(t, myTokens)
	require.Len(t, myTokens, 5)

	// Get all word stem tokens which use the same underlying types as the full
	// word tokens above (ignoring errors in this example)
	myStems, err := myText.Stems() // TokenList
	require.NoError(t, err)
	require.NotNil(t, myStems)
	require.Len(t, myStems, 5)

	// The stems are 1:1 count with the tokens
	if len(myTokens) != len(myStems) { // 5 == 5
		panic("this should never occur")
	}

	// You can also get a type count, which returns the count of each unique
	// word stem (ignoring errors) ("aardvark" has a 2 count for this example)
	myCount, err := myText.TypeCount() // map[string]int
	require.NoError(t, err)
	require.NotNil(t, myCount)
	require.Equal(t, 2, myCount["aardvark"])

	// These are a [tokenlist.TokenList], but if you need a slice of strings...
	stringTokens := myTokens.Strings() // []string
	require.Equal(t, []string{"apple", "aardvarks", "zebra", "bananna", "aardvark"}, stringTokens)

	// You can also use regular slice functions and operations on a [tokenlist.TokenList]
	length := len(myTokens) // 5
	require.Equal(t, 5, length)
	myTokens = append(myTokens, myTokens[0]) // "apple", "aardvarks", "zebra", "bananna", "aardvark", "apple"
	require.Equal(t, []string{"apple", "aardvarks", "zebra", "bananna", "aardvark", "apple"}, myTokens.Strings())
	myTokens[0] = myTokens[1] // "aardvarks", "aardvarks", "zebra", "bananna", "aardvark", "apple"
	require.Equal(t, []string{"aardvarks", "aardvarks", "zebra", "bananna", "aardvark", "apple"}, myTokens.Strings())

	// Get an individual token
	firstToken := myTokens[0] // Token
	require.Equal(t, "aardvarks", firstToken.String())

	// You can also get a token as another type
	stringToken := firstToken.String() // string
	require.Equal(t, stringToken, "aardvarks")
	runeToken := firstToken.Runes() // []rune
	require.Equal(t, runeToken, []rune("aardvarks"))
	byteToken := firstToken.Bytes() // []byte
	require.Equal(t, byteToken, []byte("aardvarks"))

	// For these examples, we need to re-create the [text.Text] with a vocabulary,
	// so the [vectorize.CountVectorizer] will work without an error to get
	// cosine similarity.
	myText, err = text.New(
		"cars have engines like motorcycles have engines",
		text.WithVocabulary([]string{"car", "engine", "brakes", "transmission"}),
	)
	require.NoError(t, err)
	require.NotNil(t, myText)
	otherText, err := text.New(
		"engines are attached to transmissions",
		text.WithVocabulary([]string{"car", "engine", "brakes", "transmission"}),
	)
	require.NoError(t, err)
	require.NotNil(t, otherText)

	// Cosine similarity with another string
	similarity, err := myText.CosineSimilarity(otherText) // ~0.5
	require.NoError(t, err)
	require.InDelta(t, 0.5, similarity, 1e-12)

	// We can also get a one-hot or frequency vectorization of our text
	myOneHotVector, err := myText.VectorizeOneHot() // vector.Vector{1, 1, 0, 0}
	require.NoError(t, err)
	require.Equal(t, vector.Vector{1, 1, 0, 0}, myOneHotVector)
	myFrequencyVector, err := myText.VectorizeFrequency() // vector.Vector{1, 2, 0, 0}
	require.NoError(t, err)
	require.Equal(t, vector.Vector{1, 2, 0, 0}, myFrequencyVector)
}
