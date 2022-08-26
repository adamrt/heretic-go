package main

import "math"

func abs(i int) int {
	return int(math.Abs(float64(i)))
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}
