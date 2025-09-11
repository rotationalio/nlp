package language

import (
	"slices"
)

// A lightweight enumeration type for a written language.
type Language uint16

const (
	Unknown Language = iota
	English
)

// Returns True if the argument [enum.Language]s contains this language.
func (l Language) In(langs ...Language) bool {
	return slices.Contains(langs, l)
}
