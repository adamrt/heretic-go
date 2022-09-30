package heretic

import "math"

func abs(i int) int {
	return int(math.Abs(float64(i)))
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// Normalize takes the min/max of face coordinates in a mesh and then normalizes
// them to the 0.0-1.0 space.
func normalize(x, min, max float64) float64 {
	return (x - min) / (max - min)
}

func sgn(a float64) float64 {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}
