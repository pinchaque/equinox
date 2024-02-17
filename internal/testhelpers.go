package equinox

import (
	"fmt"
	"math"
	"math/rand"
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

func getPoint(i uint32) *Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 3, time.UTC)
	dur, err := time.ParseDuration(fmt.Sprintf("%dm", i))
	if err != nil {
		panic(err)
	}

	s := rand.NewSource(ts.Unix()) // always use the same seed
	r := rand.New(s)               // initialize local pseudorandom generator

	animals := [...]string{"cat", "dog", "horse", "pig", "cow"}

	p := NewPoint(ts.Add(dur))
	p.attrs["color"] = fmt.Sprintf("clr%d", i)
	p.attrs["shape"] = fmt.Sprintf("shp%d", i)
	p.attrs["animal"] = animals[r.Intn(len(animals))]
	p.vals["area"] = math.Sin(float64(i))
	p.vals["temp"] = math.Cos(float64(i))
	return p
}

func testPointIO(t *testing.T, io PointIO, n int, batch int) {
	var ps []*Point
	var err error
	var mints, maxts time.Time

	for i := 0; i < n; i++ {
		p := getPoint(uint32(i))

		// remember max and min timestamps
		if i == 0 {
			mints = p.ts
			maxts = p.ts
		} else {
			if p.ts.Before(mints) {
				mints = p.ts
			}
			if p.ts.After(maxts) {
				maxts = p.ts
			}
		}

		ps = append(ps, p)
		if len(ps) >= batch { // add in batches
			err = io.Add(ps)
			if err != nil {
				t.Fatalf("unexpected error when adding %d points: %s", len(ps), err.Error())
			}

			ps = nil
		}
	}

	if len(ps) > 0 { // final batch
		io.Add(ps)
		if err != nil {
			t.Fatalf("unexpected error when adding %d points: %s", len(ps), err.Error())
		}
		ps = nil
	}

	err = io.Vacuum()
	if err != nil {
		t.Fatalf("unexpected error when vacuuming: %s", err.Error())
	}

	q := NewQuery(mints, maxts)
	cur, err := io.Search(q)
	if err != nil {
		t.Fatalf("unexpected error when initiating query %s: %s", q.String(), err.Error())
	}

	// test out cursor
	var results []*Point
	j := 1
	results, err = cur.Fetch(j)
	if err != nil {
		t.Fatalf("unexpected error when fetching %d results for query %s: %s", j, q.String(), err.Error())
	}

	// TODO: validate the results

}
