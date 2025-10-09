package vectorize

import (
	"go.rtnl.ai/nlp/vector"
)

// ############################################################################
// Vectorizer interface
// ############################################################################

type Vectorizer interface {
	Vectorize(chunk string) (vector vector.Vector, err error)
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
