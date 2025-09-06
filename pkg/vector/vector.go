package vector

// Vector is currently implemented by a `[]float64`.
type Vector []float64

// A [Vector] wrapper for the [Cosine] function.
func (v Vector) Cosine(other Vector) (cosine float64, err error) {
	return Cosine(v, other)
}

// A [Vector] wrapper for the [DotProduct] function.
func (v Vector) DotProduct(other Vector) (product float64, err error) {
	return DotProduct(v, other)
}

// A [Vector] wrapper for the [Magnitude] function.
func (v Vector) Magnitude() (length float64) {
	return Magnitude(v)
}

// A [Vector] wrapper for the [Len] function.
func (v Vector) Len() (elements int) {
	return Len(v)
}
