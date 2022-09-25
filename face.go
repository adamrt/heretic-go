package heretic

// Face represents a triangle before rasterization.

func NewFace(points [3]Vec3, texcoords [3]Tex, palette *Palette, color Color) Face {
	return Face{points, texcoords, palette, color}
}

type Face struct {
	points    [3]Vec3
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
	color Color
}

type Tex struct {
	U, V float64
}

func (t Tex) IsEmpty() bool {
	return t.U == 0.0 && t.V == 0.0
}
