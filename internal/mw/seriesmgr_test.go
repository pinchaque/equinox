package mw

import (
	"equinox/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeriesMgrSingleton(t *testing.T) {

	// whitebox testing - instance hasn't been created yet
	assert.Nil(t, seriesMgrInst)

	mgr := GetSeriesMgr()
	assert.NotNil(t, mgr)

	// whitebox testing - instance should exist
	assert.NotNil(t, seriesMgrInst)

	mgr2 := GetSeriesMgr()
	assert.Equal(t, mgr, mgr2) // should be pointer to same object
}

func TestSeriesMgrAdding(t *testing.T) {
	mgr := GetSeriesMgr()
	assert.Equal(t, 0, mgr.Size())

	s := &models.Series{Id: "foobar"}

	// shouldn't exist yet
	assert.False(t, mgr.Has(s.Id))

	// this should be an error
	r, err := mgr.Get(s.Id)
	assert.Nil(t, r)
	assert.Error(t, err)
	assert.Equal(t, `series 'foobar' does not exist`, err.Error())

	// now add it
	err = mgr.Add(s)
	assert.NoError(t, err)
	assert.Equal(t, 1, mgr.Size())

	// now get it back out
	r, err = mgr.Get(s.Id)
	assert.NoError(t, err)
	assert.Equal(t, s.Id, r.Id)
}
