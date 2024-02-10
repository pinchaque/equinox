package equinox

import (
	"maps"
	"testing"
	"time"
)

func TestSerialize(t *testing.T) {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 3, time.UTC)
	p := NewPoint(ts)
	p.attrs["color"] = "red"
	p.attrs["shape"] = "square"
	p.vals["area"] = 43.1
	p.vals["temp"] = 21.1

	s := new(Serializer)

	data, err := s.Serialize(p)

	if err != nil {
		t.Errorf("Serialization error: %s", err.Error())
	}

	p2, err := s.Deserialize(data)

	if err != nil {
		t.Errorf("Deserialization error: %s", err.Error())
	}

	if p2.ts != p.ts {
		t.Errorf("Expected %s, got %s", p.ts.UTC(), p2.ts.UTC())
	}

	if !maps.EqualFunc(p2.vals, p.vals, floatEqual) {
		t.Errorf("Expected vals '%v' got '%v'", p.vals, p2.vals)
	}

	if !maps.Equal(p2.attrs, p.attrs) {
		t.Errorf("Expected attrs '%v' got '%v'", p.attrs, p2.attrs)
	}
}
