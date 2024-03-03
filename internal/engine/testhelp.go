package engine

import (
	"equinox/internal/core"
	"equinox/internal/query"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"
)

func testAttrsToString(attrs map[string]string) string {
	var attr []string

	for k, v := range attrs {
		attr = append(attr, k+": "+v)
	}
	sort.Strings(attr) // ensure consistent output

	return strings.Join(attr, ", ")
}

func testGetAttrs() map[string]string {
	r := make(map[string]string)
	r["color"] = "blue"
	r["animal"] = "moose"
	r["shape"] = "square"
	r["index"] = "74"
	return r
}

func getPoint(i uint32) *core.Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	dur, err := time.ParseDuration(fmt.Sprintf("%dm", i))
	if err != nil {
		panic(err)
	}

	s := rand.NewSource(ts.Unix()) // always use the same seed
	r := rand.New(s)               // initialize local pseudorandom generator

	animals := [...]string{"cat", "dog", "horse", "pig", "cow"}
	shapes := [...]string{"circle", "square", "rhombus", "rectangle", "triangle", "pentagon"}
	colors := [...]string{"red", "green", "blue", "yellow", "orange", "purple", "pink", "gray", "black", "white"}

	p := core.NewPoint(ts.Add(dur))
	p.Attrs["color"] = colors[r.Intn(len(colors))]
	p.Attrs["shape"] = shapes[r.Intn(len(shapes))]
	p.Attrs["animal"] = animals[r.Intn(len(animals))]
	p.Vals["area"] = math.Sin(float64(i))
	p.Vals["temp"] = math.Cos(float64(i))
	return p
}

// gets n points starting at a, in random order
func getPoints(a uint32, n int) []*core.Point {
	var ps []*core.Point
	if n == 0 {
		return ps
	}

	for i := 0; i < n; i++ {
		ps = append(ps, getPoint(uint32(i)+a))
	}

	return ps
}

// gets n points starting at a, in random order
func getPointsShuffle(a uint32, n int) []*core.Point {
	ps := getPoints(a, n)
	rand.Shuffle(n, func(i, j int) { ps[i], ps[j] = ps[j], ps[i] })
	return ps
}

func cmpQResults(t *testing.T, q *query.Query, exp []*core.Point, act []*core.Point) {
	if len(exp) != len(act) {
		t.Fatalf("unexpected # of results for query %s: expected %d got %d", q.String(), len(exp), len(act))
	}

	// sort by ascending time
	slices.SortFunc(exp, core.PointCmp)
	slices.SortFunc(act, core.PointCmp)

	// now compare one at a time
	for i := 0; i < len(exp); i++ {
		if !exp[i].Equal(act[i]) {
			t.Errorf("unexpected point returned; got %s expected %s", act[i].String(), exp[i].String())
		}

	}
}

func testQuery(t *testing.T, io PointIO, mints time.Time, maxts time.Time, exp []*core.Point) {

	q := query.NewQuery(mints, maxts, query.True())
	qe, err := io.Search(q)
	if err != nil {
		t.Fatalf("unexpected error when initiating query %s: %s", q.String(), err.Error())
	}

	// fetch results in batches
	var results []*core.Point
	batchsize := 10
	for {
		rbatch, err := qe.Fetch(batchsize)

		if err != nil {
			t.Fatalf("unexpected error when fetching %d results for query %s: %s", batchsize, q.String(), err.Error())
			return
		}

		if false {
			t.Logf("===== Got batch of %d points ====", len(rbatch))
			for i := 0; i < len(rbatch); i++ {
				t.Logf("[%d] %s", i, rbatch[i].String())
			}
		}

		// read the last one
		if len(rbatch) == 0 {
			break
		}

		// how many we expected back; should be batchsize unless there aren't that many left
		expsize := batchsize
		if len(results)+expsize > len(exp) {
			expsize = len(exp) - len(results)
		}

		if expsize != len(rbatch) {
			t.Fatalf("unexpected # results fetched for query %s: expected %d got %d", q.String(), expsize, len(rbatch))
			break
		}

		results = append(results, rbatch...)
	}

	cmpQResults(t, q, exp, results)
}

func testPointIO(t *testing.T, io PointIO, n int, batch int) {
	exp := getPointsShuffle(0, n)
	t.Logf("testing %s with %d points and batch size %d", io.Name(), n, batch)

	var err error
	var mints, maxts time.Time
	var pbatch []*core.Point

	for i := 0; i < len(exp); i++ {
		p := exp[i]

		// remember max and min timestamps
		if i == 0 {
			mints = p.Ts
			maxts = p.Ts
		} else {
			if p.Ts.Before(mints) {
				mints = p.Ts
			}
			if p.Ts.After(maxts) {
				maxts = p.Ts
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
