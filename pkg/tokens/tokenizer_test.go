package tokens_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/tokens"
)

func TestNewRegexTokenizer(t *testing.T) {
	t.Run("SuccessDefaults", func(t *testing.T) {
		tok := tokens.NewRegexTokenizer()
		require.NotNil(t, tok)
	})

	t.Run("SuccessLanguageOption_LanguageEnglish", func(t *testing.T) {
		//setup
		lang := enum.LanguageEnglish
		langOpt := tokens.RegexTokenizerWithLanguage(lang)

		//test
		tok := tokens.NewRegexTokenizer(langOpt)
		require.NotNil(t, tok)
		require.Equal(t, lang, tok.Language())
	})

	t.Run("SuccessRegexOption_REGEX_ENGLISH_ALPHABET_ONLY", func(t *testing.T) {
		//setup
		regex := tokens.REGEX_ENGLISH_ALPHABET_ONLY
		regexOpt := tokens.RegexTokenizerWithRegex(regex)

		//test
		tok := tokens.NewRegexTokenizer(regexOpt)
		require.NotNil(t, tok)
		require.Equal(t, regex, tok.Regex())
	})
}
