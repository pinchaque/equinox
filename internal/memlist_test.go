package equinox

import (
	"testing"
)

func TestMemList(t *testing.T) {
	ml := NewMemList()
	testPointIO(t, ml, 10, 5)
}
