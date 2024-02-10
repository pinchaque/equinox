package equinox

import (
	"math"
	"testing"
	"time"
)

func floatEqual(x, y float64) bool {
	const tolerance = 0.00001
	diff := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	if math.IsNaN(diff / mean) {
		return true
	}
	return (diff / mean) < tolerance
}

func TestCreate(t *testing.T) {

	ts := time.Date(2024, 01, 10, 23, 1, 2, 3, time.UTC)
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

func TestSerialize(t *testing.T) {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 3, time.UTC)
	p := NewPoint(ts)
	p.attrs["color"] = "red"
	p.attrs["shape"] = "square"
	p.vals["area"] = 43.1
	p.vals["temp"] = 21.1

	data, err := p.Serialize()

	if err != nil {
		t.Errorf("Serialization error: %s", err.Error())
	}

	p2 := Deserialize(data)

	if p2.ts != p.ts {
		t.Errorf("Expected %s, got %s", p.ts.UTC(), p2.ts.UTC())
	}

	if len(p2.vals) != len(p.vals) {
		t.Errorf("Expected %d values, got %d", len(p.vals), len(p2.vals))
	}

	for k, v := range p.vals {
		if !floatEqual(v, p2.vals[k]) {
			t.Errorf("Expected vals[%s] to be %f, got %f", k, v, p2.vals[k])
		}
	}

	if len(p2.attrs) != len(p.attrs) {
		t.Errorf("Expected %d values, got %d", len(p.attrs), len(p2.attrs))
	}

	for k, v := range p.attrs {
		if p2.attrs[k] != v {
			t.Errorf("Expected attrs[%s] to be %s, got %s", k, v, p2.attrs[k])
		}
	}
}
