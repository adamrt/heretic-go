package main

import "math"

// FOV is just matching the width of the screen currently.  This is used to
// scale up the positions of the points from
const FOV float64 = WindowWidth

func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		width:       width,
		height:      height,
		colorBuffer: make([]Color, width*height),
	}
}

type Renderer struct {
	width, height int
	colorBuffer   []Color
}

// DrawPixel draws a single colored pixel at the specified coordinates.
func (r Renderer) DrawPixel(x, y int, color Color) {
	if x > 0 && x < int(r.width) && y > 0 && y < int(r.height) {
		r.colorBuffer[(r.width*y)+x] = color
	}
}

// DrawLine draws a solid line using the DDA algorithm.
func (r Renderer) DrawLine(x0, y0, x1, y1 int, color Color) {
	deltaX := x1 - x0
	deltaY := y1 - y0

	var longestSideLength int
	if abs(deltaX) >= abs(deltaY) {
		longestSideLength = abs(deltaX)
	} else {
		longestSideLength = abs(deltaY)
	}

	incX := float64(deltaX) / float64(longestSideLength)
	incY := float64(deltaY) / float64(longestSideLength)

	currentX := float64(x0)
	currentY := float64(y0)

	for i := 0; i <= longestSideLength; i++ {
		r.DrawPixel(int(math.Round(currentX)), int(math.Round(currentY)), color)
		currentX += incX
		currentY += incY
	}
}

// DrawGrid draws a dotted grid across entire buffer.
func (r Renderer) DrawGrid(color Color) {
	for y := 0; y < r.height; y += 10 {
		for x := 0; x < r.width; x += 10 {
			r.DrawPixel(x, y, color)
		}
	}
}

// DrawGrid draws a rectangle to the buffer.
func (r Renderer) DrawRectangle(x, y, width, height int, color Color) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			currentX := x + i
			currentY := y + j
			r.DrawPixel(currentX, currentY, color)
		}
	}
}

func (r Renderer) DrawTriangle(x0, y0, x1, y1, x2, y2 int, color Color) {
	r.DrawLine(x0, y0, x1, y1, color)
	r.DrawLine(x1, y1, x2, y2, color)
	r.DrawLine(x2, y2, x0, y0, color)
}

func (r Renderer) DrawFilledTriangle(x0, y0, x1, y1, x2, y2 int, color Color) {
	if y0 > y1 {
		y0, y1 = y1, y0
		x0, x1 = x1, x0
	}

	if y1 > y2 {
		y1, y2 = y2, y1
		x1, x2 = x2, x1

	}

	if y0 > y1 {
		y0, y1 = y1, y0
		x0, x1 = x1, x0
	}

	if y1 == y2 {
		r.fillFlatBottom(x0, y0, x1, y1, x2, y2, color)
	} else if y0 == y1 {
		r.fillFlatTop(x0, y0, x1, y1, x2, y2, color)
	} else {
		my := y1
		mx := ((x2 - x0) * (y1 - y0) / (y2 - y0)) + x0

		r.fillFlatBottom(x0, y0, x1, y1, mx, my, color)
		r.fillFlatTop(x1, y1, mx, my, x2, y2, color)
	}
}

func (r Renderer) fillFlatBottom(x0, y0, x1, y1, x2, y2 int, color Color) {
	invSlope1 := float64(x1-x0) / float64(y1-y0)
	invSlope2 := float64(x2-x0) / float64(y2-y0)

	xStart := float64(x0)
	xEnd := float64(x0)

	// Loop scanlines bottom to top
	for y := y0; y <= y2; y++ {
		r.DrawLine(int(xStart), y, int(xEnd), y, color)
		xStart += invSlope1
		xEnd += invSlope2
	}
}

func (r Renderer) fillFlatTop(x0, y0, x1, y1, x2, y2 int, color Color) {
	invSlope1 := float64(x2-x0) / float64(y2-y0)
	invSlope2 := float64(x2-x1) / float64(y2-y1)

	xStart := float64(x2)
	xEnd := float64(x2)

	// Loop scanlines bottom to top
	for y := y2; y >= y0; y-- {
		r.DrawLine(int(xStart), y, int(xEnd), y, color)
		xStart -= invSlope1
		xEnd -= invSlope2
	}
}

// Clear writes over every color in the buffer
func (r Renderer) Clear(color Color) {
	for x := 0; x < r.width; x++ {
		for y := 0; y < r.height; y++ {
			r.DrawPixel(x, y, color)
		}
	}
}
