package main

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const firstPuppyID uint32 = 1

var (
	firstPuppy = func() Puppy {
		return Puppy{
			Breed:  "Retriever",
			Colour: "Golden",
			Value:  9999.99,
		}
	}
	modifiedPuppy = func() Puppy {
		return Puppy{
			Breed:  "Bitsa",
			Colour: "Mixed",
			Value:  1.99,
		}
	}
	anotherPuppy = func() Puppy {
		return Puppy{
			Breed:  "Labrador",
			Colour: "Black",
			Value:  0,
		}
	}
	invalidPuppy = func() Puppy {
		return Puppy{
			Breed:  "Poodle",
			Colour: "White",
			Value:  -23.67,
		}
	}
)

type storerSuite struct {
	suite.Suite
	newStore func() Storer
	store    Storer
}

func (s *storerSuite) SetupTest() {
	// create test store and add the first puppy
	s.store = s.newStore()
	puppy := firstPuppy()
	_, err := s.store.CreatePuppy(puppy)
	if err != nil {
		panic("Failed to setup puppy test")
	}
}

func TestStorer(t *testing.T) {
	suite.Run(t, &storerSuite{
		newStore: func() Storer { return &SyncStore{} },
	})
	suite.Run(t, &storerSuite{
		newStore: func() Storer { return &MapStore{puppyMap: PuppyMap{}} },
	})
	//suite.Run(t, &storerSuite{newStore: NewMapStore})
}

func (s *storerSuite) TestCreate() {
	// given
	assert := tassert.New(s.T())
	newPuppy := firstPuppy()

	// when
	createdPuppyID, err := s.store.CreatePuppy(newPuppy)
	newPuppy.ID = createdPuppyID

	// then
	assert.NoError(err, "Should not get an error creating a puppy")
	foundPuppy, err := s.store.ReadPuppy(createdPuppyID)
	assert.NoError(err, "Should be able to read an newly created puppy")
	assert.Equal(newPuppy, foundPuppy, "Created puppy should be identical to the one passed to create")
}

func (s *storerSuite) TestCreateZero() {
	// given
	assert := tassert.New(s.T())
	newPuppy := anotherPuppy()

	// when
	createdPuppyID, err := s.store.CreatePuppy(newPuppy)
	newPuppy.ID = createdPuppyID

	// then
	assert.NoError(err, "Should not get an error creating a puppy with value = 0")
	foundPuppy, err := s.store.ReadPuppy(createdPuppyID)
	assert.NoError(err, "Should be able to read an newly created puppy")
	assert.Equal(newPuppy, foundPuppy, "Created puppy should be identical to the one passed to create")
}

func (s *storerSuite) TestCreateFailInvalidInput() {
	// given
	assert := tassert.New(s.T())
	newPuppy := invalidPuppy()

	// when
	_, err := s.store.CreatePuppy(newPuppy)

	// then
	assert.Error(err, "Should get an error creating a puppy with value < 0")
}

func (s *storerSuite) TestRead() {
	// given
	assert := tassert.New(s.T())
	expectedPuppy := firstPuppy()
	expectedPuppy.ID = firstPuppyID

	// when
	foundPuppy, err := s.store.ReadPuppy(firstPuppyID)

	// then
	assert.NoError(err, "Should not get an error reading puppy from store")
	assert.Equal(expectedPuppy, foundPuppy, "Store should return a puppy identical to firstPuppy")
}

func (s *storerSuite) TestReadFail() {
	// given
	assert := tassert.New(s.T())

	// when
	_, err := s.store.ReadPuppy(99)

	// then
	assert.Error(err, "Should get an error when attempting to read a non-existent puppy")
}

func (s *storerSuite) TestUpdate() {
	// given
	assert := tassert.New(s.T())
	updatePuppy := modifiedPuppy()
	updatePuppy.ID = firstPuppyID

	// when
	err := s.store.UpdatePuppy(firstPuppyID, updatePuppy)

	// then
	assert.NoError(err, "Should not get an error updating a puppy")
	foundPuppy, err := s.store.ReadPuppy(updatePuppy.ID)
	assert.NoError(err, "Should not get an error reading the updated puppy")
	assert.Equal(updatePuppy, foundPuppy, "Found puppy should be equal to updated puppy")
}

func (s *storerSuite) TestUpdateZero() {
	// given
	assert := tassert.New(s.T())
	updatePuppy := anotherPuppy()

	// when
	err := s.store.UpdatePuppy(1, updatePuppy)

	// then
	assert.NoError(err, "Should not get an error updating a puppy with value = 0")
}

func (s *storerSuite) TestUpdateFailNotFound() {
	// given
	assert := tassert.New(s.T())
	updatePuppy := anotherPuppy()

	// when
	err := s.store.UpdatePuppy(99, updatePuppy)

	// then
	assert.Error(err, "Should get an error updating a puppy")
}

func (s *storerSuite) TestUpdateFailInvalidInput() {
	// given
	assert := tassert.New(s.T())
	updatePuppy := invalidPuppy()

	// when
	err := s.store.UpdatePuppy(1, updatePuppy)

	// then
	assert.Error(err, "Should get an error updating a puppy with value < 0")
}

func (s *storerSuite) TestDeleteExisting() {
	// when
	assert := tassert.New(s.T())
	err := s.store.DeletePuppy(firstPuppyID)

	// then
	assert.NoError(err, "Should not get an error deleting an existing puppy")
	_, err = s.store.ReadPuppy(firstPuppyID)
	assert.Error(err, "Should not be able to read a deleted puppy")
}

func (s *storerSuite) TestDeleteNotExisting() {
	// when
	assert := tassert.New(s.T())
	err := s.store.DeletePuppy(99)

	// then
	assert.Error(err, "Should get an error deleting a non-existent puppy")
}
