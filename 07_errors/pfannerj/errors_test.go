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
	errFull := Errorf(ErrNotFound, "test error message %v", err)
	assert.NotEqual(err, errFull.Cause)
	assert.NotEqual("test error message", errFull.Error())
}
