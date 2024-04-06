package engine

import (
	"equinox/internal/core"
	"equinox/internal/query"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getDurMins(i int) time.Duration {
	dur, err := time.ParseDuration(fmt.Sprintf("%dm", i))
	if err != nil {
		panic(err)
	}
	return dur
}

func getPoint(i uint32) *core.Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)

	s := rand.NewSource(ts.Unix()) // always use the same seed
	r := rand.New(s)               // initialize local pseudorandom generator

	animals := [...]string{"cat", "dog", "horse", "pig", "cow"}
	shapes := [...]string{"circle", "square", "rhombus", "rectangle", "triangle", "pentagon"}
	colors := [...]string{"red", "green", "blue", "yellow", "orange", "purple", "pink", "gray", "black", "white"}

	p := core.NewPoint(ts.Add(getDurMins(int(i))))
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

// compares all extracted points from a PointIO to the expectd points
func cmpPointIO(t *testing.T, exp []*core.Point, io PointIO, s string) {
	if !assert.Equal(t, len(exp), io.Len(), "%s length", s) {
		return
	}

	act, err := io.extract()
	if !assert.Nil(t, err) {
		return
	}

	for i := 0; i < len(exp); i++ {
		assert.Equal(t, 0, core.PointCmp(exp[i], act[i]),
			"%s point[%d]\nexpected %s\ngot %s",
			s, i, exp[i].String(), act[i].String())
	}
}

func testPointIOMove(t *testing.T, fact func() PointIO) {
	f := func(n int, i int, j int) {
		ml := fact()
		ps := getPoints(0, n)
		err := ml.Add(ps...)
		assert.Nil(t, err)

		cmpPointIO(t, ps, ml, "original")

		ml2 := fact()

		// first try moving no points
		st := ps[n-1].Ts.Add(getDurMins(3))
		en := ps[n-1].Ts.Add(getDurMins(13))
		m, err := ml.Move(st, en, ml2)
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, 0, m)
		cmpPointIO(t, ps, ml, "source")
		cmpPointIO(t, make([]*core.Point, 0), ml2, "dest")

		// now move the specified points
		m, err = ml.Move(ps[i].Ts, ps[j].Ts, ml2)
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, j-i+1, m)

		psrc := make([]*core.Point, 0)
		if i > 0 {
			psrc = append(psrc, ps[0:i]...)
		}
		if j < n-1 {
			psrc = append(psrc, ps[j+1:n]...)
		}

		pdest := ps[i : j+1]

		cmpPointIO(t, psrc, ml, "source")
		cmpPointIO(t, pdest, ml2, "dest")
	}

	f(5, 0, 0) // first point
	f(5, 0, 3) // first 4 points
	f(5, 2, 4) // last 3 points
	f(5, 0, 4) // all points

	// test larger sizes
	f(100, 31, 65) // mid
	f(100, 0, 88)  // start
	f(100, 28, 99) // end

	f(1000, 314, 659) // mid
	f(1000, 0, 888)   // start
	f(1000, 283, 999) // end
}

func cmpQResults(t *testing.T, q *query.Query, exp []*core.Point, act []*core.Point) {
	if !assert.Equal(t, len(act), len(exp)) {
		return
	}

	// sort by ascending time
	slices.SortFunc(exp, core.PointCmp)
	slices.SortFunc(act, core.PointCmp)

	// now compare one at a time
	for i := 0; i < len(exp); i++ {
		assert.True(t, exp[i].Equal(act[i]))
	}
}

func testQuery(t *testing.T, io PointIO, mints time.Time, maxts time.Time, exp []*core.Point) {

	q := query.NewQuery(mints, maxts, query.True())
	qe, err := io.Search(q)
	assert.Nil(t, err)

	// fetch results in batches
	var results []*core.Point
	batchsize := 10
	for {
		rbatch, err := qe.Fetch(batchsize)
		assert.Nil(t, err)

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

		assert.Equal(t, expsize, len(rbatch))
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

	// should be empty
	assert.Equal(t, 0, io.Len())
	assert.Nil(t, io.First())
	assert.Nil(t, io.Last())

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
			err = io.Add(pbatch...)
			assert.Nil(t, err)
			pbatch = nil

			// make sure first and last are kept updated
			assert.NotNil(t, io.First())
			assert.Equal(t, mints, io.First().Ts)
			assert.NotNil(t, io.Last())
			assert.Equal(t, maxts, io.Last().Ts)
		}
	}

	if len(pbatch) > 0 { // final batch
		io.Add(pbatch...)
		assert.Nil(t, err)
		pbatch = nil

		// make sure first and last are kept updated
		assert.NotNil(t, io.First())
		assert.Equal(t, mints, io.First().Ts)
		assert.NotNil(t, io.Last())
		assert.Equal(t, maxts, io.Last().Ts)
	}

	err = io.Vacuum()
	assert.Nil(t, err)

	// basic query should return all
	testQuery(t, io, mints, maxts, exp)
	noresults := make([]*core.Point, 0)
	testQuery(t, io, mints.Add(getDurMins(-60)), mints.Add(getDurMins(-1)), noresults)
	testQuery(t, io, maxts.Add(getDurMins(1)), maxts.Add(getDurMins(60)), noresults)

}

func testPointIOFull(t *testing.T, fact func() PointIO) {
	testPointIO(t, fact(), 10, 5)
	testPointIO(t, fact(), 10, 10)
	testPointIO(t, fact(), 10, 4)
	testPointIO(t, fact(), 10, 1)
	testPointIO(t, fact(), 100, 9)
	testPointIO(t, fact(), 1000, 49)
	testPointIO(t, fact(), 1000, 50)
}
