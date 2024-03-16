package models

import "equinox/internal/engine"

// Structure representing a data series.
type Series struct {
	Id string
	IO engine.PointIO
}
