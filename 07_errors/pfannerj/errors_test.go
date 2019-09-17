package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert := assert.New(t)
	err := Errorf(ErrNotFound, "test error message")
	assert.Equal(ErrNotFound, err.Code)
	assert.Equal("test error message", err.Message)
	errWithCause := Errorf(ErrInvalidInput, "test error message with cause %v", err)
	println("***Error with cause***", errWithCause.Code)
	assert.Equal(ErrInvalidInput, errWithCause.Code)
	assert.Equal(errWithCause.Message, "test error message with cause test error message")
	assert.NotEqual(errWithCause.Cause, err)
}
