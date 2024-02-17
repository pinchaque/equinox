package equinox

import (
	"time"
)

type Cursor interface {
	Fetch(n int) ([]*Point, error)
}

type Query struct {
	start time.Time
	end   time.Time
	attrs map[string]string
}
