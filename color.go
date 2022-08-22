package main

type Color struct {
	r, g, b, a uint8
}

var (
	ColorBlack = Color{0, 0, 0, 255}
	ColorWhite = Color{255, 255, 255, 255}

	ColorRed   = Color{255, 0, 0, 255}
	ColorGreen = Color{0, 255, 0, 255}
	ColorBlue  = Color{255, 0, 255, 255}

	ColorYellow = Color{255, 255, 0, 255}
)
