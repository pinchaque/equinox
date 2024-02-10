package equinox

import (
	"maps"
	"testing"
	"time"
)

func TestSerialize(t *testing.T) {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 3, time.UTC)
	const exptime string = "2024-01-10T23:01:02.000000Z"
	const fmtstr string = "2006-01-02T15:04:05.000000Z"

	if ts.Format(fmtstr) != exptime {
		t.Errorf("Time format incorrect, expected %s got %s for UTC time %s",
			exptime, ts.Format(fmtstr), ts.UTC())
	}

	p := NewPoint(ts)
	p.attrs["color"] = "red"
	p.attrs["shape"] = "square"
	p.vals["area"] = 43.1
	p.vals["temp"] = 21.1

	s := NewSerializer()

	data, err := s.Serialize(p)

	if err != nil {
		t.Errorf("Serialization error: %s", err.Error())
	}

	// expected size: 16 + 12*num_values + 8*num_attrs = 16 + 24 + 16 = 56
	expsize := 56
	if len(data) != expsize {
		t.Errorf("Serialization returned %d bytes, expected %d", len(data), expsize)
	}

	p2, err := s.Deserialize(data)

	if err != nil {
		t.Errorf("Deserialization error: %s", err.Error())
	}

	if p2.ts.UnixMicro() != p.ts.UnixMicro() {
		t.Errorf("Expected %d, got %d", p.ts.UnixMicro(), p2.ts.UnixMicro())
	}

	if p2.ts.Format(fmtstr) != exptime {
		t.Errorf("Expected %s got %s", exptime, p2.ts.Format(fmtstr))
	}

	if !maps.EqualFunc(p2.vals, p.vals, floatEqual) {
		t.Errorf("Expected vals '%v' got '%v'", p.vals, p2.vals)
	}

	if !maps.Equal(p2.attrs, p.attrs) {
		t.Errorf("Expected attrs '%v' got '%v'", p.attrs, p2.attrs)
	}
}
