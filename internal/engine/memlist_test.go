package engine

import (
	"equinox/internal/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemListConstructBasic(t *testing.T) {
	ml := NewMemList()
	ps := getPoints(0, 10)
	var err error
	assert.Equal(t, 0, ml.Len())

	runtest := func(p []*core.Point, len int) {
		ml.Add(p...)
		err = ml.validate()
		assert.Nil(t, err)
		assert.Equal(t, len, ml.Len())
	}

	runtest(make([]*core.Point, 0), 0)
	runtest(ps[3:5], 2)
	runtest(ps[0:2], 4)
	runtest(ps[2:3], 5)
	runtest(ps[5:7], 7)
	// TODO: fix MemList so that it doesn't allow duplicate points
	// runtest(ps[5:7], 7) // we should not be allowed to add duplicates
	runtest([]*core.Point{ps[9], ps[8], ps[7]}, 10)
}

func TestMemListString(t *testing.T) {
	ml := NewMemList()

	assert.Equal(t, "MemList", ml.Name())

	ml.Add(getPoints(5, 2)...)
	exp := `MemList: {
0: [2024-01-10 23:06:02 +0000 UTC] val[area: -0.958924, temp: 0.283662] attr[animal: pig, color: purple, shape: circle]
1: [2024-01-10 23:07:02 +0000 UTC] val[area: -0.279415, temp: 0.960170] attr[animal: pig, color: purple, shape: circle]
}`
	assert.Equal(t, exp, ml.String())
}

func TestMemListConstructBatches(t *testing.T) {
	ml := NewMemList()
	ps := getPoints(0, 100)
	batch := 10
	var err error
	var pbatch []*core.Point

	for i := 0; i < len(ps); i++ {
		p := ps[i]

		pbatch = append(pbatch, p)
		if len(pbatch) >= batch { // add in batches
			err = ml.Add(pbatch...)
			assert.Nil(t, err)
			err = ml.validate()
			assert.Nil(t, err)
			pbatch = nil
		}
	}

	if len(pbatch) > 0 { // final batch
		ml.Add(pbatch...)
		assert.Nil(t, err)
		err = ml.validate()
		assert.Nil(t, err)
		pbatch = nil
	}
}

func TestMemListQuery(t *testing.T) {
	testPointIO(t, NewMemList(), 10, 5)
	testPointIO(t, NewMemList(), 10, 10)
	testPointIO(t, NewMemList(), 10, 4)
	testPointIO(t, NewMemList(), 10, 1)
	testPointIO(t, NewMemList(), 100, 9)
	testPointIO(t, NewMemList(), 1000, 49)
	testPointIO(t, NewMemList(), 1000, 50)
}

func memListCmp(t *testing.T, exp []*core.Point, ml *MemList, s string) {
	if !assert.Equal(t, len(exp), ml.Len(), "%s length", s) {
		return
	}

	i := 0
	for e := ml.buf.Front(); e != nil; e = e.Next() {
		p := e.Value.(*core.Point)
		assert.Equal(t, 0, core.PointCmp(exp[i], p), "%s point %d", s, i)
		i++
	}
}

func TestMemListMove(t *testing.T) {
	f := func(n int, i int, j int) {
		ml := NewMemList()
		ps := getPoints(0, n)
		err := ml.Add(ps...)
		assert.Nil(t, err)

		memListCmp(t, ps, ml, "original")

		ml2 := NewMemList()

		// first try moving no points
		st := ps[n-1].Ts.Add(getDurMins(3))
		en := ps[n-1].Ts.Add(getDurMins(13))
		m, err := ml.Move(st, en, ml2)
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, 0, m)
		memListCmp(t, ps, ml, "source")
		memListCmp(t, make([]*core.Point, 0), ml2, "dest")

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

		memListCmp(t, psrc, ml, "source")
		memListCmp(t, pdest, ml2, "dest")
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
