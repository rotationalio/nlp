package text

import (
	"unicode/utf8"

	"go.rtnl.ai/nlp/language"
	"go.rtnl.ai/nlp/readability"
	"go.rtnl.ai/nlp/similarity"
	"go.rtnl.ai/nlp/stem"
	"go.rtnl.ai/nlp/token"
	"go.rtnl.ai/nlp/tokenize"
	"go.rtnl.ai/nlp/tokenlist"
	"go.rtnl.ai/nlp/vector"
	"go.rtnl.ai/nlp/vectorize"
)

// ############################################################################
// Structure and Init
// ############################################################################

/*
[Text] is a one-stop shop for performing NLP operations on text.

Usage example:

	// Create a [Text] with the default settings
	myText, err := text.New("apple aardvarks zebra bananna aardvark")

	// Get all of the word tokens
	myTokens, err := myText.Tokens() // TokenList

	// Get all word stem tokens which use the same underlying types as the full
	// word tokens above (ignoring errors in this example)
	myStems, err := myText.Stems() // TokenList

	// The stems are 1:1 count with the tokens
	if len(myTokens) != len(myStems) { // 5 == 5
		panic("this should never occur")
	}

	// You can also get a type count, which returns the count of each unique
	// word stem (ignoring errors) ("aardvark" has a 2 count for this example)
	myCount, err := myText.TypeCount() // map[string]int

	// These are a [tokenlist.TokenList], but if you need a slice of strings...
	stringTokens := myTokens.Strings() // []string

	// You can also use regular slice functions and operations on a [tokenlist.TokenList]
	length := len(myTokens) // 5
	myTokens = append(myTokens, myTokens[0]) // "apple", "aardvarks", "zebra", "bananna", "aardvark", "apple"
	myTokens[0] = myTokens[1] // "aardvarks", "aardvarks", "zebra", "bananna", "aardvark", "apple"

	// Get an individual token
	firstToken := myTokens[0] // Token

	// You can also get a token as another type
	stringToken := firstToken.String() // string
	runeToken := firstToken.Runes() // []rune
	byteToken := firstToken.Bytes() // []byte

	// For these examples, we need to re-create the [Text] with a vocabulary,
	// so the [vectorize.CountVectorizer] will work without an error to get
	// cosine similarity. You could also use a different vectorization method.
	myText, err = text.New(
		"cars have engines like motorcycles have engines",
		text.WithVocabulary([]string{"car", "engine", "brakes", "transmission"}),
	)
	otherText, err := text.New(
		"engines are attached to transmissions",
		text.WithVocabulary([]string{"car", "engine", "brakes", "transmission"}),
	)

	// Cosine similarity with another string
	similarity, err := myText.CosineSimilarity(otherText) // ~0.5

	// We can also get a one-hot or frequency vectorization of our text
	myOneHotVector, err := myText.VectorizeOneHot() // vector.Vector{1, 1, 0, 0}
	myFrequencyVector, err := myText.VectorizeFrequency() // vector.Vector{1, 2, 0, 0}

	// Get readability scores (a score of 0.0 indicates that the word and/or
	// sentence count is zero)
	ease := myText.FleschKincaidReadingEase() // -5.727
	grade := myText.FleschKincaidGradeLevel() // 15.797

	// Get the counts of various things
	count := myText.WordsCount() // 7
	count = myText.SentencesCount() // 1
	count = myText.SyllablesCount() // 17
*/
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

	typeCounter          *tokenize.TypeCounter
	countVectorizer      *vectorize.CountVectorizer
	cosineSimilarizer    *similarity.CosineSimilarizer
	whitespaceTokenizer  *tokenize.WhitespaceTokenizer
	sentenceSegmenter    *tokenize.SentenceSegmenter
	sspSyllableTokenizer *tokenize.SSPSyllableTokenizer

	// ==============================
	// Caching (lazy initialization)
	// ==============================

	tokens    tokenlist.TokenList
	stems     tokenlist.TokenList
	typecount map[string]int
	words     tokenlist.TokenList
	sentences tokenlist.TokenList
	syllables [][]string
}

// Create a new [Text] from the input string with the specified [Option]s. See
// the [Text]s docstring for examples of how to use it.
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
		vectorize.CountVectorizerWithVocab(text.vocab),
		vectorize.CountVectorizerWithLang(text.lang),
		vectorize.CountVectorizerWithTokenizer(text.tokenizer),
		vectorize.CountVectorizerWithStemmer(text.stemmer),
		vectorize.CountVectorizerWithTypeCounter(text.typeCounter),
	); err != nil {
		return nil, err
	}

	// Initialize the [similarity.CosineSimilarizer]
	if text.cosineSimilarizer == nil {
		if text.cosineSimilarizer, err = similarity.NewCosineSimilarizer(
			similarity.CosineSimilarizerWithVocab(text.vocab),
			similarity.CosineSimilarizerWithLanguage(text.lang),
			similarity.CosineSimilarizerWithTokenizer(text.tokenizer),
			similarity.CosineSimilarizerWithVectorizer(text.countVectorizer),
		); err != nil {
			return nil, err
		}
	}

	// Initialize the [tokenize.WhitespaceTokenizer]
	if text.whitespaceTokenizer == nil {
		text.whitespaceTokenizer = tokenize.NewWhitespaceTokenizer()
	}

	// Initialize the [tokenize.SentenceSegmenter]
	if text.sentenceSegmenter == nil {
		text.sentenceSegmenter = tokenize.NewSentenceSegmenter(
			tokenize.SentenceSegmenterWithLanguage(text.lang),
		)
	}

	// Initialize the [tokenize.SSPSyllableTokenizer]
	if text.sspSyllableTokenizer == nil {
		if text.sspSyllableTokenizer, err = tokenize.NewSSPSyllableTokenizer(text.lang); err != nil {
			return nil, err
		}
	}

	return text, nil
}

// ############################################################################
// Tokenize
// ############################################################################

// Returns a [tokenlist.TokenList] for the [Text]s tokens using the configured
// [tokenize.Tokenizer]. This function cache the result of the operation for
// subsequent calls.
func (t *Text) Tokens() (tokens tokenlist.TokenList, err error) {
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
func (t *Text) Stems() (stems tokenlist.TokenList, err error) {
	if t.stems == nil {
		// Initialize the stems with the tokens
		var tokens tokenlist.TokenList
		if tokens, err = t.Tokens(); err != nil {
			return nil, err
		}
		t.stems = tokenlist.NewCopy(tokens)

		// Perform stemming by replacing each token with it's stem
		for i, tok := range t.stems {
			t.stems[i] = token.New(t.stemmer.Stem(tok.String()))
		}
	}
	return t.stems, nil
}

// Returns the words in the [Text] as a [tokenlist.TokenList]. Cached for faster
// subsequent calls.
func (t *Text) Words() tokenlist.TokenList {
	if t.words == nil {
		words, _ := t.whitespaceTokenizer.Tokenize(t.text) // error is ALWAYS nil
		for _, word := range words {
			t.words = append(t.words, token.New(word))
		}
	}
	return t.words
}

// Returns the sentences in the [Text] as a [tokenlist.TokenList]. Cached for
// faster subsequent calls.
func (t *Text) Sentences() tokenlist.TokenList {
	if t.sentences == nil {
		sentences, _ := t.sentenceSegmenter.Tokenize(t.text) // error is ALWAYS nil
		for _, sentence := range sentences {
			t.sentences = append(t.sentences, token.New(sentence))
		}
	}
	return t.sentences
}

// Returns the words in the [Text] tokenized as syllables as a slice of string
// slices. Cached for faster subsequent calls.
func (t *Text) Syllables() [][]string {
	if t.syllables == nil {
		t.syllables = make([][]string, 0, t.WordCount())
		for _, word := range t.Words().Strings() {
			wordSyllables, _ := t.sspSyllableTokenizer.Tokenize(word) // error is ALWAYS nil
			t.syllables = append(t.syllables, wordSyllables)
		}
	}
	return t.syllables
}

// ############################################################################
// Count
// ############################################################################

// Returns the count of the words in the [Text].
func (t *Text) WordCount() int {
	return len(t.Words()) // Words is cached
}

// Returns the count of the sentences in the [Text].
func (t *Text) SentenceCount() int {
	return len(t.Sentences()) // Sentences is cached
}

// Returns the count of the syllables in the [Text].
func (t *Text) SyllableCount() int {
	count := 0
	for _, wordSyllables := range t.Syllables() { // Syllables is cached
		count += len(wordSyllables)
	}
	return count
}

// Returns a map of the types (unique word stems) and their counts for this
// [Text]. This function cache the result of the operation for subsequent calls.
func (t *Text) TypeCount() (types map[string]int, err error) {
	if t.typecount == nil {
		// Stem the words
		var stems tokenlist.TokenList
		if stems, err = t.Stems(); err != nil {
			return nil, err
		}

		// Count the stems to get the type count
		t.typecount = t.typeCounter.CountTypes(stems.Strings())
	}
	return t.typecount, nil
}

// ############################################################################
// Vectorize
// ############################################################################

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
// Readability
// ###########################################################################

// Returns the Flesch-Kincaid Reading Ease score. Returns the value 0.0 when the
// sentence and/or word count is zero.
func (t *Text) FleschKincaidReadingEase() (score float64) {
	return readability.FleschKincaidReadingEase(t.WordCount(), t.SentenceCount(), t.SyllableCount())
}

// Returns the Flesch-Kincaid grade level. Returns the value 0.0 when the
// sentence and/or word count is zero.
func (t *Text) FleschKincaidGradeLevel() (score float64) {
	return readability.FleschKincaidGradeLevel(t.WordCount(), t.SentenceCount(), t.SyllableCount())
}

// ###########################################################################
// Misc. Properties
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

// Returns the [tokenize.WhitespaceTokenizer] configured on this [Text].
func (t *Text) WhitespaceTokenizer() *tokenize.WhitespaceTokenizer {
	return t.whitespaceTokenizer
}

// Returns the [tokenize.SentenceSegmenter] configured on this [Text].
func (t *Text) SentenceSegmenter() *tokenize.SentenceSegmenter {
	return t.sentenceSegmenter
}

// Returns the [tokenize.SSPSyllableTokenizer] configured on this [Text].
func (t *Text) SSPSyllableTokenizer() *tokenize.SSPSyllableTokenizer {
	return t.sspSyllableTokenizer
}

// ############################################################################
// Cache Getters
// ############################################################################

// Returns the raw cache value for this [Text]s tokens.
func (t *Text) TokensCache() (tokens tokenlist.TokenList) {
	return t.tokens
}

// Returns the raw cache value for this [Text]s stems.
func (t *Text) StemsCache() (stems tokenlist.TokenList) {
	return t.stems
}

// Returns the raw cache value for this [Text]s type count.
func (t *Text) TypeCountCache() (types map[string]int) {
	return t.typecount
}

// Returns the raw cache value for this [Text]s words.
func (t *Text) WordsCache() (words tokenlist.TokenList) {
	return t.words
}

// Returns the raw cache value for this [Text]s sentences.
func (t *Text) SentencesCache() (sentences tokenlist.TokenList) {
	return t.sentences
}

// Returns the raw cache value for this [Text]s syllables.
func (t *Text) SyllablesCache() (syllables [][]string) {
	return t.syllables
}
