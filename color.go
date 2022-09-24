package heretic

type Color struct {
	R, G, B, A uint8
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
