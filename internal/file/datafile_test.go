package file

import (
	"equinox/internal/core"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getPoint(i uint32) *core.Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	dur, err := time.ParseDuration(fmt.Sprintf("%dm", i))
	if err != nil {
		panic(err)
	}

	s := rand.NewSource(ts.Unix()) // always use the same seed
	r := rand.New(s)               // initialize local pseudorandom generator

	animals := [...]string{"cat", "dog", "horse", "pig", "cow"}
	shapes := [...]string{"circle", "square", "rhombus", "rectangle", "triangle", "pentagon"}
	colors := [...]string{"red", "green", "blue", "yellow", "orange", "purple", "pink", "gray", "black", "white"}

	p := core.NewPoint(ts.Add(dur))
	p.Attrs["color"] = colors[r.Intn(len(colors))]
	p.Attrs["shape"] = shapes[r.Intn(len(shapes))]
	p.Attrs["animal"] = animals[r.Intn(len(animals))]
	p.Vals["area"] = math.Sin(float64(i))
	p.Vals["temp"] = math.Cos(float64(i))
	return p
}

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
	assert.Nil(t, err)
	defer os.Remove(fn)

	ser := NewSerializer()

	b, _ := ser.Serialize(getPoint(0))

	df, err := OpenNewDF(fn, ser, uint32(len(b)))
	assert.Nil(t, err)
	defer df.Close()

	var i uint32
	for i = 0; i < 10; i++ {
		p := getPoint(i)
		err := df.Write(i, p)
		assert.Nil(t, err)
	}

	// read back from the same DataFile
	for i = 0; i < 10; i++ {
		p := getPoint(i)
		p2, err := df.Read(i)

		assert.Nil(t, err)
		assert.True(t, p.Equal(p2))
	}

	// read back from a new DataFile object
	df2, err := OpenExistingDF(fn, ser)
	assert.Nil(t, err)

	for i = 9; i >= 1; i-- {
		p := getPoint(i)
		p2, err := df2.Read(i)

		assert.Nil(t, err)
		assert.True(t, p.Equal(p2))
	}
}

func TestDFMissing(t *testing.T) {
	fn, err := tempFileName()
	assert.Nil(t, err)
	defer os.Remove(fn)

	ser := NewSerializer()

	_, err = OpenExistingDF(fn, ser)
	assert.NotNil(t, err)
}

// tests writing/reading non-sequential points
func TestDFNonseq(t *testing.T) {
	fn, err := tempFileName()
	assert.Nil(t, err)
	defer os.Remove(fn)

	ser := NewSerializer()

	b, _ := ser.Serialize(getPoint(0))

	df, err := OpenNewDF(fn, ser, uint32(len(b)))
	assert.Nil(t, err)
	defer df.Close()

	// write points 0, 2
	df.Write(0, getPoint(0))
	df.Write(2, getPoint(2))

	// point 0 and 2 should work
	for i := 0; i < 3; i += 2 {
		p := getPoint(uint32(i))
		p2, err := df.Read(uint32(i))

		assert.Nil(t, err)
		assert.True(t, p.Equal(p2))
	}

	// point 3+ should fail
	for i := 3; i < 600; i += 100 {
		_, err := df.Read(uint32(i))
		assert.NotNil(t, err)
	}

	// point 1 should fail
	for i := 1; i < 2; i++ {
		_, err := df.Read(uint32(i))
		assert.NotNil(t, err)
	}
}
