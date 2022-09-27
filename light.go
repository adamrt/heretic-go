// This file contains our different lighting types.
//
// Our lighting system is virtually non-existent currently, and is not being
// used. Coming soon!
package heretic

import (
	"image/color"
	"math"
)

type AmbientLight struct {
	Color color.NRGBA
}

type DirectionalLight struct {
	Direction Vec3
	Position  Vec3
	Color     color.NRGBA
}

func applyLightIntensity(orig color.NRGBA, factor float64) color.NRGBA {
	// Clamp from 0.0 to 1.0
	factor = math.Max(0, math.Min(factor, 1.0))

	return color.NRGBA{
		R: uint8(float64(orig.R) * factor),
		G: uint8(float64(orig.G) * factor),
		B: uint8(float64(orig.B) * factor),
		A: orig.A,
	}
}
