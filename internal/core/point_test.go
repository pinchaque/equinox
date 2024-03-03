package core

import (
	"testing"
	"time"
)

func TestPointCreateEmpty(t *testing.T) {

	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := NewPoint(ts)

	if ts != p.Ts {
		t.Errorf("Got %s, wanted %s", p.Ts.UTC(), ts.UTC())
	}

	if p.Vals == nil {
		t.Errorf("Values is nil")
	}

	if len(p.Vals) != 0 {
		t.Errorf("Expected 0 values, got %d", len(p.Vals))
	}

	if p.Attrs == nil {
		t.Errorf("Attrs is nil")
	}

	if len(p.Attrs) != 0 {
		t.Errorf("Expected 0 attributes, got %d", len(p.Attrs))
	}

	if p.Uuid.String() == "" {
		t.Errorf("Expected a guid, got empty string")
	}
}

func newPointComplete() *Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := NewPoint(ts)
	p.Attrs["shape"] = "square"
	p.Attrs["color"] = "red"
	p.Vals["area"] = 43.1
	p.Vals["temp"] = 21.1
	return p
}

func TestPointString(t *testing.T) {
	p := newPointComplete()
	exp := "[2024-01-10 23:01:02 +0000 UTC] val[area: 43.100000, temp: 21.100000] attr[color: red, shape: square]"
	if p.String() != exp {
		t.Errorf("Expected %s, got %s", exp, p.String())
	}
}

func TestPointCreateComplete(t *testing.T) {
	p := newPointComplete()
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)

	if ts != p.Ts {
		t.Errorf("Got %s, wanted %s", p.Ts.UTC(), ts.UTC())
	}

	if len(p.Vals) != 2 {
		t.Errorf("Expected 2 values, got %d", len(p.Vals))
	}

	if len(p.Attrs) != 2 {
		t.Errorf("Expected 2 attributes, got %d", len(p.Attrs))
	}
}

func TestPointEqual(t *testing.T) {
	p1 := newPointComplete()

	{ // basic equality
		p2 := newPointComplete()
		if !p1.Equal(p2) {
			t.Errorf("Expected equal, got inequal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // different timestamp
		p2 := newPointComplete()
		p2.Ts = p2.Ts.AddDate(0, 0, 1)
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // changed value
		p2 := newPointComplete()
		p2.Vals["area"] = 43.1004
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}

		if !p1.EqualTol(p2, 0.1) {
			t.Errorf("Expected equal, got inequal: %s compared to %s", p1.String(), p2.String())
		}

		if !p1.EqualTol(p2, 0.01) {
			t.Errorf("Expected equal, got inequal: %s compared to %s", p1.String(), p2.String())
		}

		if !p1.EqualTol(p2, 0.001) {
			t.Errorf("Expected equal, got inequal: %s compared to %s", p1.String(), p2.String())
		}

		if !p1.EqualTol(p2, 0.0001) {
			t.Errorf("Expected equal, got inequal: %s compared to %s", p1.String(), p2.String())
		}

		if !p1.EqualTol(p2, 0.00001) {
			t.Errorf("Expected equal, got inequal: %s compared to %s", p1.String(), p2.String())
		}

		if p1.EqualTol(p2, 0.000001) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}

		if p1.EqualTol(p2, 0.0000001) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // add value
		p2 := newPointComplete()
		p2.Vals["area2"] = 49.999
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // delete value
		p2 := newPointComplete()
		delete(p2.Vals, "area")
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // changed attr
		p2 := newPointComplete()
		p2.Attrs["color"] = "blue"
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // add attr
		p2 := newPointComplete()
		p2.Attrs["color2"] = "blue"
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // delete attr
		p2 := newPointComplete()
		delete(p2.Attrs, "color")
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}
}

func TestPointIdentical(t *testing.T) {
	/*
		ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
		p := NewPoint(ts)
		p.Attrs["shape"] = "square"
		p.Attrs["color"] = "red"
		p.Vals["area"] = 43.1
		p.Vals["temp"] = 21.1
	*/

	p1 := newPointComplete()
	p2 := newPointComplete()

	fn := func(pt1 *Point, pt2 *Point, act bool, exp bool) {
		if act != exp {
			t.Errorf("Expected identity %t got %t: [%s]%s compared to [%s]%s",
				exp, act, pt1.Uuid.String(), pt1.String(), pt2.Uuid.String(), pt2.String())
		}
	}

	fn(p1, p2, p1.Identical(p2), false)
	fn(p2, p1, p2.Identical(p1), false)
	fn(p1, p1, p1.Identical(p1), true)
	fn(p2, p2, p2.Identical(p2), true)

	// make them have the same UUIDs so now should be identical
	p1.Uuid = p2.Uuid
	fn(p1, p2, p1.Identical(p2), true)
	fn(p2, p1, p2.Identical(p1), true)
	fn(p1, p1, p1.Identical(p1), true)
	fn(p2, p2, p2.Identical(p2), true)

	// now change an attribute => no longer identical
	p1.Attrs["color"] = "mauve"
	fn(p1, p2, p1.Identical(p2), false)
	fn(p2, p1, p2.Identical(p1), false)
	fn(p1, p1, p1.Identical(p1), true)
	fn(p2, p2, p2.Identical(p2), true)
}
