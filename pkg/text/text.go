/*
[Text] is a one-stop shop for performing NLP operations on text.

Usage example:

	// Create a [Text] with the default settings
	myText := text.New("apple aardvarks zebra bananna aardvark")

	// Get all of the word tokens (ignoring errors in this example)
	myTokens, _ := myText.Tokens() // TokenList

	// Get all word stem tokens which use the same underlying types as the full
	// word tokens above (ignoring errors in this example)
	myStems, _ := myText.Stems() // TokenList

	// The stems are 1:1 count with the tokens
	if myTokens.Len() != myStems.Len() { // 5 == 5
		panic("this should never occur")
	}

	// You can also get a type count, which returns the count of each unique
	// word stem (ignoring errors) ("aardvark" has a 2 count for this example)
	myCount, _ := myText.TypeCount() // map[string]int

	// These are a [tokenlist.TokenList], but if you need a slice of strings...
	stringTokens := myTokens.Strings() // []string

	// Or to get the tokens as a slice of [token.Token] instead
	tokenSlice := myTokens.Tokens() // []Token

	// Get an individual token
	firstToken := tokenSlice[0] // Token

	// You can also get a token as another type
	stringToken := firstToken.String() // string
	runeToken := firstToken.Runes()    // []rune
	byteToken := firstToken.Bytes()    // []byte

	// For these examples, we need to re-create the [Text] with a vocabulary,
	// so the [vectorize.CountVectorizer] will work without an error to get
	// cosine similarity.
	myText := text.New(
	    "cars have engines",
	    text.WithVocabulary([]string{"car", "engine", "brakes", "transmission"}),
	)
	otherText, _ := text.New("engines go with transmissions") // no need vocab

	// Cosine similarity with another string
	similarity, _ := myText.CosineSimilarity() // float64 in range [-1.0, 1.0]

	// We can also get a one-hot or frequency vectorization of our text
	myOneHotVector, _ := myText.VectorizeOneHot()
	myFrequencyVector, _ := myText.VectorizeFrequency()
*/
package text

import (
	"unicode/utf8"

	"go.rtnl.ai/nlp/pkg/language"
	"go.rtnl.ai/nlp/pkg/similarity"
	"go.rtnl.ai/nlp/pkg/stem"
	"go.rtnl.ai/nlp/pkg/token"
	"go.rtnl.ai/nlp/pkg/tokenize"
	"go.rtnl.ai/nlp/pkg/tokenlist"
	"go.rtnl.ai/nlp/pkg/vector"
	"go.rtnl.ai/nlp/pkg/vectorize"
)

// A one-stop-shop for performing NLP operations on a string of text, such as
// tokenization, stemming, vectorization, etc.
type Text struct {
	// The string representation of the text
	text string

	// ==============================
	// Options
	// ==============================

	vocab     []string // used for the [vectorize.CountVectorizer]
	lang      language.Language
	stemmer   stem.Stemmer
	tokenizer tokenize.Tokenizer

	// ==============================
	// Standard Tools
	// ==============================

	typeCounter       *tokenize.TypeCounter
	countVectorizer   *vectorize.CountVectorizer
	cosineSimilarizer *similarity.CosineSimilarizer

	// ==============================
	// Caching (lazy initialization)
	// ==============================

	tokens    *tokenlist.TokenList
	stems     *tokenlist.TokenList
	typecount map[string]int
}

// Create a new [Text] from the input string with the specified [Option]s.
//
// Defaults:
//   - Vocabulary (use [WithVocabulary]): nil (errors will be returned from certain functions if a vocabulary is not added)
//   - Language (use [WithLanguage]): [language.English]
//   - Stemmer (use [WithStemmer]): [stem.Porter2Stemmer]
//   - Tokenizer (use [WithTokenizer]): [tokenize.RegexTokenizer]
func New(t string, options ...Option) (text *Text, err error) {
	// Initialize text
	text = &Text{
		text: t,
	}

	// OPTIONS

	// Set user options
	for _, opt := range options {
		opt(text)
	}

	// Default languge
	if text.lang == language.Unknown {
		text.lang = language.English
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

	// STANDARD TOOLS

	// Initialize the [tokenize.TypeCounter]
	if text.typeCounter == nil {
		if text.typeCounter, err = tokenize.NewTypeCounter(
			tokenize.TypeCounterWithLanguage(text.lang),
			tokenize.TypeCounterWithStemmer(text.stemmer),
			tokenize.TypeCounterWithTokenizer(text.tokenizer),
		); err != nil {
			return nil, err
		}
	}

	// Initialize the [vectorize.CountVectorizer]
	if text.countVectorizer, err = vectorize.NewCountVectorizer(
		vectorize.CountVectorizerWithLang(text.lang),
		vectorize.CountVectorizerWithTokenizer(text.tokenizer),
		vectorize.CountVectorizerWithStemmer(text.stemmer),
		vectorize.CountVectorizerWithTypeCounter(text.typeCounter),
	); err != nil {
		return nil, err
	}

	// Initialize the [similarity.CosineSimilarizer]
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
		t.typecount = t.typeCounter.CountTypes(stems.Strings())
	}
	return t.typecount, nil
}

// VectorizeFrequency returns a frequency (count) encoding vector for the [Text]
// and vocabulary. The vector returned has a value of the count of word
// instances within the chunk for each vocabulary word index.
// NOTE: You must set the vocabulary on the [Text] using [WithVocabulary] during
// creation or an error will be returned.
func (t *Text) VectorizeFrequency() (vector.Vector, error) {
	return t.countVectorizer.VectorizeFrequency(t.text, t.vocab)
}

// VectorizeOneHot returns a one-hot encoding vector for the [Text] and
// vocabulary. The vector returned has a value of 1 for each vocabulary
// word index if it is present within the text and 0 otherwise.
// NOTE: You must set the vocabulary on the [Text] using [WithVocabulary] during
// creation or an error will be returned.
func (t *Text) VectorizeOneHot() (vector.Vector, error) {
	return t.countVectorizer.VectorizeOneHot(t.text, t.vocab)
}

// Retruns a value in the range [-1.0, 1.0] that indicates if two [Text] are
// similar using the cosine similarity method.
// NOTE: If using the [vectorize.CountVectorizer], you must set the vocabulary
// on the [Text] using [WithVocabulary] during creation or an error will be
// returned.
func (t *Text) CosineSimilarity(other *Text) (similarity float64, err error) {
	return t.cosineSimilarizer.Similarity(t.text, other.text)
}

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

// Returns the vocabulary configured on this [Text].
func (t *Text) Vocab() []string {
	return t.vocab
}

// Returns the [language.Language] configured on this [Text].
func (t *Text) Language() language.Language {
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

// Returns the [tokenize.TypeCounter] configured on this [Text].
func (t *Text) TypeCounter() *tokenize.TypeCounter {
	return t.typeCounter
}

// Returns the [similarity.CosineSimilarizer] configured on this [Text].
func (t *Text) CosineSimilarizer() *similarity.CosineSimilarizer {
	return t.cosineSimilarizer
}
