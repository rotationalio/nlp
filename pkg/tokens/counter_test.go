package tokens_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/stemming"
	"go.rtnl.ai/nlp/pkg/tokens"
)

func TestNewTypeCounter(t *testing.T) {
	t.Run("SuccessDefaults", func(t *testing.T) {
		tc, err := tokens.NewTypeCounter()
		require.NoError(t, err)
		require.NotNil(t, tc)
	})

	t.Run("SuccessLanguageOption_LanguageEnglish", func(t *testing.T) {
		//setup
		lang := enum.LanguageEnglish
		langOpt := tokens.TypeCounterWithLanguage(lang)

		//test
		tc, err := tokens.NewTypeCounter(langOpt)
		require.NoError(t, err)
		require.NotNil(t, tc)
		require.Equal(t, enum.LanguageEnglish, tc.Languge())
	})

	t.Run("SuccessTokenizerOption_RegexTokenizer", func(t *testing.T) {
		//setup
		tok := tokens.NewRegexTokenizer()
		tokOpt := tokens.TypeCounterWithTokenizer(tok)

		//test
		tc, err := tokens.NewTypeCounter(tokOpt)
		require.NoError(t, err)
		require.NotNil(t, tc)
		require.Equal(t, tok, tc.Tokenizer())
	})

	t.Run("SuccessStemmerOption_Porter2Stemmer", func(t *testing.T) {
		//setup
		lang := enum.LanguageEnglish
		stemmer, err := stemming.NewPorter2Stemmer(lang)
		require.NoError(t, err)
		stemOpt := tokens.TypeCounterWithStemmer(stemmer)

		//test
		tc, err := tokens.NewTypeCounter(stemOpt)
		require.NoError(t, err)
		require.NotNil(t, tc)
		require.Equal(t, stemmer, tc.Stemmer())
	})
}

// NOTE: this test relies on specific default settings in the
// [tokens.TypeCounter] implementation, such as the default stemmer and default
// tokenizer, so if the defaults for any of the chain of tools used by the
// [tokens.TypeCounter] changes then this test will need to be repaired.
func TestTypeCounterTypeCount(t *testing.T) {
	t.Run("SuccessQuickBrownFox", func(t *testing.T) {
		//setup
		typecounter, err := tokens.NewTypeCounter()
		require.NoError(t, err)

		text := "The quick brown fox jumps over the lazy fox."
		expected := map[string]int{
			"the":   2,
			"quick": 1,
			"brown": 1,
			"fox":   2,
			"jump":  1,
			"over":  1,
			"lazi":  1,
		}

		//test
		count, err := typecounter.TypeCount(text)
		require.NoError(t, err)
		require.NotNil(t, count)
		require.InDeltaMapValues(t, expected, count, 0.0)
	})

	t.Run("SuccessQuickBrownFoxWithSymbolsAndNumbers", func(t *testing.T) {
		//setup
		typecounter, err := tokens.NewTypeCounter()
		require.NoError(t, err)

		text := "\tThe **&^$% quick% &brown$ ^fox@ %jumps!\n(over) [the] {lazy} 'fox'. 100% 99.9 F.B.I.\r\n _snake_case_"
		expected := map[string]int{
			"the":   2,
			"quick": 1,
			"brown": 1,
			"fox":   2,
			"jump":  1,
			"over":  1,
			"lazi":  1,
			"100":   1,
			// from "99.9"
			"99": 1,
			"9":  1,
			// from "F.B.I."
			"f": 1,
			"b": 1,
			"i": 1,
			// from "_snake_case_"
			"_snake_case_": 1,
		}

		//test
		count, err := typecounter.TypeCount(text)
		require.NoError(t, err)
		require.NotNil(t, count)
		require.InDeltaMapValues(t, expected, count, 0.0)
	})
}
