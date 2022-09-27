// This file contains the color buffer that will be sent to the renderer to
// display. Each color represents a pixel on-screen.
package heretic

import "image/color"

func NewColorBuffer(width, height int) ColorBuffer {
	return ColorBuffer{
		width:  width,
		height: height,
		buf:    make([]color.NRGBA, height*width),
	}
}

type ColorBuffer struct {
	height, width int
	buf           []color.NRGBA
}

func (b ColorBuffer) Set(x, y int, color color.NRGBA) {
	if x > 0 && x < int(b.width) && y > 0 && y < int(b.height) {
		b.buf[(b.width*y)+x] = color
	}
}

// Clear writes over every color in the buffer
func (b ColorBuffer) Clear(color color.NRGBA) {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			b.Set(x, y, color)
		}
	}
}

// Clear writes over every color in the buffer
func (b ColorBuffer) SetBackground(background Background) {
	for y := 0; y < b.height; y++ {
		color := background.At(y, b.height)
		for x := 0; x < b.width; x++ {
			b.Set(x, y, color)
		}
	}
}
