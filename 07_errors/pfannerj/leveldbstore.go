package main

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var levelDBPath = "/tmp/leveldb"

// LevelDBStore is a leveldb based implementation of PuppyStorer.
type LevelDBStore struct {
	sync.Mutex
	ldb    *leveldb.DB
	currID uint32
}

// NewLevelDBStore creates a new leveldb store.
func NewLevelDBStore() *LevelDBStore {
	db, err := leveldb.OpenFile(levelDBPath, nil)
	checkForDBError(err)
	return &LevelDBStore{currID: 0, ldb: db}
}

// CreatePuppy adds a new puppy to the leveldb store.
func (l *LevelDBStore) CreatePuppy(puppy Puppy) (uint32, error) {
	l.Lock()
	defer l.Unlock()
	l.currID++
	puppy.ID = l.currID
	_, ok := l.writePuppy(l.currID, puppy)
	return puppy.ID, ok
}

// ReadPuppy gets a puppy from the leveldb store with the given ID.
func (l *LevelDBStore) ReadPuppy(puppyID uint32) (Puppy, error) {
	byteID := []byte(strconv.Itoa(int(puppyID)))
	var p Puppy
	if puppy, err := l.ldb.Get(byteID, nil); err == nil {
		err := json.Unmarshal(puppy, &p)
		if err != nil {
			return p, Errorf(ErrInternalDataError, "Read failed, error retrieving corrupt data from db for puppy id %d", puppyID)
		}
		return p, nil
	}
	return p, Errorf(ErrNotFound, "Read failed, no puppy found with id %d", puppyID)
}

// UpdatePuppy modifies puppy data in the leveldb store for an existing puppy.
func (l *LevelDBStore) UpdatePuppy(puppyID uint32, puppy Puppy) error {
	l.Lock()
	defer l.Unlock()
	if _, err := l.ReadPuppy(puppyID); err != nil {
		return Errorf(ErrNotFound, "Update failed, no puppy found with id %d", puppyID)
	}
	puppy.ID = puppyID //ensure the ID within pupppy always matches the leveldb store key (puppyID)
	_, err := l.writePuppy(puppyID, puppy)
	return err
}

// DeletePuppy deletes the puppy with the given ID from the leveldb store.
func (l *LevelDBStore) DeletePuppy(puppyID uint32) error {
	l.Lock()
	defer l.Unlock()
	if _, ok := l.ReadPuppy(puppyID); ok != nil {
		return Errorf(ErrNotFound, "Delete failed, no puppy found with id %d", puppyID)
	}
	byteID := []byte(strconv.Itoa(int(puppyID)))
	err := l.ldb.Delete(byteID, nil)
	checkForDBError(err)
	return nil
}

// writePuppy is used by CreatePuppy and UpdatePuppy to validate and write puppy data to the leveldb store.
func (l *LevelDBStore) writePuppy(puppyID uint32, puppy Puppy) (uint32, error) {
	if puppy.Value < 0 {
		return puppyID, Errorf(ErrInvalidInput, "Update failed for puppy with id %d, value must not be < 0", puppyID)
	}
	puppyByte, _ := json.Marshal(puppy)
	byteID := []byte(strconv.Itoa(int(puppyID)))
	err := l.ldb.Put(byteID, puppyByte, nil)
	checkForDBError(err)
	return puppyID, nil
}

// checkForDBError causes panic if error is not nil
func checkForDBError(err error) {
	if err != nil {
		panic(err)
	}
}
