package heretic

import "image/color"

// Palette represents the 16-color palette to use during rendering a polygon.
type Palette [16]color.NRGBA

var (
	ColorBlack = color.NRGBA{0, 0, 0, 255}
	ColorWhite = color.NRGBA{255, 255, 255, 255}
	ColorGrey  = color.NRGBA{0x55, 0x55, 0x55, 255}

	ColorRed   = color.NRGBA{255, 0, 0, 255}
	ColorGreen = color.NRGBA{0, 255, 0, 255}
	ColorBlue  = color.NRGBA{0, 0, 255, 255}

	ColorYellow  = color.NRGBA{255, 255, 0, 255}
	ColorCyan    = color.NRGBA{0, 255, 255, 255}
	ColorMagenta = color.NRGBA{255, 0, 255, 255}
)

type Background struct {
	Top    color.NRGBA
	Bottom color.NRGBA
}

// These are vertical gradients so we don't care about x.  The colors need to be
// float64s before subtraction so there isn't uint8 overflow.
func (bg Background) At(y int, height int) color.NRGBA {
	d := float64(y) / float64(height)
	r := float64(bg.Bottom.R) + d*(float64(bg.Top.R)-float64(bg.Bottom.R))
	g := float64(bg.Bottom.G) + d*(float64(bg.Top.G)-float64(bg.Bottom.G))
	b := float64(bg.Bottom.B) + d*(float64(bg.Top.B)-float64(bg.Bottom.B))
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}
