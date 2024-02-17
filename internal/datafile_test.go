package equinox

import (
	"os"
	"testing"
)

func tempFileName() (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "equinox-test-*")
	if err != nil {
		return "", err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	return f.Name(), nil
}

func TestDFNew(t *testing.T) {
	fn, err := tempFileName()
	if err != nil {
		t.Fatalf("failed to get temp file name: %s", err.Error())
	}
	defer os.Remove(fn)

	ser := NewSerializer()

	b, _ := ser.Serialize(getPoint(0))

	df, err := OpenNewDF(fn, ser, uint32(len(b)))
	if err != nil {
		t.Fatalf("failed to open file %s with recsize %d: %s", fn, len(b), err.Error())
	}
	defer df.Close()

	var i uint32
	for i = 0; i < 10; i++ {
		p := getPoint(i)
		err := df.Write(i, p)
		if err != nil {
			t.Errorf("error while reading point %d: %s", i, err.Error())
		}
	}

	// read back from the same DataFile
	for i = 0; i < 10; i++ {
		p := getPoint(i)
		p2, err := df.Read(i)

		if err != nil {
			t.Errorf("error while reading point %d: %s", i, err.Error())
			continue
		}

		if !p.Equal(p2) {
			t.Errorf("expected %s, got %s", p.String(), p2.String())
		}
	}

	// read back from a new DataFile object
	df2, err := OpenExistingDF(fn, ser)
	if err != nil {
		t.Fatalf("failed to open existing file %s: %s", fn, err.Error())
	}

	for i = 9; i >= 1; i-- {
		p := getPoint(i)
		p2, err := df2.Read(i)

		if err != nil {
			t.Errorf("error while reading point %d: %s", i, err.Error())
			continue
		}

		if !p.Equal(p2) {
			t.Errorf("expected %s, got %s", p.String(), p2.String())
		}
	}

}

func TestDFMissing(t *testing.T) {
	fn, err := tempFileName()
	if err != nil {
		t.Fatalf("failed to get temp file name: %s", err.Error())
	}
	defer os.Remove(fn)

	ser := NewSerializer()

	_, err = OpenExistingDF(fn, ser)
	if err == nil {
		t.Fatalf("expected to fail to open non-existent file %s", fn)
	}
}

// tests writing/reading non-sequential points
func TestDFNonseq(t *testing.T) {
	fn, err := tempFileName()
	if err != nil {
		t.Fatalf("failed to get temp file name: %s", err.Error())
	}
	defer os.Remove(fn)

	ser := NewSerializer()

	b, _ := ser.Serialize(getPoint(0))

	df, err := OpenNewDF(fn, ser, uint32(len(b)))
	if err != nil {
		t.Fatalf("failed to open file %s with recsize %d: %s", fn, len(b), err.Error())
	}
	defer df.Close()

	// write points 0, 2
	df.Write(0, getPoint(0))
	df.Write(2, getPoint(2))

	// point 0 and 2 should work
	for i := 0; i < 3; i += 2 {
		p := getPoint(uint32(i))
		p2, err := df.Read(uint32(i))

		if err != nil {
			t.Errorf("error while reading point %d: %s", i, err.Error())
		}

		if !p.Equal(p2) {
			t.Errorf("expected %s, got %s", p.String(), p2.String())
		}
	}

	// point 3+ should fail
	for i := 3; i < 600; i += 100 {
		p2, err := df.Read(uint32(i))
		if err == nil {
			t.Errorf("should have failed to read point %d; got %s", i, p2.String())
		} else {
			//t.Logf("correctly got error when reading point %d: %s", i, err.Error())
		}
	}

	// point 1 should fail
	for i := 1; i < 2; i++ {
		p2, err := df.Read(uint32(i))
		if err == nil {
			t.Errorf("should have failed to read point %d; got %s", i, p2.String())
		} else {
			//t.Logf("correctly got error when reading point %d: %s", i, err.Error())
		}
	}
}
