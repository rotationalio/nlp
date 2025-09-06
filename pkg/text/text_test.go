package text_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/language"
	"go.rtnl.ai/nlp/pkg/stem"
	"go.rtnl.ai/nlp/pkg/text"
	"go.rtnl.ai/nlp/pkg/tokenize"
)

func TestNew(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
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
