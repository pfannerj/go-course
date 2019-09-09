package main

import (
	"fmt"
	"sync"
)

// SyncStore is a sync.Map based in-memory implementation of PuppyStorer.
type SyncStore struct {
	sync.Mutex
	sync.Map
	currID uint32
}

// NewSyncStore creates a new in-memory store with sync map intialised.
func NewSyncStore() *SyncStore {
	return &SyncStore{}
}

// CreatePuppy adds a new puppy to the sync store.
func (s *SyncStore) CreatePuppy(puppy Puppy) (uint32, error) {
	s.Lock()
	defer s.Unlock()
	s.currID++
	puppy.ID = s.currID
	if puppy.Value < 0 {
		return puppy.ID, Errorf(ErrInvalidInput, "Create failed for puppy with id %d, value must not be < 0", puppy.ID)
	}
	s.Store(puppy.ID, puppy)
	return puppy.ID, nil
}

// ReadPuppy gets a puppy from the sync store with the given ID.
func (s *SyncStore) ReadPuppy(puppyID uint32) (Puppy, error) {
	if puppy, ok := s.Load(puppyID); ok {
		puppyOut, _ := puppy.(Puppy)
		return puppyOut, nil
	}
	return Puppy{}, fmt.Errorf("no puppy found with id %d", puppyID)
}

// UpdatePuppy modifies puppy data in the sync store for an existing puppy.
func (s *SyncStore) UpdatePuppy(puppyID uint32, puppy Puppy) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Load(puppyID); !ok {
		return Errorf(ErrNotFound, "Update failed, no puppy found with id %d", puppyID)
	}
	if puppy.Value < 0 {
		return Errorf(ErrInvalidInput, "Update failed for puppy with id %d, value must not be < 0", puppyID)
	}
	puppy.ID = puppyID //ensure the ID within p always matches the sync store key (puppyID)
	s.Store(puppyID, puppy)
	return nil
}

// DeletePuppy deletes the puppy with the given ID from the sync store.
func (s *SyncStore) DeletePuppy(puppyID uint32) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Load(puppyID); ok {
		s.Delete(puppyID)
		return nil
	}
	return Errorf(ErrNotFound, "Delete failed, no puppy found with id %d", puppyID)
}
