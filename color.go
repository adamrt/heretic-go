package heretic

type Color struct {
	R, G, B, A uint8
}

func (c Color) IsTransparent() bool {
	return c.R+c.G+c.B+c.A == 0
}

// Palette represents the 16-color palette to use during rendering a polygon.
type Palette [16]Color

var (
	ColorBlack = Color{0, 0, 0, 255}
	ColorWhite = Color{255, 255, 255, 255}
	ColorGrey  = Color{0x55, 0x55, 0x55, 255}

	ColorRed   = Color{255, 0, 0, 255}
	ColorGreen = Color{0, 255, 0, 255}
	ColorBlue  = Color{0, 0, 255, 255}

	ColorYellow  = Color{255, 255, 0, 255}
	ColorCyan    = Color{0, 255, 255, 255}
	ColorMagenta = Color{255, 0, 255, 255}
)

type Background struct {
	Top    Color
	Bottom Color
}

// These are vertical gradients so we don't care about x.
func (bg Background) At(y int, height int) Color {
	d := float64(y) / float64(height)
	r := float64(bg.Bottom.R) + d*(float64(bg.Top.R)-float64(bg.Bottom.R))
	g := float64(bg.Bottom.G) + d*(float64(bg.Top.G)-float64(bg.Bottom.G))
	b := float64(bg.Bottom.B) + d*(float64(bg.Top.B)-float64(bg.Bottom.B))
	return Color{uint8(r), uint8(g), uint8(b), 255}
}
