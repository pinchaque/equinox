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
	testPointIOFull(t, func() PointIO { return NewMemList() })
}

func TestMemListMove(t *testing.T) {
	testPointIOMove(t, func() PointIO { return NewMemList() })
}
