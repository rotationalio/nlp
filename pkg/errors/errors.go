package errors

import "errors"

var (
	ErrInvalidIndex         = errors.New("the value is not a valid index for this type")
	ErrLanguageNotSupported = errors.New("the selected language is not supported")
	ErrMethodNotSupported   = errors.New("the selected method is not supported")
	ErrUndefinedValue       = errors.New("the mathematical operation has no defined value for the given arugments")
	ErrUnequalLengthVectors = errors.New("vector arguments must have an equal number of elements")
	ErrVocabularyNotSet     = errors.New("this operation requires a vocabulary be set")
)
