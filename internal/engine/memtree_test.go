package engine

import (
	"equinox/internal/core"
	"testing"
)

func TestMemTreeQuery(t *testing.T) {
	testPointIO(t, NewMemTree(), 10, 5)
	testPointIO(t, NewMemTree(), 10, 10)
	testPointIO(t, NewMemTree(), 10, 4)
	testPointIO(t, NewMemTree(), 10, 1)
	testPointIO(t, NewMemTree(), 100, 9)
	testPointIO(t, NewMemTree(), 100, 10)
	testPointIO(t, NewMemTree(), 1000, 49)
	testPointIO(t, NewMemTree(), 1000, 50)
}

func TestMemTreeString(t *testing.T) {
	mt := NewMemTree()

	exp := "MemTree"
	if mt.Name() != exp {
		t.Errorf("incorrect Name: expected %s got %s", exp, mt.Name())
	}

	mt.Add(getPoints(5, 2))
	exp = `MemTree: {
0: [2024-01-10 23:06:02 +0000 UTC] val[area: -0.958924, temp: 0.283662] attr[animal: pig, color: purple, shape: circle]
1: [2024-01-10 23:07:02 +0000 UTC] val[area: -0.279415, temp: 0.960170] attr[animal: pig, color: purple, shape: circle]
}`
	act := mt.String()
	if exp != act {
		t.Errorf("incorect MemTree string: expected '%s' got '%s'", exp, act)
	}
}
func TestMemTreeConstructBasic(t *testing.T) {
	mt := NewMemTree()
	ps := getPoints(0, 10)

	if mt.Len() != 0 {
		t.Fatalf("incorrect length expected 0 got %d", mt.Len())
	}

	runtest := func(p []*core.Point, len int) {
		mt.Add(p)
		if mt.Len() != len {
			t.Fatalf("incorrect length expected %d got %d", len, mt.Len())
		}
	}

	runtest(make([]*core.Point, 0), 0)
	runtest(ps[3:5], 2)
	runtest(ps[0:2], 4)
	runtest(ps[2:3], 5)
	runtest(ps[5:7], 7)
	runtest(ps[5:7], 7) // we should not be allowed to add duplicates
	runtest([]*core.Point{ps[9], ps[8], ps[7]}, 10)

}
