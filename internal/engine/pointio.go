package engine

import (
	"equinox/internal/core"
	"equinox/internal/query"
)

type PointIO interface {
	Add(p ...*core.Point) error
	Len() int
	Flush() error
	Vacuum() error
	Search(q *query.Query) (*query.QueryExec, error)
	Name() string
	String() string
	First() *core.Point
	Last() *core.Point
}
