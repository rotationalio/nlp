package vector

import (
	"math"

	"go.rtnl.ai/nlp/pkg/errors"
)

// Cosine returns the cosine of the angle between two vectors as a value between
// [-1.0, 1.0], as defined by SLP 3rd Edition section 6.4 fig 6.10. If the
// vectors do not have the same number of elements or either of the vectors has
// a length of zero, an error will be returned.
func Cosine(a, b Vector) (cosine float64, err error) {
	// Ensure vectors have the same number of elements
	if len(a) != len(b) {
		return 0.0, errors.ErrUnequalLengthVectors
	}

	// Calculate the dot product
	var dotprod, vlenprod float64
	if dotprod, err = DotProduct(a, b); err != nil {
		return 0.0, err
	}

	// Calculate the product of the two vector's lengths
	vlenprod = VectorLength(a) * VectorLength(b)
	if vlenprod == 0.0 {
		// Cosine is undefined for zero length vectors
		return 0.0, errors.ErrUndefinedValue
	}

	// Return final cosine value clamped to [-1.0, 1.0]
	return math.Max(-1.0, math.Min(dotprod/vlenprod, 1.0)), nil
}

// DotProduct returns the dot product of the two vectors (as defined by SLP 3rd
// Edition section 6.4 fig 6.7). If the vectors do not have the same number
// of elements, an error will be returned.
func DotProduct(a, b Vector) (product float64, err error) {
	// Ensure vectors have the same number of elements
	if len(a) != len(b) {
		return 0.0, errors.ErrUnequalLengthVectors
	}

	for i := range a {
		product += a[i] * b[i]
	}
	return product, nil
}

// VectorLength returns the vector length (as defined by SLP 3rd Edition section
// 6.4 fig 6.8).
func VectorLength(v Vector) (length float64) {
	for _, e := range v {
		length += e * e
	}
	return math.Sqrt(length)
}
