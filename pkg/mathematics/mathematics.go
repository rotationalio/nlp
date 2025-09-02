package mathematics

import "math"

// Returns x bound to the closed range `[lower, upper]`.
func BoundToRange(x, lower, upper float64) float64 {
	return math.Max(lower, math.Min(x, upper))
}
