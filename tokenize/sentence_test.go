package tokenize_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/language"
	"go.rtnl.ai/nlp/tokenize"
)

func TestNewSentenceSegmenter(t *testing.T) {
	t.Run("SuccessDefaults", func(t *testing.T) {
		segmenter := tokenize.NewSentenceSegmenter()
		require.NotNil(t, segmenter)
		require.Equal(t, language.English, segmenter.Language())
	})

	t.Run("SuccessLanguageOption_LANGUAGE_ENGLISH", func(t *testing.T) {
		lang := language.English
		segmenter := tokenize.NewSentenceSegmenter(tokenize.SentenceSegmenterWithLanguage(lang))
		require.NotNil(t, segmenter)
		require.Equal(t, lang, segmenter.Language())
	})
}

func TestSentenceSegmenter(t *testing.T) {
	segmenter := tokenize.NewSentenceSegmenter(tokenize.SentenceSegmenterWithLanguage(language.English))
	require.NotNil(t, segmenter)

	expected := []string{
		"The quick brown fox, Mr. Fox, jumped over the quicker-- 105.4% quicker, in fact-- brown fox, the Hon. Judge Fox, because it owed the quicker fox $3.14!",
		"Isn't that amazing!?",
		"I think so!",
		"Crazy times, indeed.",
		"Ellipses...",
		"Interrobang!?",
		"This last sentence has no punctuation at the end",
	}

	sentences, err := segmenter.Tokenize("The quick brown fox, Mr. Fox, jumped over the quicker-- 105.4% quicker, in fact-- \n brown fox, the Hon. Judge Fox, because it owed the quicker fox $3.14! Isn't that amazing!?\n I think so!\t Crazy times, indeed. Ellipses... Interrobang!? This last sentence has no punctuation at the end")
	require.Nil(t, err, "error should ALWAYS be nil for the SentenceSegmenter.Tokenize function")
	require.Equal(t, expected, sentences)

}
