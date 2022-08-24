package main

// Face represents a triangle before rasterization.
type Face struct {
	points [3]Vec3
	color  Color
}

// Triangle represents a triangle after rasterization.
type Triangle struct {
	points       [3]Vec2
	color        Color
	averageDepth float64
}
