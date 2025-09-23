package language

import (
	"slices"
)

// ############################################################################
// [Language] Enumeration (lightweight)
// ############################################################################

// A lightweight enumeration type for a written language.
type Language uint16

const (
	Unknown Language = iota
	English
)

// ############################################################################
// Methods
// ############################################################################

// Returns True if the argument [enum.Language]s contains this language.
func (l Language) In(langs ...Language) bool {
	return slices.Contains(langs, l)
}

// ############################################################################
// Constants for [English] language
// ############################################################################

// Punctuation marks which end an English sentence.
var SENTENCE_PUNCTUATION_ENGLISH = ".!?"

// Words that look like they end a sentence but usually do not. Non-exhaustive.
var SENTENCE_STOP_WORDS_ENGLISH = []string{"Mr.", "Mrs.", "Ms.", "Dr.", "Hon."}
