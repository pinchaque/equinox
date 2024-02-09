package equinox

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {

	ts := time.Date(2024, 01, 10, 23, 1, 2, 3, time.UTC)
	p := NewPoint(ts)

	if ts != p.ts {
		t.Errorf("Got %s, wanted %s", p.ts.UTC(), ts.UTC())
	}

	if len(p.vals) != 0 {
		t.Errorf("Expected 0 values, got %d", len(p.vals))
	}

	if len(p.attrs) != 0 {
		t.Errorf("Expected 0 attributes, got %d", len(p.attrs))
	}
}
