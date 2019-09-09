package main

import (
	"fmt"
)

// Error codes
const (
	// ErrNotFound is used when attempting to read a non-existing entry
	ErrNotFound int = iota
	// ErrInvalidInput is used when the incoming request is invalid
	ErrInvalidInput
	// ErrInternalDataError is used for issues retrieving data from the DB
	ErrInternalDataError
)

// Error defines an error that separates internal and external error messages
type Error struct {
	Message string
	Code    int
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%v\n%v", e.Message, e.Cause)
}

// Errorf creates a new Error with formatting
func Errorf(code int, format string, args ...interface{}) *Error {
	return ErrorEf(code, nil, format, args...)
}

// ErrorEf creates a new Error with causing error and formatting
func ErrorEf(code int, cause error, format string, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Code:    code,
		Cause:   cause,
	}
}
