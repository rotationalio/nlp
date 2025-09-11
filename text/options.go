package text

import (
	"go.rtnl.ai/nlp/language"
	"go.rtnl.ai/nlp/stem"
	"go.rtnl.ai/nlp/tokenize"
)

// An [Option] configures a [Text] in [New].
type Option func(text *Text)

// Returns a function that sets the vocabulary on a [Text].
func WithVocabulary(vocab []string) Option {
	return func(text *Text) {
		text.vocab = vocab
	}
}

// Returns a function that sets the [language.Language] on a [Text].
func WithLanguage(lang language.Language) Option {
	return func(text *Text) {
		text.lang = lang
	}
}

// Returns a function that sets the [stem.Stemmer] on a [Text].
func WithStemmer(stemmer stem.Stemmer) Option {
	return func(text *Text) {
		text.stemmer = stemmer
	}
}

// Returns a function that sets the [tokenize.Tokenizer] on a [Text].
func WithTokenizer(tokenizer tokenize.Tokenizer) Option {
	return func(text *Text) {
		text.tokenizer = tokenizer
	}
}
