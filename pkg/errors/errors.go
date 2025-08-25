package errors

import "errors"

var (
	ErrUnequalLengthVectors = errors.New("vector arguments must have an equal number of elements")
	ErrLanguageNotSupported = errors.New("the selected language is not supported")
	ErrMethodNotSupported   = errors.New("the selected method is not supported")
	ErrUndefinedValue       = errors.New("the mathematical operation has no defined value for the given arugments")
)
