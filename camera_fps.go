// This file contains a simple fps camera.
package heretic

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Camera interface {
	ProcessMouseMovement(xrel, yrel, delta float64)
	ProcessMouseWheel(y float64, delta float64)
	ProcessMouseButton(button MouseButton, pressed bool)
	ProcessKeyboardInput(state []uint8, delta float64)
	ViewMatrix() Matrix
}

type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonMiddle
	MouseButtonRight
)

type FPSCamera struct {
	eye   Vec3
	front Vec3
	up    Vec3

	yaw   float64
	pitch float64

	speed float64

	// We store this so we don't have to calculate it per frame per vertex
	// unless it moves.
	viewMatrix Matrix

	rightButtonHeld bool
	leftButtonHeld  bool
}

func NewFPSCamera(eye, front, up Vec3) *FPSCamera {
	c := FPSCamera{
		front: front,
		eye:   eye,
		up:    up,
		speed: 2.0,
	}
	c.updateViewMatrix()
	return &c
}

func (c *FPSCamera) ViewMatrix() Matrix { return c.viewMatrix }

func (c *FPSCamera) updateViewMatrix() {
	c.viewMatrix = LookAt(c.eye, c.eye.Add(c.front), c.up)
}

func (c *FPSCamera) ProcessMouseMovement(xrel, yrel, delta float64) {
	sensitiviy := 0.3
	if c.leftButtonHeld {
		c.yaw += xrel * sensitiviy * delta
		c.pitch += yrel * sensitiviy * delta

		if c.pitch > 1.5 {
			c.pitch = 1.5
		}
		if c.pitch < -1.5 {
			c.pitch = -1.5
		}

		target := Vec3{0, 0, 1}
		cameraRotation := MatrixIdentity().Mul(MatrixMakeRotY(c.yaw)).Mul(MatrixMakeRotX(c.pitch))
		direction := cameraRotation.MulVec4(target.Vec4()).Vec3()
		c.front = direction
		c.updateViewMatrix()
	} else if c.rightButtonHeld {
		// X
		c.eye = c.eye.Add(c.right().Mul(xrel * delta * sensitiviy))

		// Y
		c.eye = c.eye.Add(c.up.Mul(yrel * delta * sensitiviy))
		c.updateViewMatrix()
	}
}

func (c *FPSCamera) ProcessMouseWheel(y float64, delta float64) {
	velocity := c.front.Mul(c.speed * y * delta)
	c.eye = c.eye.Add(velocity)
	c.updateViewMatrix()
}

func (c *FPSCamera) ProcessMouseButton(button MouseButton, pressed bool) {
	if button == MouseButtonLeft {
		c.leftButtonHeld = pressed
	}
	if button == MouseButtonRight {
		c.rightButtonHeld = pressed
	}
}

func (c *FPSCamera) ProcessKeyboardInput(state []uint8, delta float64) {
	w := state[sdl.GetScancodeFromKey(sdl.K_w)] != 0
	s := state[sdl.GetScancodeFromKey(sdl.K_s)] != 0
	a := state[sdl.GetScancodeFromKey(sdl.K_a)] != 0
	d := state[sdl.GetScancodeFromKey(sdl.K_d)] != 0

	if w {
		c.eye = c.eye.Add(c.front.Mul(c.speed * delta))
	}
	if s {
		c.eye = c.eye.Sub(c.front.Mul(c.speed * delta))
	}
	if a {
		c.eye = c.eye.Add(c.front.Cross(c.up).Normalize().Mul(c.speed * delta))
	}
	if d {
		c.eye = c.eye.Sub(c.front.Cross(c.up).Normalize().Mul(c.speed * delta))
	}

	if w || s || a || d {
		c.updateViewMatrix()
	}

}

func (c *FPSCamera) right() Vec3 {
	return c.front.Cross(c.up).Normalize()
}
