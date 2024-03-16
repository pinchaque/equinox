package core

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
	Id    *Id
}

// Creates a new point at the specified time with a random ID. Initializes
// Vals and Attrs to be empty maps.
func NewPoint(ts time.Time) *Point {
	p := NewPointEmptyId(ts)
	p.GenerateId()
	return p
}

// Creates a new point at the specified time with an ID of zero. This is useful
// in cases where we need point bounds for a query that we can cmopare against
// real points.
func NewPointEmptyId(ts time.Time) *Point {
	p := Point{Ts: ts}
	p.Vals = make(map[string]float64)
	p.Attrs = make(map[string]string)
	p.Id = nil
	return &p
}

// Creates a new completely empty point. This is useful for when we need a blank
// point object for unmarshaling.
func NewPointEmpty() *Point {
	p := Point{}
	p.Vals = make(map[string]float64)
	p.Attrs = make(map[string]string)
	p.Id = nil
	return &p
}

// Compare function to sort points by their timestamps. Attributes and ID
// are not taken into account. Returns -1 if a < b, 0 if equal, 1 if b > a.
func PointCmp(a, b *Point) int {
	if a.Ts.UnixMicro() < b.Ts.UnixMicro() {
		return -1
	} else if a.Ts.UnixMicro() > b.Ts.UnixMicro() {
		return 1
	} else {
		return 0
	}
}

// Generates a new Id for this Point.
func (p *Point) GenerateId() {
	p.Id = NewId()
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

// Returns a clone of this point with all data members copied. Modifying the
// clone won't modify the original.
func (p *Point) Clone() *Point {
	r := NewPoint(p.Ts)
	r.Id = p.Id.Clone()

	for k, v := range p.Vals {
		r.Vals[k] = v
	}

	for k, v := range p.Attrs {
		r.Attrs[k] = v
	}

	return r
}

// Returns true if this point is identical to the other point, which means
// the timestamps, attributes, AND ID all match.
func (p *Point) Identical(oth *Point) bool {
	return p.Equal(oth) &&
		p.Id != nil &&
		oth.Id != nil &&
		p.Id.String() == oth.Id.String()
}

// Returns true if two Points are "equal", which means that the timestamp,
// values, and attributes are equal. IDs are ignored. Note that this checks
// for exact floating point equality. Use EqualTol if you want to allow for
// some error tolerance.
func (p *Point) Equal(other *Point) bool {
	return (p.Ts.UnixMicro() == other.Ts.UnixMicro() &&
		maps.Equal(p.Vals, other.Vals) &&
		maps.Equal(p.Attrs, other.Attrs))
}

// Returns true if two Points are "equal" within the specified floating point
// tolerance. This means that the timestamp, values, and attributes are equal.
// IDs are ignored.
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
