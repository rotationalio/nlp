package vectorize

import (
	"go.rtnl.ai/nlp/pkg/errors"
	"go.rtnl.ai/nlp/pkg/language"
	"go.rtnl.ai/nlp/pkg/stem"
	"go.rtnl.ai/nlp/pkg/tokenize"
	"go.rtnl.ai/nlp/pkg/vector"
)

// ############################################################################
// Vectorizer interface
// ############################################################################

type Vectorizer interface {
	Vectorize(chunk string) (vector vector.Vector, err error)
}

// ############################################################################
// CountVectorizer
// ############################################################################

// CountVectorizer can be used to vectorize text using the frequency or one-hot
// text vectorization algorithms.
type CountVectorizer struct {
	vocab       []string
	lang        language.Language
	tokenizer   tokenize.Tokenizer
	stemmer     stem.Stemmer
	typeCounter *tokenize.TypeCounter
	method      VectorizationMethod
}

// Returns a new [CountVectorizer] instance.
//
// Defaults:
//   - Vocab: nil
//   - Lang: [language.English]
//   - Tokenizer: [tokenize.RegexTokenizer]
//   - Stemmer: [stem.Porter2Stemmer]
//   - TypeCounter: [tokenize.TypeCounter]
//   - Method: [VectorizeOneHot]
func NewCountVectorizer(opts ...CountVectorizerOption) (vectorizer *CountVectorizer, err error) {
	// Set options
	vectorizer = &CountVectorizer{}
	for _, fn := range opts {
		fn(vectorizer)
	}

	// Set defaults

	if vectorizer.lang == language.Unknown {
		vectorizer.lang = language.English
	}

	if vectorizer.tokenizer == nil {
		vectorizer.tokenizer = tokenize.NewRegexTokenizer(tokenize.RegexTokenizerWithLanguage(vectorizer.lang))
	}

	if vectorizer.stemmer == nil {
		if vectorizer.stemmer, err = stem.NewPorter2Stemmer(vectorizer.lang); err != nil {
			return nil, err
		}
	}

	if vectorizer.typeCounter == nil {
		if vectorizer.typeCounter, err = tokenize.NewTypeCounter(
			tokenize.TypeCounterWithLanguage(vectorizer.lang),
			tokenize.TypeCounterWithTokenizer(vectorizer.tokenizer),
			tokenize.TypeCounterWithStemmer(vectorizer.stemmer),
		); err != nil {
			return nil, err
		}
	}

	if vectorizer.method == VectorizeUnknown {
		vectorizer.method = VectorizeOneHot
	}

	return vectorizer, nil
}

// Returns the [CountVectorizer]s configured vocabulary.
func (c *CountVectorizer) Vocab() []string {
	return c.vocab
}

// Returns the [CountVectorizer]s configured [language.Language].
func (c *CountVectorizer) Language() language.Language {
	return c.lang
}

// Returns the [CountVectorizer]s configured [tokenize.Tokenizer].
func (c *CountVectorizer) Tokenizer() tokenize.Tokenizer {
	return c.tokenizer
}

// Returns the [CountVectorizer]s configured [stem.Stemmer].
func (c *CountVectorizer) Stemmer() stem.Stemmer {
	return c.stemmer
}

// Returns the [CountVectorizer]s configured [tokenize.TypeCounter].
func (c *CountVectorizer) TypeCounter() *tokenize.TypeCounter {
	return c.typeCounter
}

// Returns the [CountVectorizer]s configured [VectorizationMethod].
func (c *CountVectorizer) Method() VectorizationMethod {
	return c.method
}

// Vectorizes the chunk of text using the pre-configured vocabulary and
// [VectorizationMethod].
func (v *CountVectorizer) Vectorize(chunk string) (vector vector.Vector, err error) {
	// We need to have set a vocabulary if we wish to use this function
	if v.vocab == nil {
		return nil, errors.ErrVocabularyNotSet
	}

	// Call the method function
	switch v.method {
	case VectorizeOneHot:
		return v.VectorizeOneHot(chunk, v.vocab)
	case VectorizeFrequency:
		return v.VectorizeFrequency(chunk, v.vocab)
	}
	return nil, errors.ErrMethodNotSupported
}

// VectorizeFrequency returns a frequency (count) encoding vector for the given
// chunk of text and given vocabulary. The vector returned has a value of
// the count of word instances within the chunk for each vocabulary word index.
func (v *CountVectorizer) VectorizeFrequency(chunk string, vocab []string) (vector vector.Vector, err error) {
	// Type count the text
	var types map[string]int
	if types, err = v.typeCounter.TypeCount(chunk); err != nil {
		return nil, err
	}

	// Create the vector from the vocabulary
	vector = make([]float64, len(vocab))
	for i, word := range vocab {
		// Stem the vocab word with the same stemmer as the type counter uses
		stem := v.typeCounter.Stemmer().Stem(word)
		if count, ok := types[stem]; ok {
			vector[i] = float64(count)
		}
	}

	return vector, nil
}

// VectorizeOneHot returns a one-hot encoding vector for the given text chunk
// and given vocabulary. The vector returned has a value of 1 for each
// vocabulary word index if it is present within the text and 0 otherwise.
func (v *CountVectorizer) VectorizeOneHot(chunk string, vocab []string) (vector vector.Vector, err error) {
	// Get the frequency encoding
	if vector, err = v.VectorizeFrequency(chunk, vocab); err != nil {
		return nil, err
	}

	// Then convert it to a one-hot encoding
	for i, e := range vector {
		if e != 0.0 {
			vector[i] = 1
		}
	}

	return vector, nil
}

// ############################################################################
// VectorizationMethod "enum"
// ############################################################################

type VectorizationMethod uint8

const (
	VectorizeUnknown VectorizationMethod = iota
	VectorizeOneHot
	VectorizeFrequency
)

// ############################################################################
// CountVectorizerOption
// ############################################################################

// TypeCounterOption functions modify a [CountVectorizer].
type CountVectorizerOption func(c *CountVectorizer)

// CountVectorizerWithLang sets the vocabulary to use with the
// [CountVectorizer].
func CountVectorizerWithVocab(vocab []string) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.vocab = vocab
	}
}

// CountVectorizerWithLang sets the [language.Language] to use with the
// [CountVectorizer].
func CountVectorizerWithLang(lang language.Language) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.lang = lang
	}
}

// CountVectorizerWithTokenizer sets the [tokenize.Tokenizer] to use with the
// [CountVectorizer].
func CountVectorizerWithTokenizer(tokenizer tokenize.Tokenizer) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.tokenizer = tokenizer
	}
}

// CountVectorizerWithStemmer sets the [stem.Stemmer] to use with the
// [CountVectorizer].
func CountVectorizerWithStemmer(stemmer stem.Stemmer) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.stemmer = stemmer
	}
}

// CountVectorizerWithTypeCounter sets the [tokenize.TypeCounter] to use with the
// [CountVectorizer].
func CountVectorizerWithTypeCounter(typecounter *tokenize.TypeCounter) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.typeCounter = typecounter
	}
}

// CountVectorizerWithMethod sets the [VectorizationMethod] to use with the
// [CountVectorizer].
func CountVectorizerWithMethod(method VectorizationMethod) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.method = method
	}
}
