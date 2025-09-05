package tokenize

import (
	"regexp"

	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/tokenlist"
)

// ############################################################################
// Tokenizer interface
// ############################################################################

type Tokenizer interface {
	Tokenize(chunk string) (tokens *tokenlist.TokenList, err error)
}

// ############################################################################
// Regex Expressions for Tokenizing
// ############################################################################

const (
	// English word tokenization for words using the "word character" ('\w')
	// which includes letters, numbers, and underscores
	REGEX_ENGLISH_WORDS = `\b\w+\b`
	// English word tokenization for words using lowercase and uppercase letters
	REGEX_ENGLISH_ALPHABET_ONLY = `\b[A-Za-z]+\b`
)

// ############################################################################
// RegexTokenizer
// ############################################################################

// Ensure [RegexTokenizer] meets the [Tokenizer] interface requirements.
var _ Tokenizer = &RegexTokenizer{}

// RegexTokenizer can be used to tokenize text; create with [NewRegexTokenizer].
type RegexTokenizer struct {
	lang  enum.Language
	regex string
}

// Returns a new [RegexTokenizer] instance.
//
// Defaults:
//   - Language: [LanguageEnglish]
//   - Regex: [REGEX_ENGLISH_WORDS]
func NewRegexTokenizer(opts ...RegexTokenizerOption) *RegexTokenizer {
	// Set options
	tokenizer := &RegexTokenizer{}
	for _, fn := range opts {
		fn(tokenizer)
	}

	// Set defaults
	if tokenizer.lang == enum.LanguageUnknown {
		tokenizer.lang = enum.LanguageEnglish
	}
	if tokenizer.regex == "" {
		tokenizer.regex = REGEX_ENGLISH_WORDS
	}

	return tokenizer
}

// Returns the [RegexTokenizer]s configured [enum.Language]
func (t *RegexTokenizer) Language() enum.Language {
	return t.lang
}

// Returns the [RegexTokenizer]s configured regular expression.
func (t *RegexTokenizer) Regex() string {
	return t.regex
}

// Tokenizes a text string using [regexp.Regexp.FindAllString].
func (t *RegexTokenizer) Tokenize(chunk string) (tokens *tokenlist.TokenList, err error) {
	// Compile regexp
	var r *regexp.Regexp
	if r, err = regexp.Compile(t.regex); err != nil {
		return nil, err
	}

	// Tokenize with regex
	tokens = tokenlist.New(r.FindAllString(chunk, -1))

	return tokens, nil
}

// ############################################################################
// RegexTokenizerOption
// ############################################################################

// RegexTokenizerOption functions modify a [RegexTokenizer].
type RegexTokenizerOption func(t *RegexTokenizer)

// Returns a function which sets the [enum.Language] to use with the [RegexTokenizer].
func RegexTokenizerWithLanguage(lang enum.Language) RegexTokenizerOption {
	return func(t *RegexTokenizer) {
		t.lang = lang
	}
}

// Returns a function which sets the regular expression to use with the
// [RegexTokenizer].
func RegexTokenizerWithRegex(regex string) RegexTokenizerOption {
	return func(t *RegexTokenizer) {
		t.regex = regex
	}
}
