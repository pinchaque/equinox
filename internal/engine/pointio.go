package engine

import (
	"equinox/internal/core"
	"equinox/internal/query"
	"time"
)

type PointIO interface {
	// Adds the specified points to the storage engine
	Add(p ...*core.Point) error

	// Returns the number of points stored in the engine
	Len() int

	// Flushes points to persistent storage, if applicable
	Flush() error

	// Performs cleanup tasks on the storage engine
	Vacuum() error

	// Searches the storage engine for points matching the specified query
	Search(q *query.Query) (*query.QueryExec, error)

	// Returns a short name of the storage engine
	Name() string

	// Returns a detailed string including all data stored in the engine,
	// mostly useful for debugging
	String() string

	// Returns the first chronological point in the storage engine
	First() *core.Point

	// Returns the last chronological point in the storage engine
	Last() *core.Point

	// Moves all points in [start, end] to another storage engine, returning
	// the number of points moved.
	Move(start time.Time, end time.Time, dest PointIO) (int, error)

	// Extracts all contained points, useful for debugging and testing
	extract() ([]*core.Point, error)
}
