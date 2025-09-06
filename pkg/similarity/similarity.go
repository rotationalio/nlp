package similarity

import (
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/mathematics"
	"go.rtnl.ai/nlp/pkg/tokenize"
	"go.rtnl.ai/nlp/pkg/vector"
	"go.rtnl.ai/nlp/pkg/vectorize"
)

// ############################################################################
// Similarizer interface
// ############################################################################

// A Similarizer compares the similarity of two strings.
type Similarizer interface {
	Similarity(a, b string) (similarity float64, err error)
}

// ############################################################################
// CosineSimilarizer
// ############################################################################

// CosineSimilarizer can be used to calculate the cosine similarity of two text
// strings using the cosine of their vectors.
type CosineSimilarizer struct {
	vocab      []string
	lang       enum.Language
	tokenizer  tokenize.Tokenizer
	vectorizer vectorize.Vectorizer
}

// Returns a new [CosineSimilarizer] with the vocabulary and options set.
//
// Defaults:
//   - Vocab: nil
//   - Lang: [enum.LanguageEnglish]
//   - Tokenizer: [tokenize.RegexTokenizer]
//   - Vectorizer: [vectorize.CountVectorizer]
func NewCosineSimilarizer(opts ...CosineSimilarizerOption) (similarizer *CosineSimilarizer, err error) {
	// Set options
	similarizer = &CosineSimilarizer{}
	for _, fn := range opts {
		fn(similarizer)
	}

	//Set defaults

	if similarizer.lang == enum.LanguageUnknown {
		similarizer.lang = enum.LanguageEnglish
	}

	if similarizer.tokenizer == nil {
		similarizer.tokenizer = tokenize.NewRegexTokenizer(tokenize.RegexTokenizerWithLanguage(similarizer.lang))
	}

	if similarizer.vectorizer == nil {
		if similarizer.vectorizer, err = vectorize.NewCountVectorizer(
			vectorize.CountVectorizerWithVocab(similarizer.vocab),
			vectorize.CountVectorizerWithLang(similarizer.lang),
		); err != nil {
			return nil, err
		}
	}

	return similarizer, nil
}

// Returns the [CosineSimilarizer]s configured vocabulary.
func (c *CosineSimilarizer) Vocab() []string {
	return c.vocab
}

// Returns the [CosineSimilarizer]s configured [enum.Language].
func (c *CosineSimilarizer) Language() enum.Language {
	return c.lang
}

// Returns the [CosineSimilarizer]s configured [tokenize.Tokenizer].
func (c *CosineSimilarizer) Tokenizer() tokenize.Tokenizer {
	return c.tokenizer
}

// Returns the [CosineSimilarizer]s configured [vector.Vectorizer].
func (c *CosineSimilarizer) Vectorizer() vectorize.Vectorizer {
	return c.vectorizer
}

// Similarity returns a value in the range [-1.0, 1.0] that indicates if two
// strings are similar using the cosine similarity method. Must have set a
// vocabulary using [With] in [NewCosineSimilarizer] to use this function.
func (s *CosineSimilarizer) Similarity(a, b string) (similarity float64, err error) {
	//Vectorize the strings
	var vecA, vecB vector.Vector
	if vecA, err = s.vectorizer.Vectorize(a); err != nil {
		return 0.0, err
	}
	if vecB, err = s.vectorizer.Vectorize(b); err != nil {
		return 0.0, err
	}

	// Calculate the cosine of the angle between the vectors as the similarity
	if similarity, err = vector.Cosine(vecA, vecB); err != nil {
		return 0.0, err
	}

	// Return value bounded to [-1.0, 1.0]
	return mathematics.BoundToRange(similarity, -1.0, 1.0), nil
}

// ############################################################################
// SimilarityOption
// ############################################################################

// A CosineSimilarizerOption function sets options for a [CosineSimilarizer].
type CosineSimilarizerOption func(s *CosineSimilarizer)

// Returns a function which sets a [CosineSimilarizer]s vocabulary.
func CosineSimilarizerWithVocab(vocab []string) CosineSimilarizerOption {
	return func(s *CosineSimilarizer) {
		s.vocab = vocab
	}
}

// Returns a function which sets a [CosineSimilarizer]s [enum.Language].
func CosineSimilarizerWithLanguage(lang enum.Language) CosineSimilarizerOption {
	return func(s *CosineSimilarizer) {
		s.lang = lang
	}
}

// Returns a function which sets a [CosineSimilarizer]s [tokenize.Tokenizer].
func CosineSimilarizerWithTokenizer(tokenizer tokenize.Tokenizer) CosineSimilarizerOption {
	return func(s *CosineSimilarizer) {
		s.tokenizer = tokenizer
	}
}

// Returns a function which sets a [CosineSimilarizer]s [vector.Vectorizer].
func CosineSimilarizerWithVectorizer(vectorizer vectorize.Vectorizer) CosineSimilarizerOption {
	return func(s *CosineSimilarizer) {
		s.vectorizer = vectorizer
	}
}
