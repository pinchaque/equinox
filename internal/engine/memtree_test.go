package engine

import (
	"equinox/internal/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemTreeQuery(t *testing.T) {
	testPointIOFull(t, func() PointIO { return NewMemTree() })
}

func TestMemTreeString(t *testing.T) {
	mt := NewMemTree()
	assert.Equal(t, "MemTree", mt.Name())

	mt.Add(getPoints(5, 2)...)
	exp := `MemTree: {
0: [2024-01-10 23:06:02 +0000 UTC] val[area: -0.958924, temp: 0.283662] attr[animal: pig, color: purple, shape: circle]
1: [2024-01-10 23:07:02 +0000 UTC] val[area: -0.279415, temp: 0.960170] attr[animal: pig, color: purple, shape: circle]
}`

	assert.Equal(t, exp, mt.String())
}
func TestMemTreeConstructBasic(t *testing.T) {
	mt := NewMemTree()
	ps := getPoints(0, 10)

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
