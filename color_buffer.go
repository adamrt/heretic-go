package heretic

func NewColorBuffer(width, height int) ColorBuffer {
	return ColorBuffer{
		width:  width,
		height: height,
		buf:    make([]Color, height*width),
	}
}

type ColorBuffer struct {
	height, width int
	buf           []Color
}

func (b ColorBuffer) Set(x, y int, color Color) {
	if x > 0 && x < int(b.width) && y > 0 && y < int(b.height) {
		b.buf[(b.width*y)+x] = color
	}
}

// Clear writes over every color in the buffer
func (b ColorBuffer) Clear(color Color) {
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
