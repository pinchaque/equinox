package query

import (
	"equinox/internal/core"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t2 := time.Date(2024, 01, 12, 13, 0, 0, 0, time.UTC)
	t4 := time.Date(2024, 01, 14, 13, 0, 0, 0, time.UTC)

	q := NewQuery(t2, t4, True())
	exp := "[2024-01-12 13:00:00 +0000 UTC-2024-01-14 13:00:00 +0000 UTC] [true]"

	assert.Equal(t, exp, q.String())

	// make sure start/end inversion is accounted for
	q = NewQuery(t4, t2, True())
	assert.Equal(t, exp, q.String())
}

func TestCmpTime(t *testing.T) {
	t1 := time.Date(2024, 01, 11, 13, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 01, 12, 13, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 01, 13, 13, 0, 0, 0, time.UTC)
	t4 := time.Date(2024, 01, 14, 13, 0, 0, 0, time.UTC)
	t5 := time.Date(2024, 01, 15, 13, 0, 0, 0, time.UTC)

	q := NewQuery(t2, t4, True())

	fn := func(ts time.Time, exp int) {
		act := q.CmpTime(core.NewPoint(ts))
		assert.Equal(t, exp, act)
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
func TestQueryAttrs(t *testing.T) {
	ts1 := time.Date(2024, 01, 11, 13, 0, 0, 0, time.UTC)
	ts2 := time.Date(2024, 01, 12, 13, 0, 0, 0, time.UTC)
	ts3 := time.Date(2024, 01, 13, 13, 0, 0, 0, time.UTC)

	p := core.NewPoint(ts1)
	p.Attrs = make(map[string]string)
	p.Attrs["color"] = "blue"
	p.Attrs["animal"] = "moose"
	p.Attrs["shape"] = "square"
	p.Attrs["index"] = "74"

	fnsub := func(q *Query, exp_ts bool, exp_attr bool) {
		// check attr only match
		assert.Equal(t, exp_attr, q.MatchAttr(p))

		// check time only match
		assert.Equal(t, exp_ts, q.MatchTime(p))

		// check timestamp+attr match
		assert.Equal(t, exp_ts && exp_attr, q.Match(p))
	}

	fn := func(qa FilterAttr, exp_attr bool) {
		// check a few different time ranges
		fnsub(NewQuery(ts1, ts3, qa), true, exp_attr)
		fnsub(NewQuery(ts1, ts1, qa), true, exp_attr)
		fnsub(NewQuery(ts2, ts3, qa), false, exp_attr)
		fnsub(NewQuery(ts3, ts3, qa), false, exp_attr)
	}

	// things that are true
	t1 := Equal("color", "blue")
	t2 := Equal("animal", "moose")
	t3 := Regex("index", `^\d+$`)
	t4 := Equal("shape", "square")

	// things that are false
	f1 := Regex("color", "^x")
	f2 := Regex("animal", "mo{3,5}se")
	f3 := Equal("index", "777")
	f4 := Regex("flavor", "sour")

	fn(True(), true)
	fn(Exists("color"), true)
	fn(Exists("flavor"), false)
	fn(t1, true)
	fn(t2, true)
	fn(t3, true)
	fn(t4, true)
	fn(f1, false)
	fn(f2, false)
	fn(f3, false)
	fn(f4, false)
	fn(Or(Exists("flavor"), f1, f2, f3, f4), false)
	fn(Or(Exists("flavor"), t1, f2, f3, f4), true)
	fn(Or(And(f1, f2, f3, f4), Not(t1), Or(Not(f4), t2, t3)), true)
	fn(Or(And(f1, f2, f3, f4), Not(t4), Or(Not(t4), f3)), false)

}
