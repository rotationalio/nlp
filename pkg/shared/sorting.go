package shared

import (
	"slices"
	"strings"
)

// Sorts the slice of strings by length (longest to shortest) and sorts
// lexicographically for equal length strings.
func SortByLengthAndLexicographically(l []string) {
	slices.SortFunc(l, func(a, b string) int {
		if len(a) == len(b) {
			return strings.Compare(a, b)
		}
		return len(b) - len(a)
	})
}
