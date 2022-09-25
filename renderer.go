package heretic

import "math"

func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		width:       width,
		height:      height,
		colorBuffer: make([]Color, width*height),
		zBuffer:     make([]float64, width*height),
	}
}

type Renderer struct {
	width, height int
	colorBuffer   []Color
	zBuffer       []float64
}

// DrawPixel draws a single colored pixel at the specified coordinates.
func (r Renderer) DrawPixel(x, y int, color Color) {
	if x > 0 && x < int(r.width) && y > 0 && y < int(r.height) {
		r.colorBuffer[(r.width*y)+x] = color
	}
}

func (r Renderer) DrawTexel(x, y int, a, b, c Vec4, auv, buv, cuv Tex, lightIntensity float64, texture Texture, palette *Palette) {
	pointP := Vec2{float64(x), float64(y)}

	weights := barycentricWeights(a.Vec2(), b.Vec2(), c.Vec2(), pointP)

	alpha := weights.X
	beta := weights.Y
	gamma := weights.Z

	var interpolatedU, interpolatedV, interpolatedReciprocalW float64

	interpolatedU = (auv.U/a.W)*alpha + (buv.U/b.W)*beta + (cuv.U/c.W)*gamma
	interpolatedV = (auv.V/a.W)*alpha + (buv.V/b.W)*beta + (cuv.V/c.W)*gamma

	// FIXME: move this calculation out of the function as it only needs to
	// be calcualted once per triangle.
	interpolatedReciprocalW = (1/a.W)*alpha + (1/b.W)*beta + (1/c.W)*gamma

	interpolatedU /= interpolatedReciprocalW
	interpolatedV /= interpolatedReciprocalW

	// FIXME: Texture Width is hardcoded
	textureX := int(math.Abs(interpolatedU*float64(texture.width))) % texture.width
	textureY := int(math.Abs(interpolatedV*float64(texture.height))) % texture.height

	// Adjust 1/w so the pixels that are closer to the cam have smaller values
	interpolatedReciprocalW = 1.0 - interpolatedReciprocalW

	// Only draw pixel if depth value is less than one previously stored in zbuffer.
	if interpolatedReciprocalW < r.ZBufferAt(x, y) {
		textureColor := texture.data[(textureY*texture.width)+textureX]
		// If there is a palette, the current color components will
		// represent the index into the palette.
		if palette != nil {
			textureColor = palette[textureColor.R]
		}

		textureColorWithLight := textureColor
		// Disabling this until we get proper lighting
		// textureWithLightColor := applyLightIntensity(textureColor, lightIntensity)

		// This handels transparent colors when there is a palette (FFT).
		if textureColor.IsTransparent() && palette != nil {
			return
		}
		r.DrawPixel(x, y, textureColorWithLight)
		r.ZBufferSet(x, y, interpolatedReciprocalW)
	}
}

func (r Renderer) DrawTrianglePixel(x, y int, a, b, c Vec4, color Color) {
	pointP := Vec2{float64(x), float64(y)}

	weights := barycentricWeights(a.Vec2(), b.Vec2(), c.Vec2(), pointP)

	alpha := weights.X
	beta := weights.Y
	gamma := weights.Z

	// FIXME: move this calculation out of the function as it only needs to
	// be calcualted once per triangle.
	interpolatedReciprocalW := (1/a.W)*alpha + (1/b.W)*beta + (1/c.W)*gamma

	// Adjust 1/w so the pixels that are closer to the cam have smaller values
	interpolatedReciprocalW = 1.0 - interpolatedReciprocalW

	// Only draw pixel if depth value is less than one previously stored in zbuffer.
	if interpolatedReciprocalW < r.ZBufferAt(x, y) {
		r.DrawPixel(x, y, color)
		r.ZBufferSet(x, y, interpolatedReciprocalW)
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
func (r Renderer) DrawFilledTriangle(
	x0, y0 int, z0 float64, w0 float64,
	x1, y1 int, z1 float64, w1 float64,
	x2, y2 int, z2 float64, w2 float64,
	color Color,
) {
	if y0 > y1 {
		y0, y1 = y1, y0
		x0, x1 = x1, x0
		z0, z1 = z1, z0
		w0, w1 = w1, w0
	}

	if y1 > y2 {
		y1, y2 = y2, y1
		x1, x2 = x2, x1
		z1, z2 = z2, z1
		w1, w2 = w2, w1

	}

	if y0 > y1 {
		y0, y1 = y1, y0
		x0, x1 = x1, x0
		z0, z1 = z1, z0
		w0, w1 = w1, w0
	}

	a := Vec4{float64(x0), float64(y0), z0, w0}
	b := Vec4{float64(x1), float64(y1), z1, w1}
	c := Vec4{float64(x2), float64(y2), z2, w2}

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
				r.DrawTrianglePixel(x, y, a, b, c, color)
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
				r.DrawTrianglePixel(x, y, a, b, c, color)
			}
		}
	}
}

func (r Renderer) DrawTexturedTriangle(
	x0, y0 int, z0, w0 float64, at Tex,
	x1, y1 int, z1, w1 float64, bt Tex,
	x2, y2 int, z2, w2 float64, ct Tex,
	lightIntensity float64,
	texture Texture,
	palette *Palette,
) {

	if y0 > y1 {
		y0, y1 = y1, y0
		x0, x1 = x1, x0
		z0, z1 = z1, z0
		w0, w1 = w1, w0
		at.U, bt.U = bt.U, at.U
		at.V, bt.V = bt.V, at.V
	}

	if y1 > y2 {
		y1, y2 = y2, y1
		x1, x2 = x2, x1
		z1, z2 = z2, z1
		w1, w2 = w2, w1
		bt.U, ct.U = ct.U, bt.U
		bt.V, ct.V = ct.V, bt.V

	}

	if y0 > y1 {
		y0, y1 = y1, y0
		x0, x1 = x1, x0
		z0, z1 = z1, z0
		w0, w1 = w1, w0
		at.U, bt.U = bt.U, at.U
		at.V, bt.V = bt.V, at.V
	}

	// FIXME: Flip the texture coordinates (handle this on import?)
	at.V = 1 - at.V
	bt.V = 1 - bt.V
	ct.V = 1 - ct.V

	a := Vec4{float64(x0), float64(y0), z0, w0}
	b := Vec4{float64(x1), float64(y1), z1, w1}
	c := Vec4{float64(x2), float64(y2), z2, w2}

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
				r.DrawTexel(x, y, a, b, c, at, bt, ct, lightIntensity, texture, palette)
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
				r.DrawTexel(x, y, a, b, c, at, bt, ct, lightIntensity, texture, palette)
			}
		}
	}
}

// Clear writes over every color in the buffer
func (r Renderer) ColorBufferColor(color Color) {
	for x := 0; x < r.width; x++ {
		for y := 0; y < r.height; y++ {
			r.DrawPixel(x, y, color)
		}
	}
}

// Clear writes over every color in the buffer
func (r Renderer) ColorBufferBackground(bg Background) {
	for y := 0; y < r.height; y++ {
		color := bg.At(y, r.height)
		for x := 0; x < r.width; x++ {
			r.DrawPixel(x, y, color)
		}
	}
}

// Clear writes over every color in the buffer
func (r Renderer) ZBufferClear() {
	for x := 0; x < r.width; x++ {
		for y := 0; y < r.height; y++ {
			r.zBuffer[(y*r.width)+x] = 1.0
		}
	}
}

func (r Renderer) ZBufferAt(x, y int) float64 {
	if x < 0 || x >= r.width || y < 0 || y >= r.height {
		return 1.0
	}
	return r.zBuffer[(y*r.width)+x]
}

func (r Renderer) ZBufferSet(x, y int, v float64) {
	if x < 0 || x >= r.width || y < 0 || y >= r.height {
		return
	}
	r.zBuffer[(y*r.width)+x] = v
}

func barycentricWeights(a, b, c, p Vec2) Vec3 {
	ab := b.Sub(a)
	bc := c.Sub(b)
	ac := c.Sub(a)
	ap := p.Sub(a)
	bp := p.Sub(b)

	// Calcualte the area of the full triangle ABC using cross product (area of parallelogram)
	triangleArea := (ab.X*ac.Y - ab.Y*ac.X)

	// Weight alpha is the area of subtriangle BCP divided by the area of the full triangle ABC
	alpha := (bc.X*bp.Y - bp.X*bc.Y) / triangleArea

	// Weight beta is the area of subtriangle ACP divided by the area of the full triangle ABC
	beta := (ap.X*ac.Y - ac.X*ap.Y) / triangleArea

	// Weight gamma is easily found since barycentric cooordinates always add up to 1
	gamma := 1 - alpha - beta

	weights := Vec3{alpha, beta, gamma}
	return weights
}
