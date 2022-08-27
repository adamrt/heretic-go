package main

// Triangle represents a triangle after rasterization.
type Triangle struct {
	points    [3]Vec4
	texcoords [3]Tex
	color     Color

	lightIntensity float64
}
