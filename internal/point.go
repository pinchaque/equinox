package equinox

import (
	"fmt"
	"maps"
	"math"
	"sort"
	"strings"
	"time"
)

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

// function to sort points
// returns -1 if a < b, 0 if equal, 1 if b > a
func PointCmp(a, b *Point) int {
	if a.ts.UnixMicro() < b.ts.UnixMicro() {
		return -1
	} else if a.ts.UnixMicro() > b.ts.UnixMicro() {
		return 1
	} else {
		return 0
	}
}

func (p *Point) String() string {
	var val, attr []string

	for k, v := range p.vals {
		val = append(val, k+": "+fmt.Sprintf("%f", v))
	}
	sort.Strings(val) // ensure consistent output

	for k, v := range p.attrs {
		attr = append(attr, k+": "+v)
	}
	sort.Strings(attr) // ensure consistent output

	return fmt.Sprintf("[%s] val[%s] attr[%s]",
		p.ts.UTC(),
		strings.Join(val, ", "),
		strings.Join(attr, ", "))
}

// equal (including exact floating point equality)
func (p *Point) Equal(other *Point) bool {
	return (p.ts.UnixMicro() == other.ts.UnixMicro() &&
		maps.Equal(p.vals, other.vals) &&
		maps.Equal(p.attrs, other.attrs))
}

// equal within a given floating point tolerance
func (p *Point) EqualTol(other *Point, tol float64) bool {
	cmp := func(x, y float64) bool {
		diff := math.Abs(x - y)
		mean := math.Abs(x+y) / 2.0
		if math.IsNaN(diff / mean) {
			return true
		}
		return (diff / mean) < tol
	}

	return (p.ts.UnixMicro() == other.ts.UnixMicro() &&
		maps.EqualFunc(p.vals, other.vals, cmp) &&
		maps.Equal(p.attrs, other.attrs))
}
