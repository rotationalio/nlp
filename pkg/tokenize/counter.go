package tokenize

import (
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/stem"
)

// ############################################################################
// TypeCounter
// ############################################################################

// TypeCounter can be used to perform type counting on text; create with [NewTypeCounter].
type TypeCounter struct {
	lang      enum.Language
	tokenizer Tokenizer
	stemmer   stem.Stemmer
}

// Returns a new [TypeCounter] instance. Defaults to the default [RegexTokenizer] and
// [stem.Stemmer] options. Modified by passing [TypeCounterOption] functions into
// relevant function calls.
//
// Defaults:
//   - Language: [LanguageEnglish]
//   - Tokenizer: [RegexTokenizer]
//   - Stemmer: [Porter2Stemmer]
func NewTypeCounter(opts ...TypeCounterOption) (tc *TypeCounter, err error) {
	// Set options
	tc = &TypeCounter{}
	for _, fn := range opts {
		fn(tc)
	}

	// Set defaults
	if tc.lang == enum.LanguageUnknown {
		tc.lang = enum.LanguageEnglish
	}

	if tc.tokenizer == nil {
		tc.tokenizer = NewRegexTokenizer()
	}

	if tc.stemmer == nil {
		if tc.stemmer, err = stem.NewPorter2Stemmer(tc.lang); err != nil {
			return nil, err
		}
	}

	return tc, nil
}

// Returns the [TypeCounter]s configured [enum.Language].
func (c *TypeCounter) Languge() enum.Language {
	return c.lang
}

// Returns the [TypeCounter]s configured [Tokenizer].
func (c *TypeCounter) Tokenizer() Tokenizer {
	return c.tokenizer
}

// Returns the [TypeCounter]s configured [stem.Stemmer].
func (c *TypeCounter) Stemmer() stem.Stemmer {
	return c.stemmer
}

// Returns a map of the types (unique words) and their counts for the string.
func (c *TypeCounter) TypeCount(chunk string) (types map[string]int, err error) {
	// Tokenize
	var tokens []string
	if tokens, err = c.tokenizer.Tokenize(chunk); err != nil {
		return nil, err
	}

	// Stem
	for i, tok := range tokens {
		tokens[i] = c.stemmer.Stem(tok)
	}

	// Count
	return c.CountTypes(tokens), nil
}

// CountTypes returns a the count of each type (unique word) in the given token
// list.
func (c *TypeCounter) CountTypes(tokens []string) (types map[string]int) {
	sz := len(tokens) / 5 // map size selected arbitrarily
	types = make(map[string]int, sz)
	for _, tok := range tokens {
		types[tok] += 1
	}
	return types
}

// ############################################################################
// TypeCounterOptions
// ############################################################################

// TypeCounterOption functions modify a [TypeCounter].
type TypeCounterOption func(t *TypeCounter)

// TypeCounterWithLanguage sets the [enum.Language] to be used for a [TypeCounter].
func TypeCounterWithLanguage(lang enum.Language) TypeCounterOption {
	return func(t *TypeCounter) {
		t.lang = lang
	}
}

// TypeCounterWithTokenizer sets the [Tokenizer] to be used for a [TypeCounter].
func TypeCounterWithTokenizer(tokenizer Tokenizer) TypeCounterOption {
	return func(t *TypeCounter) {
		t.tokenizer = tokenizer
	}
}

// TypeCounterWithStemmer sets the [stem.Stemmer] to be used for a [TypeCounter].
func TypeCounterWithStemmer(stemmer stem.Stemmer) TypeCounterOption {
	return func(t *TypeCounter) {
		t.stemmer = stemmer
	}
}
