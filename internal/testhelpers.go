package equinox

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"testing"
	"time"
)

func getPoint(i uint32) *Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 3, time.UTC)
	dur, err := time.ParseDuration(fmt.Sprintf("%dm", i))
	if err != nil {
		panic(err)
	}

	s := rand.NewSource(ts.Unix()) // always use the same seed
	r := rand.New(s)               // initialize local pseudorandom generator

	animals := [...]string{"cat", "dog", "horse", "pig", "cow"}
	shapes := [...]string{"circle", "square", "rhombus", "rectangle", "triangle", "pentagon"}
	colors := [...]string{"red", "green", "blue", "yellow", "orange", "purple", "pink", "gray", "black", "white"}

	p := NewPoint(ts.Add(dur))
	p.attrs["color"] = colors[r.Intn(len(colors))]
	p.attrs["shape"] = shapes[r.Intn(len(shapes))]
	p.attrs["animal"] = animals[r.Intn(len(animals))]
	p.vals["area"] = math.Sin(float64(i))
	p.vals["temp"] = math.Cos(float64(i))
	return p
}

// gets n points starting at a, in random order
func getPoints(a uint32, n int) []*Point {
	var ps []*Point
	if n == 0 {
		return ps
	}

	for i := 0; i < n; i++ {
		ps = append(ps, getPoint(uint32(i)+a))
	}

	rand.Shuffle(n, func(i, j int) { ps[i], ps[j] = ps[j], ps[i] })

	return ps
}

func cmpQResults(t *testing.T, q *Query, exp []*Point, act []*Point) {
	if len(exp) != len(act) {
		t.Fatalf("unexpected # of results for query %s: expected %d got %d", q.String(), len(exp), len(act))
	} else {
		t.Logf("num results for query %s: expected %d got %d", q.String(), len(exp), len(act))
	}

	// sort by ascending time
	fn := func(a, b *Point) int { return int(a.ts.UnixMicro() - b.ts.UnixMicro()) }
	slices.SortFunc(exp, fn)
	slices.SortFunc(act, fn)

	// now compare one at a time
	for i := 0; i < len(exp); i++ {
		if !exp[i].Equal(act[i]) {
			t.Errorf("unexpected point returned; got %s expected %s", act[i].String(), exp[i].String())
		} else {
			t.Logf("expected point returned; got %s expected %s", act[i].String(), exp[i].String())
		}

	}
}

func testQuery(t *testing.T, io PointIO, mints time.Time, maxts time.Time, exp []*Point) {

	q := NewQuery(mints, maxts)
	cur, err := io.Search(q)
	if err != nil {
		t.Fatalf("unexpected error when initiating query %s: %s", q.String(), err.Error())
	}

	var results []*Point
	j := 1
	// TODO figure out how many to fetch and validate that it got how many we asked for
	results, err = cur.Fetch(j)
	if err != nil {
		t.Fatalf("unexpected error when fetching %d results for query %s: %s", j, q.String(), err.Error())
	}

	cmpQResults(t, q, exp, results)
}

func testPointIO(t *testing.T, io PointIO, n int, batch int) {
	exp := getPoints(0, n)
	t.Logf("testing %s with %d points and batch size %d", io.String(), n, batch)
	t.Logf("generated %d points", len(exp))

	var err error
	var mints, maxts time.Time
	var pbatch []*Point

	for i := 0; i < len(exp); i++ {
		p := exp[i]

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

		pbatch = append(pbatch, p)
		if len(pbatch) >= batch { // add in batches
			err = io.Add(pbatch)
			if err != nil {
				t.Fatalf("unexpected error when adding %d points: %s", len(pbatch), err.Error())
			}

			pbatch = nil
		}
	}

	if len(pbatch) > 0 { // final batch
		io.Add(pbatch)
		if err != nil {
			t.Fatalf("unexpected error when adding %d points: %s", len(pbatch), err.Error())
		}
		pbatch = nil
	}

	err = io.Vacuum()
	if err != nil {
		t.Fatalf("unexpected error when vacuuming: %s", err.Error())
	}

	// basic query should return all
	testQuery(t, io, mints, maxts, exp)
}
