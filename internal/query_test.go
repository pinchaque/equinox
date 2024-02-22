package equinox

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t2 := time.Date(2024, 01, 12, 13, 0, 0, 0, time.UTC)
	t4 := time.Date(2024, 01, 14, 13, 0, 0, 0, time.UTC)

	q := NewQuery(t2, t4)
	exp := "[2024-01-12 13:00:00 +0000 UTC-2024-01-14 13:00:00 +0000 UTC] []"
	if q.String() != exp {
		t.Errorf("expected string %s got %s", exp, q.String())
	}

	// make sure start/end inversion is accounted for
	q = NewQuery(t4, t2)
	if q.String() != exp {
		t.Errorf("expected string %s got %s", exp, q.String())
	}
}

func TestCmpTime(t *testing.T) {
	t1 := time.Date(2024, 01, 11, 13, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 01, 12, 13, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 01, 13, 13, 0, 0, 0, time.UTC)
	t4 := time.Date(2024, 01, 14, 13, 0, 0, 0, time.UTC)
	t5 := time.Date(2024, 01, 15, 13, 0, 0, 0, time.UTC)

	q := NewQuery(t2, t4)

	fn := func(ts time.Time, exp int) {
		act := q.CmpTime(NewPoint(ts))
		if act != exp {
			t.Errorf("checking %s against query %s: expected %d but got %d",
				ts.UTC(), q.String(), exp, act)
		}
	}

	fn(t1, -1)
	fn(time.Date(2024, 01, 12, 12, 59, 59, 0, time.UTC), -1)
	fn(t2, 0)
	fn(t3, 0)
	fn(t4, 0)
	fn(time.Date(2024, 01, 14, 13, 0, 0, 1, time.UTC), 0) // nanoseconds ignored
	fn(time.Date(2024, 01, 14, 13, 0, 1, 0, time.UTC), 1)
	fn(t5, 1)
}
