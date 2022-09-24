package heretic

import "math"

type AmbientLight struct {
	Color Color
}

type DirectionalLight struct {
	Direction Vec3
	Position  Vec3
	Color     Color
}

func applyLightIntensity(orig Color, factor float64) Color {
	// Clamp from 0.0 to 1.0
	factor = math.Max(0, math.Min(factor, 1.0))

	return Color{
		R: uint8(float64(orig.R) * factor),
		G: uint8(float64(orig.G) * factor),
		B: uint8(float64(orig.B) * factor),
		A: orig.A,
	}
}
