package mw

import (
	"equinox/internal/models"
	"fmt"
)

// Manages access to underlying data series objects, providing caching and
// lookup
type seriesMgr struct {
	series map[string]*models.Series
}

// Singleton instance of seriesMgr
var seriesMgrInst *seriesMgr

// Returns singleton instance of the data series manager.
func GetSeriesMgr() *seriesMgr {
	if seriesMgrInst == nil {
		seriesMgrInst = &seriesMgr{series: make(map[string]*models.Series)}
	}

	return seriesMgrInst
}

// Returns number of elements currently in the series manager
func (sm *seriesMgr) Size() int {
	return len(seriesMgrInst.series)
}

// Retrieves the data series with the given ID, returning an error if it does
// not exist.
func (sm *seriesMgr) Get(id string) (*models.Series, error) {
	s, exist := seriesMgrInst.series[id]
	if !exist {
		return nil, fmt.Errorf("series '%s' does not exist", id)
	}

	return s, nil
}

// Returns true if the data series with given id already exists, false othersie
func (sm *seriesMgr) Has(id string) bool {
	_, exist := seriesMgrInst.series[id]
	return exist
}

// Adds the given series to the manager if one with that id doesn't already
// exist. If it exists then nothing is added an an error is returned.
func (sm *seriesMgr) Add(s *models.Series) error {
	_, exist := seriesMgrInst.series[s.Id]
	if exist {
		return fmt.Errorf("series '%s' already exists", s.Id)
	}

	seriesMgrInst.series[s.Id] = s

	return nil
}
