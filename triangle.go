package heretic

import (
	"image/color"
)

type Triangle struct {
	// Points represents a vertices before rasterization.
	Points []Vec3

	// Projected represents a vertices after rasterization.
	Projected []Vec4

	Texcoords []Tex

	// Palette represents the 16-color Palette to use during rendering a
	// polygon.  This is due to FFT texture storage. The raw texture pixel
	// value is an index for a palettes. Each map has 16 palettes of 16
	// colors each. Each polygon references on of the 16 palettes to use.
	// Eventually Renderer.DrawTexel() function uses uses the pallet.
	Palette Palette

	// Color is used when there is no texture or when there is a texture,
	// but the polygon has no palette.
	Color color.NRGBA

	LightIntensity float64
}

// Normal calculates and returns the face normal for the triangle.
// This is a left handed system.
func (t Triangle) Normal() Vec3 {
	a := t.Projected[0].Vec3()
	b := t.Projected[1].Vec3()
	c := t.Projected[2].Vec3()
	vectorAB := b.Sub(a).Normalize()
	vectorAC := c.Sub(a).Normalize()
	normal := vectorAB.Cross(vectorAC).Normalize()
	return normal
}
