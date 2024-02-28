package equinox

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t2 := time.Date(2024, 01, 12, 13, 0, 0, 0, time.UTC)
	t4 := time.Date(2024, 01, 14, 13, 0, 0, 0, time.UTC)

	q := NewQuery(t2, t4, NewQATrue())
	exp := "[2024-01-12 13:00:00 +0000 UTC-2024-01-14 13:00:00 +0000 UTC] [true]"
	if q.String() != exp {
		t.Errorf("expected string %s got %s", exp, q.String())
	}

	// make sure start/end inversion is accounted for
	q = NewQuery(t4, t2, NewQATrue())
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

	q := NewQuery(t2, t4, NewQATrue())

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
func TestQueryAttrs(t *testing.T) {
	ts1 := time.Date(2024, 01, 11, 13, 0, 0, 0, time.UTC)
	ts2 := time.Date(2024, 01, 12, 13, 0, 0, 0, time.UTC)
	ts3 := time.Date(2024, 01, 13, 13, 0, 0, 0, time.UTC)

	/*
		r["color"] = "blue"
		r["animal"] = "moose"
		r["shape"] = "square"
		r["index"] = "74"
	*/
	p := NewPoint(ts1)
	p.Attrs = testGetAttrs()

	fnsub := func(q *Query, exp_ts bool, exp_attr bool) {
		var act bool

		// check attr only match
		act = q.MatchAttr(p)
		if act != exp_attr {
			t.Errorf("MatchAttr(%s) for query {{%s}}: expected %t but got %t",
				p.String(), q.String(), exp_attr, act)
		}

		// check time only match
		act = q.MatchTime(p)
		if act != exp_ts {
			t.Errorf("MatchTime(%s) for query {{%s}}: expected %t but got %t",
				p.String(), q.String(), exp_ts, act)
		}

		// check timestamp+attr match
		exp_both := exp_ts && exp_attr
		act = q.Match(p)
		if act != exp_both {
			t.Errorf("Match(%s) for query {{%s}}: expected %t but got %t",
				p.String(), q.String(), exp_both, act)
		}
	}

	fn := func(qa QueryAttr, exp_attr bool) {
		// check a few different time ranges
		fnsub(NewQuery(ts1, ts3, qa), true, exp_attr)
		fnsub(NewQuery(ts1, ts1, qa), true, exp_attr)
		fnsub(NewQuery(ts2, ts3, qa), false, exp_attr)
		fnsub(NewQuery(ts3, ts3, qa), false, exp_attr)
	}

	// things that are true
	t1 := NewQAEqual("color", "blue")
	t2 := NewQAEqual("animal", "moose")
	t3 := NewQARegex("index", `^\d+$`)
	t4 := NewQAEqual("shape", "square")

	// things that are false
	f1 := NewQARegex("color", "^x")
	f2 := NewQARegex("animal", "mo{3,5}se")
	f3 := NewQAEqual("index", "777")
	f4 := NewQARegex("flavor", "sour")

	fn(NewQATrue(), true)
	fn(NewQAExists("color"), true)
	fn(NewQAExists("flavor"), false)
	fn(t1, true)
	fn(t2, true)
	fn(t3, true)
	fn(t4, true)
	fn(f1, false)
	fn(f2, false)
	fn(f3, false)
	fn(f4, false)
	fn(NewQAOr(NewQAExists("flavor"), f1, f2, f3, f4), false)
	fn(NewQAOr(NewQAExists("flavor"), t1, f2, f3, f4), true)
	fn(NewQAOr(NewQAAnd(f1, f2, f3, f4), NewQANot(t1), NewQAOr(NewQANot(f4), t2, t3)), true)
	fn(NewQAOr(NewQAAnd(f1, f2, f3, f4), NewQANot(t4), NewQAOr(NewQANot(t4), f3)), false)

}
