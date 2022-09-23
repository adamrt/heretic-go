package heretic

import "math"

type Light struct {
	direction Vec3
}

func applyLightIntensity(orig Color, factor float64) Color {
	// Clamp from 0.0 to 1.0
	factor = math.Max(0, math.Min(factor, 1.0))

	return Color{
		r: uint8(float64(orig.r) * factor),
		g: uint8(float64(orig.g) * factor),
		b: uint8(float64(orig.b) * factor),
		a: orig.a,
	}
}
