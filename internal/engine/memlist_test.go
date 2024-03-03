package engine

import (
	"equinox/internal/core"
	"testing"
)

func TestMemListConstructBasic(t *testing.T) {
	ml := NewMemList()
	ps := getPoints(0, 10)
	var err error

	runtest := func(p []*core.Point, len int) {
		ml.Add(p)
		err = ml.validate()
		if err != nil {
			t.Fatalf("validation failed: %s", err.Error())
		}
		if ml.Len() != len {
			t.Fatalf("incorrect length expected %d got %d", len, ml.Len())
		}
	}

	runtest(ps[3:5], 2)
	runtest(ps[0:2], 4)
	runtest(ps[2:3], 5)
	runtest(ps[5:7], 7)
	runtest([]*core.Point{ps[9], ps[8], ps[7]}, 10)
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
			err = ml.Add(pbatch)
			if err != nil {
				t.Fatalf("unexpected error when adding %d points: %s", len(pbatch), err.Error())
			}
			err = ml.validate()
			if err != nil {
				t.Fatalf("Validation failed: %s", err.Error())
			}
			pbatch = nil
		}
	}

	if len(pbatch) > 0 { // final batch
		ml.Add(pbatch)
		if err != nil {
			t.Fatalf("unexpected error when adding %d points: %s", len(pbatch), err.Error())
		}
		err = ml.validate()
		if err != nil {
			t.Fatalf("Validation failed: %s", err.Error())
		}
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
