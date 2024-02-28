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
	Ts    time.Time
	Vals  map[string]float64
	Attrs map[string]string
}

func NewPoint(ts time.Time) *Point {
	p := Point{Ts: ts}
	p.Vals = make(map[string]float64)
	p.Attrs = make(map[string]string)
	return &p
}

// Compare function to sort points by time only. Returns -1 if a < b, 0 if equal, 1 if b > a.
func PointCmpTime(a, b *Point) int {
	if a.Ts.UnixMicro() < b.Ts.UnixMicro() {
		return -1
	} else if a.Ts.UnixMicro() > b.Ts.UnixMicro() {
		return 1
	} else {
		return 0
	}
}

// Compare function to sort points by all fields including attributes. Returns -1 if a < b, 0 if equal, 1 if b > a.
func PointCmp(a, b *Point) int {
	return PointCmpTime(a, b)
}

func (p *Point) String() string {
	var val, attr []string

	for k, v := range p.Vals {
		val = append(val, k+": "+fmt.Sprintf("%f", v))
	}
	sort.Strings(val) // ensure consistent output

	for k, v := range p.Attrs {
		attr = append(attr, k+": "+v)
	}
	sort.Strings(attr) // ensure consistent output

	return fmt.Sprintf("[%s] val[%s] attr[%s]",
		p.Ts.UTC(),
		strings.Join(val, ", "),
		strings.Join(attr, ", "))
}

// Returns true if this point is less than the other
func (p *Point) Less(oth *Point) bool {
	return PointCmpTime(p, oth) < 0
}

// equal (including exact floating point equality)
func (p *Point) Equal(other *Point) bool {
	return (p.Ts.UnixMicro() == other.Ts.UnixMicro() &&
		maps.Equal(p.Vals, other.Vals) &&
		maps.Equal(p.Attrs, other.Attrs))
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

	return (p.Ts.UnixMicro() == other.Ts.UnixMicro() &&
		maps.EqualFunc(p.Vals, other.Vals, cmp) &&
		maps.Equal(p.Attrs, other.Attrs))
}
