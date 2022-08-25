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
func (r Renderer) DrawPixel(x, y int, color Color) {
	if x > 0 && x < int(r.width) && y > 0 && y < int(r.height) {
		r.colorBuffer[(r.width*y)+x] = color
	}
}

func (r Renderer) DrawTexel(x, y int, pointA, pointB, pointC Vec4, u0, v0, u1, v1, u2, v2 float64, texture Texture) {
	pointP := Vec2{float64(x), float64(y)}
	a := pointA.Vec2()
	b := pointB.Vec2()
	c := pointC.Vec2()

	weights := barycentricWeights(a, b, c, pointP)

	alpha := weights.x
	beta := weights.y
	gamma := weights.z

	var interpolatedU, interpolatedV, interpolatedReciprocalW float64

	interpolatedU = (u0/pointA.w)*alpha + (u1/pointB.w)*beta + (u2/pointC.w)*gamma
	interpolatedV = (v0/pointA.w)*alpha + (v1/pointB.w)*beta + (v2/pointC.w)*gamma
	interpolatedReciprocalW = (1/pointA.w)*alpha + (1/pointB.w)*beta + (1/pointC.w)*gamma

	interpolatedU /= interpolatedReciprocalW
	interpolatedV /= interpolatedReciprocalW

	// FIXME: Texture Width is hardcoded
	textureX := int(math.Abs(interpolatedU*float64(texture.width))) % texture.width
	textureY := int(math.Abs(interpolatedV*float64(texture.height))) % texture.height

	r.DrawPixel(x, y, texture.data[(textureY*64)+textureX])
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

// Draw a filled triangle with the flat-top/flat-bottom method.  We split the
// original triangle in two, half flat-bottom and half flat-top
//
//          (x0,y0)
//            / \
//           /   \
//          /     \
//         /       \
//        /         \
//   (x1,y1)------(Mx,My)
//       \_           \
//          \_         \
//             \_       \
//                \_     \
//                   \    \
//                     \_  \
//                        \_\
//                           \
//                         (x2,y2)
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

// Draw a filled a triangle with a flat bottom
//
//        (x0,y0)
//          / \
//         /   \
//        /     \
//       /       \
//      /         \
//  (x1,y1)------(x2,y2)
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

// Draw a filled a triangle with a flat top
//
//  (x0,y0)------(x1,y1)
//      \         /
//       \       /
//        \     /
//         \   /
//          \ /
//        (x2,y2)
//
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

func (r Renderer) DrawTexturedTriangle(
	x0, y0 int, z0, w0 float64, u0, v0 float64,
	x1, y1 int, z1, w1 float64, u1, v1 float64,
	x2, y2 int, z2, w2 float64, u2, v2 float64,
	texture Texture,
) {

	if y0 > y1 {
		y0, y1 = y1, y0
		x0, x1 = x1, x0
		z0, z1 = z1, z0
		w0, w1 = w1, w0
		u0, u1 = u1, u0
		v0, v1 = v1, v0
	}

	if y1 > y2 {
		y1, y2 = y2, y1
		x1, x2 = x2, x1
		z1, z2 = z2, z1
		w1, w2 = w2, w1
		u1, u2 = u2, u1
		v1, v2 = v2, v1

	}

	if y0 > y1 {
		y0, y1 = y1, y0
		x0, x1 = x1, x0
		z0, z1 = z1, z0
		w0, w1 = w1, w0
		u0, u1 = u1, u0
		v0, v1 = v1, v0
	}

	pointA := Vec4{float64(x0), float64(y0), z0, w0}
	pointB := Vec4{float64(x1), float64(y1), z1, w1}
	pointC := Vec4{float64(x2), float64(y2), z2, w2}

	//
	// Top part of triangle
	//
	var invSlope1, invSlope2 float64

	if y1-y0 != 0 {
		invSlope1 = float64(x1-x0) / math.Abs(float64(y1-y0))
	}
	if y2-y0 != 0 {
		invSlope2 = float64(x2-x0) / math.Abs(float64(y2-y0))
	}

	if y1-y0 != 0 {
		for y := y0; y <= y1; y++ {
			xStart := int(float64(x1) + float64(y-y1)*invSlope1)
			xEnd := int(float64(x0) + float64(y-y0)*invSlope2)

			// If xStart is to the left of xEnd
			if xStart > xEnd {
				xStart, xEnd = xEnd, xStart
			}

			for x := xStart; x < xEnd; x++ {
				r.DrawTexel(x, y, pointA, pointB, pointC, u0, v0, u1, v1, u2, v2, texture)
			}
		}
	}

	//
	// Bottom part of triangle
	//
	invSlope1, invSlope2 = 0.0, 0.0

	if y2-y1 != 0 {
		invSlope1 = float64(x2-x1) / math.Abs(float64(y2-y1))
	}
	if y2-y0 != 0 {
		invSlope2 = float64(x2-x0) / math.Abs(float64(y2-y0))
	}

	if y2-y1 != 0 {
		for y := y1; y <= y2; y++ {
			xStart := int(float64(x1) + float64(y-y1)*invSlope1)
			xEnd := int(float64(x0) + float64(y-y0)*invSlope2)

			// If xStart is to the left of xEnd
			if xStart > xEnd {
				xStart, xEnd = xEnd, xStart
			}

			for x := xStart; x < xEnd; x++ {
				r.DrawTexel(x, y, pointA, pointB, pointC, u0, v0, u1, v1, u2, v2, texture)
			}
		}
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

func barycentricWeights(a, b, c, p Vec2) Vec3 {
	ab := b.Sub(a)
	bc := c.Sub(b)
	ac := c.Sub(a)
	ap := p.Sub(a)
	bp := p.Sub(b)

	// Calcualte the area of the full triangle ABC using cross product (area of parallelogram)
	triangle_area := (ab.x*ac.y - ab.y*ac.x)

	// Weight alpha is the area of subtriangle BCP divided by the area of the full triangle ABC
	alpha := (bc.x*bp.y - bp.x*bc.y) / triangle_area

	// Weight beta is the area of subtriangle ACP divided by the area of the full triangle ABC
	beta := (ap.x*ac.y - ac.x*ap.y) / triangle_area

	// Weight gamma is easily found since barycentric cooordinates always add up to 1
	gamma := 1 - alpha - beta

	weights := Vec3{alpha, beta, gamma}
	return weights
}
