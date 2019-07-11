package main

import (
	"bytes"
	"testing"
)

func TestMain(t *testing.T) {
	var buf bytes.Buffer
	mainout = &buf
	main()
	expected := `Map Puppy Created with ID: 1
Map Puppy read: &{1 Labrador Brown 999.99}
Map Puppy updated: 1
Map Puppy deleted: true
Sync Puppy Created with ID: 1
Sync Puppy read: &{1 Labrador Brown 999.99}
Sync Puppy updated: 1
Sync Puppy deleted: true
`
	actual := buf.String()
	if expected != actual {
		t.Errorf("Unexpected output in main(), Expected: %v, Actual: %v", expected, actual)
	}
}