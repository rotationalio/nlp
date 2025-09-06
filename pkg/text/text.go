package text

import (
	"unicode/utf8"

	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/similarity"
	"go.rtnl.ai/nlp/pkg/stem"
	"go.rtnl.ai/nlp/pkg/token"
	"go.rtnl.ai/nlp/pkg/tokenize"
	"go.rtnl.ai/nlp/pkg/tokenlist"
	"go.rtnl.ai/nlp/pkg/vector"
	"go.rtnl.ai/nlp/pkg/vectorize"
)

// A one-stop-shop for performing NLP operations on a string of text, such as
// stemming, vectorization, or similarity to another string.
type Text struct {
	// The string representation of the text.
	text string

	// ========================
	// Options
	// ========================

	// The [enum.Language] of this text.
	lang enum.Language

	// The [stem.Stemmer] to use for stemming this text's tokens.
	stemmer stem.Stemmer

	// The [tokenize.Tokenizer] to use for tokenizing this text.
	tokenizer tokenize.Tokenizer

	// The [tokenize.TypeCounter] to use for tokenizing this text.
	counter *tokenize.TypeCounter

	// ========================
	// Standard Tools
	// ========================

	// FrequencyVectorizer
	countVectorizer *vectorize.CountVectorizer

	// CosineSimilarizer
	cosineSimilarizer *similarity.CosineSimilarizer

	// ========================
	// Caching
	// ========================

	// Cache of tokens of this text; lazily initialized.
	tokens *tokenlist.TokenList

	// Cache of stem tokens of this text; lazily initialized.
	stems *tokenlist.TokenList

	// Cache of type count of this text; lazily initialized.
	typecount map[string]int
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

	// Default type counter
	if text.counter == nil {
		if text.counter, err = tokenize.NewTypeCounter(
			tokenize.TypeCounterWithLanguage(text.lang),
			tokenize.TypeCounterWithStemmer(text.stemmer),
			tokenize.TypeCounterWithTokenizer(text.tokenizer),
		); err != nil {
			return nil, err
		}
	}

	// Initialize the CountVectorizer
	if text.countVectorizer, err = vectorize.NewCountVectorizer(
		vectorize.CountVectorizerWithLang(text.lang),
		vectorize.CountVectorizerWithTokenizer(text.tokenizer),
		vectorize.CountVectorizerWithStemmer(text.stemmer),
		vectorize.CountVectorizerWithTypeCounter(text.counter),
	); err != nil {
		return nil, err
	}

	if text.cosineSimilarizer, err = similarity.NewCosineSimilarizer(
		similarity.CosineSimilarizerWithLanguage(text.lang),
		similarity.CosineSimilarizerWithTokenizer(text.tokenizer),
		similarity.CosineSimilarizerWithVectorizer(text.countVectorizer),
	); err != nil {
		return nil, err
	}

	return text, nil
}

// Returns a [tokenlist.TokenList] for the [Text]s tokens using the configured
// [tokenize.Tokenizer]. This function cache the result of the operation for
// subsequent calls.
func (t *Text) Tokens() (tokens *tokenlist.TokenList, err error) {
	if t.tokens == nil {
		var toks []string
		if toks, err = t.tokenizer.Tokenize(t.text); err != nil {
			return nil, err
		}
		t.tokens = tokenlist.New(toks)
	}
	return t.tokens, nil
}

// Returns a [tokenlist.TokenList] for the [Text]s stems using the configured
// [stem.Stemmer]. This function cache the result of the operation for
// subsequent calls.
func (t *Text) Stems() (stems *tokenlist.TokenList, err error) {
	if t.stems == nil {
		// Initialize the stems with the tokens
		var tokens *tokenlist.TokenList
		if tokens, err = t.Tokens(); err != nil {
			return nil, err
		}
		t.stems = tokenlist.NewCopy(tokens)

		// Perform stemming by replacing each token with it's stem
		for i, tok := range t.stems.Tokens() {
			t.stems.Replace(i, token.New(t.stemmer.Stem(tok.String())))
		}
	}
	return t.stems, nil
}

// Returns a map of the types (unique word stems) and their counts for this
// [Text]. This function cache the result of the operation for subsequent calls.
func (t *Text) TypeCount() (types map[string]int, err error) {
	if t.typecount == nil {
		// Stem the words
		var stems *tokenlist.TokenList
		if stems, err = t.Stems(); err != nil {
			return nil, err
		}
		// Count the stems to get the type count
		t.typecount = t.counter.CountTypes(stems.Strings())
	}
	return t.typecount, nil
}

// VectorizeFrequency returns a frequency (count) encoding vector for the [Text]
// and vocabulary. The vector returned has a value of the count of word
// instances within the chunk for each vocabulary word index.
// TODO (sc-34048): replace the vocab with a vocab.Vocab that is storable and etc.
func (t *Text) VectorizeFrequency(vocab []string) (vector.Vector, error) {
	return t.countVectorizer.VectorizeFrequency(t.text, vocab)
}

// VectorizeFrequency returns a frequency (count) encoding vector for the [Text]
// and vocabulary. The vector returned has a value of 1 for each vocabulary
// word index if it is present within the text and 0 otherwise.
// TODO (sc-34048): replace the vocab with a vocab.Vocab that is storable and etc.
func (t *Text) VectorizeOneHot(vocab []string) (vector.Vector, error) {
	return t.countVectorizer.VectorizeOneHot(t.text, vocab)
}

// FIXME: we have a vocabulary problem :()
//TODO func (t *Text) CosineSimilarity(other *Text) float64 {}

// ###########################################################################
// Properties
// ###########################################################################

// Returns the number of UTF-8 runes (aka: characters) in the [Text].
func (t *Text) Len() int {
	return utf8.RuneCountInString(t.text)
}

// Returns the number of bytes in the [Text] (like `len(aString)`).
func (t *Text) ByteLen() int {
	return len(t.text)
}

// ###########################################################################
// Getters
// ###########################################################################

// Returns the [Text] as a string.
func (t *Text) Text() string {
	return t.text
}

// Returns the [Text] as a string.
func (t *Text) String() string {
	return t.text
}

// Returns the [Text] as a slice of runes.
func (t *Text) Runes() []rune {
	return []rune(t.text)
}

// Returns the [Text] as a slice of bytes.
func (t *Text) Bytes() []byte {
	return []byte(t.text)
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

// Returns the [vectorize.CountVectorizer] configured on this [Text].
func (t *Text) CountVectorizer() *vectorize.CountVectorizer {
	return t.countVectorizer
}

// Returns the [similarity.CosineSimilarizer] configured on this [Text].
func (t *Text) CosineSimilarizer() *similarity.CosineSimilarizer {
	return t.cosineSimilarizer
}
