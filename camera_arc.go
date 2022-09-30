// This file contains a simple Arcball Camera that will rotate around the target.
//
// FIXME: when rotating to the left or right, then dragging up, we should take
//        into consideration the side we are looking at and make that side go up
//        instead of the front going up or down regardless of viewing angle.
// FIXME: Add mouse wheel movement
// FIXME: Add panning
//
package heretic

import "math"

type ArcCamera struct {
	eye     Vec3
	front   Vec3
	up      Vec3
	worldUp Vec3

	viewMatrix Matrix

	height, width float64

	rightButtonHeld bool
	leftButtonHeld  bool
}

func NewArcCamera(eye, front, up Vec3, width, height int) *ArcCamera {
	c := ArcCamera{
		eye:     eye,
		front:   front,
		up:      up,
		worldUp: Vec3{0, 1, 0},
		height:  float64(height),
		width:   float64(width),
	}
	c.updateViewMatrix()
	return &c
}

func (c *ArcCamera) updateViewMatrix() { c.viewMatrix = LookAt(c.eye, c.front, c.up) }

func (c *ArcCamera) ViewMatrix() Matrix { return c.viewMatrix }

func (c *ArcCamera) ProcessKeyboardInput(state []uint8, delta float64) {}
func (c *ArcCamera) ProcessMouseWheel(y float64, delta float64)        {}

func (c *ArcCamera) ProcessMouseMovement(xrel, yrel, delta float64) {
	if c.leftButtonHeld {
		// Calculate the amount of rotation given the mouse movement.
		var deltaAngleX float64 = (2.0 * math.Pi / c.width)
		var deltaAngleY float64 = (math.Pi / c.height)
		var xAngle float64 = float64(xrel) * deltaAngleX
		var yAngle float64 = float64(yrel) * deltaAngleY

		cosAngle := float64(c.front.Dot(c.worldUp))
		if cosAngle*sgn(deltaAngleY) > 0.99 {
			yAngle = 0.0
		}

		position := c.eye.Vec4()
		pivot := c.front.Vec4()

		// step 2: Rotate the camera around the pivot point on the first axis.
		rotationMatrixX := MatrixMakeRotX(yAngle)
		position = rotationMatrixX.MulVec4(position.Sub(pivot)).Add(pivot)

		rotationMatrixY := MatrixMakeRotY(xAngle)
		position = rotationMatrixY.MulVec4(position.Sub(pivot)).Add(pivot)
		c.eye = position.Vec3()

		c.updateViewMatrix()
	}
}

func (c *ArcCamera) ProcessMouseButton(button MouseButton, pressed bool) {
	if button == MouseButtonLeft {
		c.leftButtonHeld = pressed
	}
	if button == MouseButtonRight {
		c.rightButtonHeld = pressed
	}
}
