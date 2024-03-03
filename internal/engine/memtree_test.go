package engine

import (
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
