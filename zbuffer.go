// This file contains the zbuffer that will be used to determine what pixels are
// in front of the other during rendering. It's to overcome the shortcomings of
// the painters algorithm.
package heretic

func NewZBuffer(width, height int) ZBuffer {
	return ZBuffer{
		width:  width,
		height: height,
		buf:    make([]float64, height*width),
	}
}

type ZBuffer struct {
	height, width int
	buf           []float64
}

func (z ZBuffer) At(x, y int) float64 {
	if x < 0 || x >= z.width || y < 0 || y >= z.height {
		return 1.0
	}
	return z.buf[(y*z.width)+x]
}

func (z ZBuffer) Set(x, y int, v float64) {
	if x < 0 || x >= z.width || y < 0 || y >= z.height {
		return
	}
	z.buf[(y*z.width)+x] = v
}

func (z ZBuffer) Clear() {
	for x := 0; x < z.width; x++ {
		for y := 0; y < z.height; y++ {
			z.Set(x, y, 1.0)
		}
	}
}
