package equinox

import (
	"testing"
)

func TestMemListConstructBasic(t *testing.T) {
	ml := NewMemList()
	ps := getPoints(0, 10)
	var err error

	ml.Add(ps[3:5])
	err = ml.validate()
	if err != nil {
		t.Fatalf("validation failed: %s", err.Error())
	}
	//t.Logf("after adding 3:4: %s", ml.String())
	if ml.Len() != 2 {
		t.Fatalf("incorrect length expected 2 got %d", ml.Len())
	}

	ml.Add(ps[0:2])
	err = ml.validate()
	if err != nil {
		t.Fatalf("validation failed: %s", err.Error())
	}
	if ml.Len() != 4 {
		t.Fatalf("incorrect length expected 4 got %d", ml.Len())
	}
	//t.Logf("after adding 0:1: %s", ml.String())

	ml.Add(ps[2:3])
	err = ml.validate()
	if err != nil {
		t.Fatalf("Validation failed: %s", err.Error())
	}
	if ml.Len() != 5 {
		t.Fatalf("incorrect length expected 5 got %d", ml.Len())
	}
	//t.Logf("after adding 2: %s", ml.String())

	ml.Add(ps[5:7])
	err = ml.validate()
	if err != nil {
		t.Fatalf("Validation failed: %s", err.Error())
	}
	if ml.Len() != 7 {
		t.Fatalf("incorrect length expected 7 got %d", ml.Len())
	}
	//t.Logf("after adding 5:6: %s", ml.String())

	ml.Add([]*Point{ps[9], ps[8], ps[7]})
	err = ml.validate()
	if err != nil {
		t.Fatalf("Validation failed: %s", err.Error())
	}
	if ml.Len() != 10 {
		t.Fatalf("incorrect length expected 10 got %d", ml.Len())
	}
	//t.Logf("after adding 9:7: %s", ml.String())
}

func TestMemListConstructBatches(t *testing.T) {
	ml := NewMemList()
	ps := getPoints(0, 100)
	batch := 10
	var err error
	var pbatch []*Point

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
	ml := NewMemList()
	testPointIO(t, ml, 10, 5)
}
