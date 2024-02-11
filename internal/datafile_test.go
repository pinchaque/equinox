package equinox

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"
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

func getPoint(i int) *Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 3, time.UTC)
	dur, err := time.ParseDuration(fmt.Sprintf("%dm", i))
	if err != nil {
		panic(err)
	}
	p := NewPoint(ts.Add(dur))
	p.attrs["color"] = fmt.Sprintf("clr%d", i)
	p.attrs["shape"] = fmt.Sprintf("shp%d", i)
	p.vals["area"] = math.Sin(float64(i))
	p.vals["temp"] = math.Cos(float64(i))
	return p
}

func TestDFWriteRead(t *testing.T) {
	fn, err := tempFileName()
	if err != nil {
		t.Fatalf("failed to get temp file name: %s", err.Error())
	}
	//defer os.Remove(fn)
	t.Logf("temp file: %s", fn)

	df := OpenFile(fn)
	defer df.CloseFile()

	for i := 0; i < 10; i++ {
		p := getPoint(i)
		df.Write(i, p)
	}

	for i := 0; i < 10; i++ {
		p := getPoint(i)
		p2, err := df.Read(i)

		if err != nil {
			t.Errorf("error while reading poing %d: %s", i, err.Error())
		}

		if !p.Equal(p2) {
			t.Errorf("expected %s, got %s", p.String(), p2.String())
		}
	}
}
