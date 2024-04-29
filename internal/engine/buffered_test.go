package engine

import (
	"equinox/internal/core"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createBuffered() PointIO {
	return NewBuffered(NewMemTree(), NewMemList())
}

func TestBufferedQuery(t *testing.T) {
	testPointIOFull(t, func() PointIO { return createBuffered() })
}

func TestBufferedMove(t *testing.T) {
	testPointIOMove(t, func() PointIO { return createBuffered() })
}

func TestBufferedString(t *testing.T) {
	mt := createBuffered()
	assert.Equal(t, "Buffered", mt.Name())

	mt.Add(testGetPoints(5, 2)...)
	exp := `Buffered MemTree to MemList with 15m0s delay
MemTree: {
0: [2024-01-10 23:06:02 +0000 UTC] val[area: -0.958924, temp: 0.283662] attr[animal: pig, color: purple, shape: circle]
1: [2024-01-10 23:07:02 +0000 UTC] val[area: -0.279415, temp: 0.960170] attr[animal: pig, color: purple, shape: circle]
}
MemList: {

}`

	assert.Equal(t, exp, mt.String())
}

func TestBufferedConstructBasic(t *testing.T) {
	mt := createBuffered()
	ps := testGetPoints(0, 10)

	assert.Equal(t, 0, mt.Len())

	runtest := func(p []*core.Point, len int) {
		mt.Add(p...)
		assert.Equal(t, len, mt.Len())
	}

	runtest(make([]*core.Point, 0), 0)
	runtest(ps[3:5], 2)
	runtest(ps[0:2], 4)
	runtest(ps[2:3], 5)
	runtest(ps[5:7], 7)
	runtest(ps[5:7], 7) // we should not be allowed to add duplicates
	runtest([]*core.Point{ps[9], ps[8], ps[7]}, 10)

}

func TestBufferedFlush(t *testing.T) {
	mt := NewBuffered(NewMemTree(), NewMemList())
	ps := testGetPoints(0, 10)

	assert.Equal(t, 0, mt.Len())

	testlens := func(buflen int, archlen int, s string) {
		assert.Equal(t, buflen, mt.Buf.Len(), s)
		assert.Equal(t, archlen, mt.Archive.Len(), s)
		assert.Equal(t, buflen+archlen, mt.Len(), s)
	}
	var err error

	testlens(0, 0, "empty to start")
	err = mt.Flush()
	assert.Nil(t, err)

	// add 3 points and flush them through
	err = mt.Add(ps[0:3]...)
	assert.Nil(t, err)
	testlens(3, 0, "all 3 points should be in buffer")
	err = mt.Flush()
	assert.Nil(t, err)
	testlens(0, 3, "3 points flushed to archive")

	// add 4 more points and flush them through
	err = mt.Add(ps[3:7]...)
	assert.Nil(t, err)
	testlens(4, 3, "4 new points in buffer, 3 still in archive")
	err = mt.Flush()
	assert.Nil(t, err)
	testlens(0, 7, "all 7 points in archive")

	// now add a current point that shouldn't get moved
	p := testGetPoint(0)
	p.Ts = time.Now().UTC()
	err = mt.Add()
	assert.Nil(t, err)
	testlens(1, 7, "7 points in archive, 1 new in buffer")
	err = mt.Flush()
	assert.Nil(t, err)
	testlens(1, 7, "7 points in archive, 1 still in buffer after flush")
}
