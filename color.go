package main

type Color struct {
	r, g, b, a uint8
}

var (
	ColorBlack = Color{0, 0, 0, 255}
	ColorWhite = Color{255, 255, 255, 255}
	ColorGrey  = Color{0x55, 0x55, 0x55, 255}

	ColorRed   = Color{255, 0, 0, 255}
	ColorGreen = Color{0, 255, 0, 255}
	ColorBlue  = Color{255, 0, 255, 255}

	ColorYellow = Color{255, 255, 0, 255}
)

func NewColorBuffer(width, height int) ColorBuffer {
	return ColorBuffer{
		width:  width,
		height: height,
		buf:    make([]Color, width*height),
	}
}

type ColorBuffer struct {
	width, height int
	buf           []Color
}

func (b ColorBuffer) Set(x, y int, color Color) {
	b.buf[(b.width*y)+x] = color
}

func (b ColorBuffer) Clear(color Color) {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			b.Set(x, y, color)
		}
	}
}
