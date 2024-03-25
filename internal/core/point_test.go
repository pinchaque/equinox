package core

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPointCreate(t *testing.T) {

	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := NewPoint(ts)
	assert.Equal(t, ts, p.Ts)
	assert.NotNil(t, p.Vals)
	assert.Equal(t, 0, len(p.Vals))
	assert.NotNil(t, p.Attrs)
	assert.Equal(t, 0, len(p.Attrs))
	assert.NotEqual(t, "", p.Id.String())

	// generate new Id and make sure it is different
	oldId := p.Id.String()
	p.GenerateId()
	assert.NotEqual(t, oldId, p.Id.String())
}
func TestPointCreateEmptyId(t *testing.T) {

	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := NewPointEmptyId(ts)

	assert.Equal(t, ts, p.Ts)
	assert.NotNil(t, p.Vals)
	assert.Equal(t, 0, len(p.Vals))
	assert.NotNil(t, p.Attrs)
	assert.Equal(t, 0, len(p.Attrs))
	assert.Nil(t, p.Id)
}

func TestPointCreateEmpty(t *testing.T) {

	ts := time.Time{}
	p := NewPointEmpty()

	assert.Equal(t, ts, p.Ts)
	assert.NotNil(t, p.Vals)
	assert.Equal(t, 0, len(p.Vals))
	assert.NotNil(t, p.Attrs)
	assert.Equal(t, 0, len(p.Attrs))
	assert.Nil(t, p.Id)
}

func testNewPointComplete() *Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := NewPoint(ts)
	p.Attrs["shape"] = "square"
	p.Attrs["color"] = "red"
	p.Vals["area"] = 43.1
	p.Vals["temp"] = 21.1
	return p
}

func TestPointString(t *testing.T) {
	p := testNewPointComplete()
	exp := "[2024-01-10 23:01:02 +0000 UTC] val[area: 43.100000, temp: 21.100000] attr[color: red, shape: square]"
	assert.Equal(t, exp, p.String())

}

func TestPointJSON(t *testing.T) {
	p := testNewPointComplete()
	p.Ts = time.Date(2024, 01, 10, 23, 1, 2, 123456789, time.UTC) // add microsecs
	p.Id.val = 485782                                             // need consistent ID
	b, err := json.Marshal(p)
	assert.NoError(t, err)
	exp := `{"Ts":"2024-01-10T23:01:02.123456789Z","Vals":{"area":43.1,"temp":21.1},"Attrs":{"color":"red","shape":"square"},"Id":"AAAAAAAHaZY="}`
	assert.Equal(t, exp, string(b))

	// now try unmarshaling
	p2 := &Point{} // empty point
	err = json.Unmarshal(b, p2)
	assert.NoError(t, err)
	assert.Equal(t, true, p.Equal(p2), "Orig Point:\n%s\nUnmarshaled:\n%s\n", p.String(), p2.String())
}
func TestPointJSONEmpty(t *testing.T) {
	p := NewPointEmpty()
	b, err := json.Marshal(p)
	assert.NoError(t, err)
	exp := `{"Ts":"0001-01-01T00:00:00Z","Vals":{},"Attrs":{},"Id":null}`
	assert.Equal(t, exp, string(b))

	// now try unmarshaling
	p2 := &Point{} // empty point
	err = json.Unmarshal(b, p2)
	assert.NoError(t, err)
	assert.Nil(t, p2.Id)
	assert.Equal(t, true, p.Equal(p2), "Orig Point:\n%s\nUnmarshaled:\n%s\n", p.String(), p2.String())
}

func TestPointCreateComplete(t *testing.T) {
	p := testNewPointComplete()
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)

	assert.Equal(t, ts, p.Ts)
	assert.NotNil(t, p.Vals)
	assert.Equal(t, 2, len(p.Vals))
	assert.NotNil(t, p.Attrs)
	assert.Equal(t, 2, len(p.Attrs))
	assert.NotNil(t, p.Id)
}

func TestPointEqual(t *testing.T) {
	p1 := testNewPointComplete()

	cmp := func(pt1 *Point, pt2 *Point, exp int) {
		act := PointCmp(pt1, pt2)
		assert.Equal(t, exp, act)
	}

	{ // basic equality
		p2 := testNewPointComplete()
		assert.True(t, p1.Equal(p2))
		cmp(p1, p2, 0)
	}

	{ // different timestamp
		p2 := testNewPointComplete()
		p2.Ts = p2.Ts.AddDate(0, 0, 1)
		assert.False(t, p1.Equal(p2))
		cmp(p1, p2, -1)
		cmp(p2, p1, 1)
	}

	{ // changed value
		p2 := testNewPointComplete()
		cmp(p1, p2, 0)
		p2.Vals["area"] = 43.1004
		cmp(p1, p2, 0) // only timestamp matters
		assert.False(t, p1.Equal(p2))

		assert.True(t, p1.EqualTol(p2, 0.1))
		assert.True(t, p1.EqualTol(p2, 0.01))
		assert.True(t, p1.EqualTol(p2, 0.001))
		assert.True(t, p1.EqualTol(p2, 0.0001))
		assert.True(t, p1.EqualTol(p2, 0.00001))
		assert.False(t, p1.EqualTol(p2, 0.000001))
		assert.False(t, p1.EqualTol(p2, 0.0000001))

		// test 0.00 value
		p2.Vals["area"] = 0.000
		assert.True(t, p2.EqualTol(p2, 0.001))
	}

	{ // add value
		p2 := testNewPointComplete()
		p2.Vals["area2"] = 49.999
		assert.False(t, p1.Equal(p2))
	}

	{ // delete value
		p2 := testNewPointComplete()
		delete(p2.Vals, "area")
		assert.False(t, p1.Equal(p2))
	}

	{ // changed attr
		p2 := testNewPointComplete()
		cmp(p1, p2, 0)
		p2.Attrs["color"] = "blue"
		cmp(p1, p2, 0) // only timestamp matters
		assert.False(t, p1.Equal(p2))
	}

	{ // add attr
		p2 := testNewPointComplete()
		cmp(p1, p2, 0)
		p2.Attrs["color2"] = "blue"
		cmp(p1, p2, 0) // only timestamp matters
		assert.False(t, p1.Equal(p2))
	}

	{ // delete attr
		p2 := testNewPointComplete()
		delete(p2.Attrs, "color")
		assert.False(t, p1.Equal(p2))
	}
}

func TestPointIdentical(t *testing.T) {
	/*
		ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
		p := NewPoint(ts)
		p.Attrs["shape"] = "square"
		p.Attrs["color"] = "red"
		p.Vals["area"] = 43.1
		p.Vals["temp"] = 21.1
	*/

	p1 := testNewPointComplete()
	p2 := testNewPointComplete()

	fn := func(pt1 *Point, pt2 *Point, exp bool) {
		act := pt1.Identical(pt2)
		assert.Equal(t, exp, act)
	}

	fn(p1, p2, false)
	fn(p2, p1, false)
	fn(p1, p1, true)
	fn(p2, p2, true)

	// make them have the same IDs so now should be identical
	p1.Id = p2.Id
	fn(p1, p2, true)
	fn(p2, p1, true)
	fn(p1, p1, true)
	fn(p2, p2, true)

	// now change an attribute => no longer identical
	p1.Attrs["color"] = "mauve"
	fn(p1, p2, false)
	fn(p2, p1, false)
	fn(p1, p1, true)
	fn(p2, p2, true)
}

func TestPointClone(t *testing.T) {
	/*
		ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
		p := NewPoint(ts)
		p.Attrs["shape"] = "square"
		p.Attrs["color"] = "red"
		p.Vals["area"] = 43.1
		p.Vals["temp"] = 21.1
	*/

	p1 := testNewPointComplete()
	p2 := p1.Clone()

	// these should be identical and equal
	assert.True(t, p1.Equal(p2))
	assert.True(t, p1.Identical(p2))

	// if we change one it shouldn't affect the other
	p1.Attrs["shape"] = "triangle"
	assert.Equal(t, "square", p2.Attrs["shape"])

	p1.Vals["area"] = 1.0
	assert.Equal(t, 43.1, p2.Vals["area"])

	p1.Ts = time.Date(2024, 01, 11, 23, 1, 2, 0, time.UTC)
	assert.Equal(t, "2024-01-10 23:01:02 +0000 UTC", p2.Ts.String())

	p1.Id.val = 3
	assert.NotEqual(t, 3, p2.Id.val)
}
