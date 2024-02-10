package equinox

import (
	"fmt"
	"maps"
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

	if !maps.EqualFunc(p2.vals, p.vals, floatEqual) {
		t.Errorf("Expected vals '%s' got '%s'", fmt.Sprint(p.vals), fmt.Sprint(p2.vals))
	}

	if !maps.Equal(p2.attrs, p.attrs) {
		t.Errorf("Expected attrs '%s' got '%s'", fmt.Sprint(p.attrs), fmt.Sprint(p2.attrs))
	}
}
