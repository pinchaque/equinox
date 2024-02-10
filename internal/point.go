package equinox

import "time"

type Point struct {
	ts    time.Time
	vals  map[string]float64
	attrs map[string]string
}

func NewPoint(ts time.Time) *Point {
	p := Point{ts: ts}
	p.vals = make(map[string]float64)
	p.attrs = make(map[string]string)
	return &p
}
