package main

import "math"

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
func (r *Renderer) DrawPixel(x, y int, color Color) {
	if x > 0 && x < int(r.width) && y > 0 && y < int(r.height) {
		r.colorBuffer[(r.width*y)+x] = color
	}
}

// DrawLine draws a solid line using the DDA algorithm.
func (r *Renderer) DrawLine(a, b Vec2, color Color) {
	deltaX := b.x - a.x
	deltaY := b.y - a.y

	var sideLength int
	if math.Abs(deltaX) >= math.Abs(deltaY) {
		sideLength = int(math.Abs(deltaX))
	} else {
		sideLength = int(math.Abs(deltaY))
	}

	incX := deltaX / float64(sideLength)
	incY := deltaY / float64(sideLength)

	currentX := a.x
	currentY := a.y

	for i := 0; i < sideLength; i++ {
		r.DrawPixel(int(math.Round(currentX)), int(math.Round(currentY)), color)
		currentX += incX
		currentY += incY
	}
}

// DrawGrid draws a dotted grid across entire buffer.
func (r *Renderer) DrawGrid(color Color) {
	for y := 0; y < r.height; y += 10 {
		for x := 0; x < r.width; x += 10 {
			r.DrawPixel(x, y, color)
		}
	}
}

// DrawGrid draws a rectangle to the buffer.
func (r *Renderer) DrawRectangle(x, y, width, height int, color Color) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			currentX := x + i
			currentY := y + j
			r.DrawPixel(currentX, currentY, color)
		}
	}
}

func (r *Renderer) DrawTriangle(t Triangle, color Color) {
	r.DrawLine(t.points[0], t.points[1], color)
	r.DrawLine(t.points[1], t.points[2], color)
	r.DrawLine(t.points[2], t.points[0], color)
}

// Clear writes over every color in the buffer
func (r *Renderer) Clear(color Color) {
	for x := 0; x < r.width; x++ {
		for y := 0; y < r.height; y++ {
			r.DrawPixel(x, y, color)
		}
	}
}
