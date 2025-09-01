package vector

import (
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/errors"
	"go.rtnl.ai/nlp/pkg/stemming"
	"go.rtnl.ai/nlp/pkg/tokens"
)

// ############################################################################
// Vectorizer interface
// ############################################################################

type Vectorizer interface {
	Vectorize(chunk string) (vector Vector, err error)
}

// ############################################################################
// CountVectorizer
// ############################################################################

// CountVectorizer can be used to vectorize text using the frequency or one-hot
// text vectorization algorithms.
type CountVectorizer struct {
	vocab       []string
	lang        enum.Language
	tokenizer   tokens.Tokenizer
	stemmer     stemming.Stemmer
	typeCounter *tokens.TypeCounter
	method      VectorizationMethod
}

// Returns a new [CountVectorizer] instance.
//
// Defaults:
//   - Lang: [LanguageEnglish]
//   - Tokenizer: [RegexTokenizer] using Lang above and it's own defaults
//   - Stemmer: [Porter2Stemmer] using Lang above and it's own defaults
//   - TypeCounter: [TypeCounter] using Lang, Stemmer, and Tokenizer above
//   - Method: [VectorizeOneHot]
func NewCountVectorizer(vocab []string, opts ...CountVectorizerOption) (vectorizer *CountVectorizer, err error) {
	// Set options
	vectorizer = &CountVectorizer{}
	for _, fn := range opts {
		fn(vectorizer)
	}

	// Set vocab (a required option)
	vectorizer.vocab = vocab

	// Set defaults

	if vectorizer.lang == enum.LanguageUnknown {
		vectorizer.lang = enum.LanguageEnglish
	}

	if vectorizer.tokenizer == nil {
		vectorizer.tokenizer = tokens.NewRegexTokenizer(tokens.RegexTokenizerWithLanguage(vectorizer.lang))
	}

	if vectorizer.stemmer == nil {
		if vectorizer.stemmer, err = stemming.NewPorter2Stemmer(vectorizer.lang); err != nil {
			return nil, err
		}
	}

	if vectorizer.typeCounter == nil {
		if vectorizer.typeCounter, err = tokens.NewTypeCounter(
			tokens.TypeCounterWithLanguage(vectorizer.lang),
			tokens.TypeCounterWithTokenizer(vectorizer.tokenizer),
			tokens.TypeCounterWithStemmer(vectorizer.stemmer),
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

// Returns the [CountVectorizer]s configured [enum.Language].
func (c *CountVectorizer) Language() enum.Language {
	return c.lang
}

// Returns the [CountVectorizer]s configured [tokens.Tokenizer].
func (c *CountVectorizer) Tokenizer() tokens.Tokenizer {
	return c.tokenizer
}

// Returns the [CountVectorizer]s configured [stemming.Stemmer].
func (c *CountVectorizer) Stemmer() stemming.Stemmer {
	return c.stemmer
}

// Returns the [CountVectorizer]s configured [tokens.TypeCounter].
func (c *CountVectorizer) TypeCounter() *tokens.TypeCounter {
	return c.typeCounter
}

// Returns the [CountVectorizer]s configured [VectorizationMethod].
func (c *CountVectorizer) Method() VectorizationMethod {
	return c.method
}

// Vectorizes the chunk of text.
func (v *CountVectorizer) Vectorize(chunk string) (vector Vector, err error) {
	switch v.method {
	case VectorizeOneHot:
		return v.VectorizeOneHot(chunk)
	case VectorizeFrequency:
		return v.VectorizeFrequency(chunk)
	}
	return nil, errors.ErrMethodNotSupported
}

// VectorizeFrequency returns a frequency (count) encoding vector for the given
// chunk of text and given vocabulary map. The vector returned has a value of
// the count of word instances within the chunk for each vocabulary word index.
func (v *CountVectorizer) VectorizeFrequency(chunk string) (vector Vector, err error) {
	// Type count the chunk
	var types map[string]int
	if types, err = v.typeCounter.TypeCount(chunk); err != nil {
		return nil, err
	}

	// Create the vector from the vocabulary
	vector = make(Vector, len(v.vocab))
	for i, word := range v.vocab {
		// Stem the vocab word with the same stemmer as the type counter uses
		stem := v.typeCounter.Stemmer().Stem(word)
		if count, ok := types[stem]; ok {
			vector[i] = float64(count)
		}
	}

	return vector, nil
}

// VectorizeOneHot returns a one-hot encoding vector for the given chunk of text
// and given vocabulary map. The vector returned has a value of 1 for each
// vocabulary word index if it is present within the chunk of text and 0
// otherwise.
func (v *CountVectorizer) VectorizeOneHot(chunk string) (vector Vector, err error) {
	// Get the frequency encoding
	if vector, err = v.VectorizeFrequency(chunk); err != nil {
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
	VectorizeUnknown = iota
	VectorizeOneHot
	VectorizeFrequency
)

// ############################################################################
// CountVectorizerOption
// ############################################################################

// TypeCounterOption functions modify a [CountVectorizer].
type CountVectorizerOption func(c *CountVectorizer)

// CountVectorizerWithLang sets the [enum.Language] to use with the
// [CountVectorizer].
func CountVectorizerWithLang(lang enum.Language) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.lang = lang
	}
}

// CountVectorizerWithTokenizer sets the [tokens.Tokenizer] to use with the
// [CountVectorizer].
func CountVectorizerWithTokenizer(tokenizer tokens.Tokenizer) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.tokenizer = tokenizer
	}
}

// CountVectorizerWithStemmer sets the [stemming.Stemmer] to use with the
// [CountVectorizer].
func CountVectorizerWithStemmer(stemmer stemming.Stemmer) CountVectorizerOption {
	return func(c *CountVectorizer) {
		c.stemmer = stemmer
	}
}

// CountVectorizerWithTypeCounter sets the [tokens.TypeCounter] to use with the
// [CountVectorizer].
func CountVectorizerWithTypeCounter(typecounter *tokens.TypeCounter) CountVectorizerOption {
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
