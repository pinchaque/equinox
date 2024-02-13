package equinox

import (
	"testing"
	"time"
)

func TestPointCreateEmpty(t *testing.T) {

	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := NewPoint(ts)

	if ts != p.ts {
		t.Errorf("Got %s, wanted %s", p.ts.UTC(), ts.UTC())
	}

	if p.vals == nil {
		t.Errorf("Values is nil")
	}

	if len(p.vals) != 0 {
		t.Errorf("Expected 0 values, got %d", len(p.vals))
	}

	if p.attrs == nil {
		t.Errorf("Attrs is nil")
	}

	if len(p.attrs) != 0 {
		t.Errorf("Expected 0 attributes, got %d", len(p.attrs))
	}
}

func newPointComplete() *Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := NewPoint(ts)
	p.attrs["shape"] = "square"
	p.attrs["color"] = "red"
	p.vals["area"] = 43.1
	p.vals["temp"] = 21.1
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

	if ts != p.ts {
		t.Errorf("Got %s, wanted %s", p.ts.UTC(), ts.UTC())
	}

	if len(p.vals) != 2 {
		t.Errorf("Expected 2 values, got %d", len(p.vals))
	}

	if len(p.attrs) != 2 {
		t.Errorf("Expected 2 attributes, got %d", len(p.attrs))
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
		p2.ts = p2.ts.AddDate(0, 0, 1)
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // changed value
		p2 := newPointComplete()
		p2.vals["area"] = 43.1004
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
		p2.vals["area2"] = 49.999
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // delete value
		p2 := newPointComplete()
		delete(p2.vals, "area")
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // changed attr
		p2 := newPointComplete()
		p2.attrs["color"] = "blue"
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // add attr
		p2 := newPointComplete()
		p2.attrs["color2"] = "blue"
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}

	{ // delete attr
		p2 := newPointComplete()
		delete(p2.attrs, "color")
		if p1.Equal(p2) {
			t.Errorf("Expected inequal, got equal: %s compared to %s", p1.String(), p2.String())
		}
	}
}
