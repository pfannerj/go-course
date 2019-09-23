package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert := assert.New(t)
	errWithoutCause := Errorf(ErrNotFound, "test error message")
	assert.Equal(ErrNotFound, errWithoutCause.Code)
	assert.Equal("test error message", errWithoutCause.Message)
	errMsg1 := errWithoutCause.Error()
	assert.Equal("test error message\n<nil>", errMsg1)
	errWithCause := Errorf(ErrInvalidInput, "test error message with cause %d", ErrInternalDataError)
	fmt.Println("*** Before Error call ***", errWithCause.Cause)
	assert.Error(errWithCause)
	errMsg2 := errWithCause.Error()
	assert.Equal("test error message with cause 2\n - error cause: [2]", errMsg2)
	fmt.Println("*** After Error call ***", errWithCause.Cause)
	assert.Equal(errWithCause.Code, ErrInvalidInput)
	assert.Equal(errWithCause.Message, "test error message with cause 2")
	assert.Equal(errWithCause.Cause, fmt.Errorf(" - error cause: [%d]", ErrInternalDataError))
}

// func TestErrorf(t *testing.T) {
// 	err := Errorf(ErrInternalDataError, "internal data error")
// 	assert.Equal(t, ErrInternalDataError, err.Code)
// 	errMessage := err.Error()
// 	assert.Equal(t, "internal data error\n<nil>", errMessage)
// }
