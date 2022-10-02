package heretic

import "image/color"

func NewFrameBuffer(width, height int) FrameBuffer {
	return FrameBuffer{
		width:  width,
		height: height,
		zbuf:   make([]float64, height*width),
		cbuf:   make([]color.NRGBA, height*width),
	}
}

type FrameBuffer struct {
	height, width int
	zbuf          []float64
	cbuf          []color.NRGBA
}

// Clear writes over every color in the buffer
func (b FrameBuffer) Clear(color color.NRGBA) {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			b.SetColor(x, y, color)
		}
	}
}

func (b FrameBuffer) SetColor(x, y int, color color.NRGBA) {
	if x > 0 && x < int(b.width) && y > 0 && y < int(b.height) {
		b.cbuf[(b.width*y)+x] = color
	}
}

// Clear writes over every color in the buffer
func (b FrameBuffer) SetBackground(background Background) {
	for y := 0; y < b.height; y++ {
		color := background.At(y, b.height)
		for x := 0; x < b.width; x++ {
			b.SetColor(x, y, color)
		}
	}
}

func (b FrameBuffer) ClearDepth() {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			b.SetDepth(x, y, 1.0)
		}
	}
}

func (b FrameBuffer) DepthAt(x, y int) float64 {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return 1.0
	}
	return b.zbuf[(y*b.width)+x]
}

func (b FrameBuffer) SetDepth(x, y int, v float64) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return
	}
	b.zbuf[(y*b.width)+x] = v
}
