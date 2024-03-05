package core

import (
	"testing"
	"time"
)

func TestPointCreate(t *testing.T) {

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

	if p.Id.val <= 0 {
		t.Errorf("Expected an id >= 0, got %d", p.Id.val)
	}

	if p.Id.String() == "" {
		t.Errorf("Expected an id, got empty string")
	}
}
func TestPointCreateEmptyId(t *testing.T) {

	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := NewPointEmptyId(ts)

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

	if p.Id.val != 0 {
		t.Errorf("Expected an id of 0, got %d", p.Id.val)
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

	cmp := func(pt1 *Point, pt2 *Point, exp int) {
		act := PointCmp(pt1, pt2)
		if act != exp {
			t.Errorf("PointCmp(%s, %s) expected %d got %d", pt1.String(), pt2.String(), exp, act)
		}
	}

	{ // basic equality
		p2 := newPointComplete()
		if !p1.Equal(p2) {
			t.Errorf("Expected equal, got inequal: %s compared to %s", p1.String(), p2.String())
		}

		cmp(p1, p2, 0)
	}

	{ // different timestamp
		p2 := newPointComplete()
		p2.Ts = p2.Ts.AddDate(0, 0, 1)
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}

		cmp(p1, p2, -1)
		cmp(p2, p1, 1)
	}

	{ // changed value
		p2 := newPointComplete()
		cmp(p1, p2, 0)
		p2.Vals["area"] = 43.1004
		cmp(p1, p2, 0) // only timestamp matters
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

		// test 0.00 value
		p2.Vals["area"] = 0.000
		if !p2.EqualTol(p2, 0.001) {
			t.Errorf("Expected equal, got inequal: %s compared to %s", p2.String(), p2.String())
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
		cmp(p1, p2, 0)
		p2.Attrs["color"] = "blue"
		cmp(p1, p2, 0) // only timestamp matters
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // add attr
		p2 := newPointComplete()
		cmp(p1, p2, 0)
		p2.Attrs["color2"] = "blue"
		cmp(p1, p2, 0) // only timestamp matters
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

	fn := func(pt1 *Point, pt2 *Point, exp bool) {
		act := pt1.Identical(pt2)
		if act != exp {
			t.Errorf("Expected identity %t got %t: [%s]%s compared to [%s]%s",
				exp, act, pt1.Id.String(), pt1.String(), pt2.Id.String(), pt2.String())
		}
	}

	fn(p1, p2, false)
	fn(p2, p1, false)
	fn(p1, p1, true)
	fn(p2, p2, true)

	// make them have the same IDs so now should be identical
	p1.Id = p2.Id
	fn(p1, p2, true)
	fn(p2, p1, true)
	fn(p1, p1, true)
	fn(p2, p2, true)

	// now change an attribute => no longer identical
	p1.Attrs["color"] = "mauve"
	fn(p1, p2, false)
	fn(p2, p1, false)
	fn(p1, p1, true)
	fn(p2, p2, true)
}

func TestPointClone(t *testing.T) {
	/*
		ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
		p := NewPoint(ts)
		p.Attrs["shape"] = "square"
		p.Attrs["color"] = "red"
		p.Vals["area"] = 43.1
		p.Vals["temp"] = 21.1
	*/

	p1 := newPointComplete()
	p2 := p1.Clone()

	// these should be identical and equal
	if !p1.Equal(p2) {
		t.Errorf("Expected equal but got inequal: [%s]%s compared to [%s]%s",
			p1.Id.String(), p1.String(), p2.Id.String(), p2.String())
	}

	if !p1.Identical(p2) {
		t.Errorf("Expected identical but got different: [%s]%s compared to [%s]%s",
			p1.Id.String(), p1.String(), p2.Id.String(), p2.String())
	}

	// if we change one it shouldn't affect the other
	p1.Attrs["shape"] = "triangle"
	if p2.Attrs["shape"] != "square" {
		t.Errorf("changing p1.Attrs[shape] to triangle also changed p2.Attrs[shape] to %s", p2.Attrs["shape"])
	}

	p1.Vals["area"] = 1.0
	if p2.Vals["area"] != 43.1 {
		t.Errorf("changing p1.Vals[area] to triangle also changed p2.Vals[area] to %f", p2.Vals["area"])
	}

	p1.Ts = time.Date(2024, 01, 11, 23, 1, 2, 0, time.UTC)
	if p2.Ts.String() != "2024-01-10 23:01:02 +0000 UTC" {
		t.Errorf("changing p1.Ts to %s also changed p2.Ts to %s", p1.Ts.String(), p2.Ts.String())
	}

	p1.Id.val = 3
	if p2.Id.val == 3 {
		t.Errorf("changing p1.Id to 3 also changed p2.Ts to 3")
	}
}
