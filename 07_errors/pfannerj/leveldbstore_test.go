package main

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	tassert "github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestInternalDataError(t *testing.T) {
	// Setup corrupt data to cause internal error
	func() {
		assert := tassert.New(t)
		db, _ := leveldb.OpenFile(levelDBPath, nil)
		defer db.Close()
		puppy, err := json.Marshal("Corrupt puppy data")
		assert.NoError(err, "Should not get an error marshalling corrupt puppy data")
		puppyID := []byte(strconv.Itoa(999))
		err = db.Put(puppyID, puppy, nil)
		assert.NoError(err, "Should not get an error writing corrupt data to db")
	}()
	l := NewLevelDBStore()
	defer l.ldb.Close()
	_, err := l.ReadPuppy(999)
	assert.Error(t, err, "Should get an error reading the corrupt puppy")
}

func TestCheckForDBErrorPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("checkForDBError did panic")
		}
	}()
	err := Errorf(ErrInternalDataError, "test data error")
	checkForDBError(err)
}

func TestCheckForDBErrorDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("checkForDBError did not panic")
		}
	}()
	checkForDBError(nil)
}
