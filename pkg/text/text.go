package text

import (
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/stem"
	"go.rtnl.ai/nlp/pkg/token"
	"go.rtnl.ai/nlp/pkg/tokenize"
	"go.rtnl.ai/nlp/pkg/tokenlist"
)

// A one-stop-shop for performing NLP operations on a string of text, such as
// stemming, vectorization, or similarity to another string.
type Text struct {
	// The string representation of the text.
	text string

	// The tokens of this text; lazily initialized.
	tokens *tokenlist.TokenList

	// The stem tokens of this text; lazily initialized.
	stems *tokenlist.TokenList

	// The [enum.Language] of this text.
	lang enum.Language

	// The [stemming.Stemmer] to use for stemming this text's tokens.
	stemmer stem.Stemmer

	// The [tokens.Tokenizer] to use for tokenizing this text.
	tokenizer tokenize.Tokenizer
}

// Create a new [Text] from the input string with the specified [Option]s.
//
// Defaults:
// * Language: [enum.LanguageEnglish]
// * Stemmer: [stem.Porter2Stemmer]
// * Tokenizer: [tokenize.RegexTokenizer]
func New(t string, options ...Option) (text *Text, err error) {
	// Initialize text
	text = &Text{
		text:   t,
		tokens: nil, // will tokenize on the first call to Tokens()
		stems:  nil, // will stem the tokens on the first call to Stems()
	}

	// Set user options
	for _, opt := range options {
		opt(text)
	}

	// Default languge
	if text.lang == enum.LanguageUnknown {
		text.lang = enum.LanguageEnglish
	}

	// Default stemmer
	if text.stemmer == nil {
		if text.stemmer, err = stem.NewPorter2Stemmer(text.lang); err != nil {
			return nil, err
		}
	}

	// Default tokenizer
	if text.tokenizer == nil {
		text.tokenizer = tokenize.NewRegexTokenizer(
			tokenize.RegexTokenizerWithLanguage(text.lang),
		)
	}

	return text, nil
}

// Returns a [tokenize.TokenList] for the [Text]s tokens using the configured
// [tokenize.Tokenizer]. This function will tokenize the text on it's first call,
// then cache the tokens for future calls.
func (t *Text) Tokens() (tokens *tokenlist.TokenList, err error) {
	if t.tokens == nil {
		if t.tokens, err = t.tokenizer.Tokenize(t.text); err != nil {
			return nil, err
		}
	}
	return t.tokens, nil
}

// Returns a [tokenize.TokenList] for the [Text]s token stems using the configured
// [stem.Stemmer]. This function will stem the tokens on it's first call,
// then cache the stems for future calls. If the text has not been tokenized,
// then it will also be tokenized and cached on the first call.
func (t *Text) Stems() (stems *tokenlist.TokenList, err error) {
	if t.stems == nil {
		// Initialize the stems with the tokens
		t.stems = tokenlist.NewCopy(t.tokens)

		// Perform stemming by replacing each token with it's stem
		for i, tok := range t.stems.Tokens() {
			t.stems.Replace(i, token.New(t.stemmer.Stem(tok.String())))
		}
	}
	return t.stems, nil
}

// TODO func (t *Text) Types() *TokenList { return t.Stems() }
// TODO func (t *Text) Len() int {}
// TODO func (t *Text) Runes() []rune {}
// TODO func (t *Text) TokenCount() map[Token]int {}
// TODO func (t *Text) TypeCount() map[Token]int {}
// TODO func (t *Text) CountVectorize() vector.Vector {}
// TODO func (t *Text) CosineSimilarity(other *Text) float64 {}

// ###########################################################################
// Getters
// ###########################################################################

// Returns the string that this text was created with.
func (t *Text) Text() string {
	return t.text
}

// Returns the [enum.Language] configured on this [Text].
func (t *Text) Language() enum.Language {
	return t.lang
}

// Returns the [stem.Stemmer] configured on this [Text].
func (t *Text) Stemmer() stem.Stemmer {
	return t.stemmer
}

// Returns the [tokenize.Tokenizer] configured on this [Text].
func (t *Text) Tokenizer() tokenize.Tokenizer {
	return t.tokenizer
}
