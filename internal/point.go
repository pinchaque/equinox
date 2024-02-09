package equinox

import "time"

type Point struct {
	ts    time.Time
	vals  map[string]string
	attrs map[string]string
}

func NewPoint(ts time.Time) *Point {
	p := Point{ts: ts}
	return &p
}
