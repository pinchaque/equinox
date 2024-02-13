package equinox

import "math"

func floatEqual(x, y float64) bool {
	const tolerance = 0.00001
	diff := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	if math.IsNaN(diff / mean) {
		return true
	}
	return (diff / mean) < tolerance
}
