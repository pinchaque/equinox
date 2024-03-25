package file

import (
	"equinox/internal/core"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	const exptime string = "2024-01-10T23:01:02.000000Z"
	const fmtstr string = "2006-01-02T15:04:05.000000Z"

	assert.Equal(t, exptime, ts.Format(fmtstr))

	p := core.NewPoint(ts)
	p.Attrs["color"] = "red"
	p.Attrs["shape"] = "square"
	p.Vals["area"] = 43.1
	p.Vals["temp"] = 21.1

	s := NewSerializer()

	data, err := s.Serialize(p)

	assert.Nil(t, err)

	// expected size: 16 + 12*num_values + 8*num_attrs = 16 + 24 + 16 = 56
	assert.Equal(t, 56, len(data))

	p2, err := s.Deserialize(data)

	assert.Nil(t, err)
	assert.True(t, p2.Equal(p))

}
