package errors

import "errors"

var (
	ErrInvalidIndex         = errors.New("the value is not a valid index for this type")
	ErrLanguageNotSupported = errors.New("the selected language is not supported")
	ErrMethodNotSupported   = errors.New("the selected method is not supported")
	ErrMissingConfig        = errors.New("missing a required configuration value")
	ErrUndefinedValue       = errors.New("the mathematical operation has no defined value for the given arugments")
	ErrUnequalLengthVectors = errors.New("vector arguments must have an equal number of elements")
)

// Call to stdlib's [errors.New]:
//
// New returns an error that formats as the
// given text. Each call to New returns a distinct error value even if the text
// is identical.
var New func(text string) error = errors.New

// Call to stdlib's [errors.Join]:
//
// Join returns an error that wraps the given errors. Any nil error values are
// discarded. Join returns nil if every value in errs is nil. The error formats
// as the concatenation of the strings obtained by calling the Error method of
// each element of errs, with a newline between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
var Join func(errs ...error) error = errors.Join
