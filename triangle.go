package heretic

import "image/color"

// Triangle represents a triangle after rasterization.
type Triangle struct {
	points    [3]Vec4
	texcoords [3]Tex

	// Palette represents the 16-color palette to use during rendering a
	// polygon.  This is due to FFT texture storage. The raw texture pixel
	// value is an index for a palettes. Each map has 16 palettes of 16
	// colors each. Each polygon references on of the 16 palettes to use. It
	// is just passed from Face to Triangle and not used until
	// Renderer.DrawTexel() function.
	palette *Palette

	// Color is used when there is no texture or when there is a texture,
	// but the polygon has no palette.
	color color.NRGBA

	lightIntensity float64
}

func (t Triangle) Normal() Vec3 {
	a := t.points[0].Vec3()
	b := t.points[1].Vec3()
	c := t.points[2].Vec3()
	vectorAB := b.Sub(a).Normalize()
	vectorAC := c.Sub(a).Normalize()
	normal := vectorAB.Cross(vectorAC).Normalize() // Left handed system
	return normal
}

func (t Triangle) HasTexture() bool {
	for _, tc := range t.texcoords {
		if !tc.IsEmpty() {
			return true
		}
	}
	return false
}
