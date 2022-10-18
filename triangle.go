package heretic

import (
	"image/color"
)

type Triangle struct {
	// Points represents a vertices before rasterization.
	Points []Vec3

	// Projected represents a vertices after rasterization.
	Projected []Vec4

	Normals   []Vec3
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
	a := t.Points[0]
	b := t.Points[1]
	c := t.Points[2]
	vectorAB := b.Sub(a).Normalize()
	vectorAC := c.Sub(a).Normalize()
	normal := vectorAB.Cross(vectorAC).Normalize()
	return normal
}

// FIXME: This is a clunky check. We should have triangle.Texcoords be nil
// instead of this. But the fft importer currently returns 3 empty texcoords
// regardless of whats in them. You cant just return nil until we add additional
// checks into the clipping functionality because it indexes texcoords.
func (t Triangle) HasTexture() bool {
	for _, tc := range t.Texcoords {
		if !tc.IsEmpty() {
			return true
		}
	}
	return false
}
