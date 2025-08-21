package similarity

import (
	"go.rtnl.ai/nlp/pkg/enum"
	"go.rtnl.ai/nlp/pkg/tokens"
	"go.rtnl.ai/nlp/pkg/vector"
)

// ############################################################################
// Similarizer interface
// ############################################################################

type Similarizer interface {
	Similarity(a, b string) (similarity float64, err error)
}

// ############################################################################
// CosineSimilarizer
// ############################################################################

// CosineSimilarizer can be used to calculate the cosine similarity of two text
// chunks.
type CosineSimilarizer struct {
	vocab      []string
	lang       enum.Language
	tokenizer  tokens.Tokenizer
	vectorizer vector.Vectorizer
}

// Returns a new [CosineSimilarizer] with the vocabulary and options set.
//
// Defaults:
//   - Lang: [LanguageEnglish]
//   - Tokenizer: [RegexTokenizer] with the Lang above and it's own defaults
//   - Vectorizer: [CountVectorizer] with the given vocabulary, Lang above, and
//     it's own defaults
func NewCosineSimilarizer(vocab []string, opts ...CosineSimilarizerOption) (similarizer *CosineSimilarizer, err error) {
	// Set options
	similarizer = &CosineSimilarizer{}
	for _, fn := range opts {
		fn(similarizer)
	}

	// Set vocab (a required option)
	similarizer.vocab = vocab

	//Set defaults
	if similarizer.lang == enum.LanguageUnknown {
		similarizer.lang = enum.LanguageEnglish
	}
	if similarizer.tokenizer == nil {
		similarizer.tokenizer = tokens.NewRegexTokenizer(tokens.RegexTokenizerWithLanguage(similarizer.lang))
	}
	if similarizer.vectorizer == nil {
		if similarizer.vectorizer, err = vector.NewCountVectorizer(similarizer.vocab, vector.CountVectorizerWithLang(similarizer.lang)); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// Returns the [CosineSimilarizer]s configured vocabulary.
func (c *CosineSimilarizer) Vocab() []string {
	return c.vocab
}

// Returns the [CosineSimilarizer]s configured [enum.Language].
func (c *CosineSimilarizer) Language() enum.Language {
	return c.lang
}

// Returns the [CosineSimilarizer]s configured [tokens.Tokenizer].
func (c *CosineSimilarizer) Tokenizer() tokens.Tokenizer {
	return c.tokenizer
}

// Returns the [CosineSimilarizer]s configured [vector.Vectorizer].
func (c *CosineSimilarizer) Vectorizer() vector.Vectorizer {
	return c.vectorizer
}

// Similarity returns a value in the range [-1.0, 1.0] that indicates if two
// strings are similar using the cosine similarity method.
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

	return similarity, nil
}

// ############################################################################
// SimilarityOption
// ############################################################################

// A CosineSimilarizerOption function sets options for a [CosineSimilarizer].
type CosineSimilarizerOption func(s *CosineSimilarizer)

// Returns a function which sets a [CosineSimilarizer]s [enum.Language].
func CosineSimilarizerWithLanguage(lang enum.Language) CosineSimilarizerOption {
	return func(s *CosineSimilarizer) {
		s.lang = lang
	}
}

// Returns a function which sets a [CosineSimilarizer]s [tokens.Tokenizer].
func CosineSimilarizerWithTokenizer(tokenizer tokens.Tokenizer) CosineSimilarizerOption {
	return func(s *CosineSimilarizer) {
		s.tokenizer = tokenizer
	}
}

// Returns a function which sets a [CosineSimilarizer]s [vector.Vectorizer].
func CosineSimilarizerWithVectorizer(vectorizer vector.Vectorizer) CosineSimilarizerOption {
	return func(s *CosineSimilarizer) {
		s.vectorizer = vectorizer
	}
}
