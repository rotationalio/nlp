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

func TestRegexTokenizer(t *testing.T) {
	testcases := []struct {
		Name     string
		Text     string
		Regex    string
		Expected []string
	}{
		{
			Name:  "QuickBrownFox-REGEX_ENGLISH_WORDS",
			Text:  "The quick brown fox jumps over the lazy fox.",
			Regex: tokens.REGEX_ENGLISH_WORDS,
			Expected: []string{
				"The",
				"quick",
				"brown",
				"fox",
				"jumps",
				"over",
				"the",
				"lazy",
				"fox",
			},
		},
		{
			Name:  "QuickBrownFox-REGEX_ENGLISH_ALPHABET_ONLY",
			Text:  "The quick brown fox jumps over the lazy fox.",
			Regex: tokens.REGEX_ENGLISH_ALPHABET_ONLY,
			Expected: []string{
				"The",
				"quick",
				"brown",
				"fox",
				"jumps",
				"over",
				"the",
				"lazy",
				"fox",
			},
		},
		{
			Name:  "QuickBrownFoxWithSymbolsAndNumbers-REGEX_ENGLISH_WORDS",
			Text:  "\tThe **&^$% quick% &brown$ ^fox@ %jumps!\n(over) [the] {lazy} 'fox'. 100% 99.9 F.B.I.\r\n _snake_case_",
			Regex: tokens.REGEX_ENGLISH_WORDS,
			Expected: []string{
				"The",
				"quick",
				"brown",
				"fox",
				"jumps",
				"over",
				"the",
				"lazy",
				"fox",
				"100",
				// from "99.9"
				"99",
				"9",
				// from "F.B.I."
				"F",
				"B",
				"I",
				// from "_snake_case_"
				"_snake_case_",
			},
		},
		{
			Name:  "QuickBrownFoxWithSymbolsAndNumbers-REGEX_ENGLISH_ALPHABET_ONLY",
			Text:  "\tThe **&^$% quick% &brown$ ^fox@ %jumps!\n(over) [the] {lazy} 'fox'. 100% 99.9 F.B.I.\r\n _snake_case_",
			Regex: tokens.REGEX_ENGLISH_ALPHABET_ONLY,
			Expected: []string{
				"The",
				"quick",
				"brown",
				"fox",
				"jumps",
				"over",
				"the",
				"lazy",
				"fox",
				// does not capture the numbers
				// from "F.B.I."
				"F",
				"B",
				"I",
				// does not capture the _snake_case_
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			//setup
			tokenizer := tokens.NewRegexTokenizer(tokens.RegexTokenizerWithRegex(tc.Regex))
			require.NotNil(t, tokenizer)

			//test
			tokens, err := tokenizer.Tokenize(tc.Text)
			require.NoError(t, err)
			require.NotNil(t, tokens)
			require.Equal(t, tc.Expected, tokens)
		})
	}
}
