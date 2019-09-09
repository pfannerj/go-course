package main

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var levelDBPath = "/tmp/leveldb"

// LevelDBStore provides sync for leveldb
type LevelDBStore struct {
	sync.Mutex
	ldb    *leveldb.DB
	currID uint32
}

// NewLevelDBStore creates new storer for leveldb.
func NewLevelDBStore() *LevelDBStore {
	db, err := leveldb.OpenFile(levelDBPath, nil)
	dbErrorPanic(err)
	return &LevelDBStore{currID: 0, ldb: db}
}

// CreatePuppy creates puppy
func (l *LevelDBStore) CreatePuppy(puppy Puppy) (uint32, error) {
	l.Lock()
	defer l.Unlock()
	l.currID++
	puppy.ID = l.currID
	_, err := l.putPuppy(l.currID, puppy)
	return puppy.ID, err
}

// ReadPuppy reads puppy from backend
func (l *LevelDBStore) ReadPuppy(puppyID uint32) (Puppy, error) {
	byteID := []byte(strconv.Itoa(int(puppyID)))
	var p Puppy
	if puppy, err := l.ldb.Get(byteID, nil); err == nil {
		err := json.Unmarshal(puppy, &p)
		if err != nil {
			return p, Errorf(ErrInternalDataError, "Read failed, error retrieving corrupt data from database for puppy id %d", puppyID)
		}
		return p, nil
	}
	return p, Errorf(ErrNotFound, "Read failed, no puppy found with id %d", puppyID)
}

// UpdatePuppy updates puppy
func (l *LevelDBStore) UpdatePuppy(puppyID uint32, puppy Puppy) error {
	l.Lock()
	defer l.Unlock()
	if _, err := l.ReadPuppy(puppyID); err != nil {
		return Errorf(ErrNotFound, "Update failed, no puppy found with id %d", puppyID)
	}
	_, err := l.putPuppy(puppyID, puppy)
	return err
}

// DeletePuppy deletes puppy
func (l *LevelDBStore) DeletePuppy(puppyID uint32) error {
	l.Lock()
	defer l.Unlock()
	if _, err := l.ReadPuppy(puppyID); err != nil {
		return Errorf(ErrNotFound, "Delete failed, no puppy found with id %d", puppyID)
	}
	byteID := []byte(strconv.Itoa(int(puppyID)))
	err := l.ldb.Delete(byteID, nil)
	dbErrorPanic(err)
	return nil
}

// putPuppy stores puppy in backend
func (l *LevelDBStore) putPuppy(puppyID uint32, puppy Puppy) (uint32, error) {
	if puppy.Value < 0 {
		return puppyID, Errorf(ErrInvalidInput, "Update failed for puppy with id %d, value must not be < 0", puppyID)
	}
	puppyByte, _ := json.Marshal(puppy)
	byteID := []byte(strconv.Itoa(int(puppyID)))
	err := l.ldb.Put(byteID, puppyByte, nil)
	dbErrorPanic(err)
	return puppyID, nil
}

// dbErrorPanic causes panic in error is not nil
func dbErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
