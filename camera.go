// This file contains a perspective camera. It's very basic currently.
package heretic

import (
	"github.com/veandco/go-sdl2/sdl"
)

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

	rightButtonHeld bool
	leftButtonHeld  bool
}

func NewFPSCamera(eye, front, up Vec3) FPSCamera {
	return FPSCamera{
		front: front,
		eye:   eye,
		up:    up,
		speed: 2.0,
	}
}

func (c *FPSCamera) processMouseMovement(xrel, yrel, delta float64) {
	sensitiviy := 0.3
	if c.leftButtonHeld {
		c.yaw += xrel * sensitiviy * delta
		c.pitch += yrel * sensitiviy * delta

		if c.pitch > 89.0 {
			c.pitch = 89.0
		}
		if c.pitch < -89.0 {
			c.pitch = -89.0
		}

		target := Vec3{0, 0, 1}
		cameraRotation := MatrixIdentity().Mul(MatrixMakeRotY(c.yaw)).Mul(MatrixMakeRotX(c.pitch))
		direction := cameraRotation.MulVec4(target.Vec4()).Vec3()
		c.front = direction

	} else if c.rightButtonHeld {
		// X
		c.eye = c.eye.Add(c.right().Mul(xrel * delta * sensitiviy))

		// Y
		c.eye = c.eye.Add(c.up.Mul(yrel * delta * sensitiviy))

	}
}

func (c *FPSCamera) processMouseWheel(y float64, delta float64) {
	velocity := c.front.Mul(c.speed * y * delta)
	c.eye = c.eye.Add(velocity)
}

func (c *FPSCamera) processMouseButton(button MouseButton, pressed bool) {
	if button == MouseButtonLeft {
		c.leftButtonHeld = pressed
	}
	if button == MouseButtonRight {
		c.rightButtonHeld = pressed
	}
}

func (c *FPSCamera) processKeyboardInput(state []uint8, delta float64) {
	if state[sdl.GetScancodeFromKey(sdl.K_w)] != 0 {
		c.eye = c.eye.Add(c.front.Mul(c.speed * delta))
	}
	if state[sdl.GetScancodeFromKey(sdl.K_s)] != 0 {
		c.eye = c.eye.Sub(c.front.Mul(c.speed * delta))
	}
	if state[sdl.GetScancodeFromKey(sdl.K_a)] != 0 {
		c.eye = c.eye.Add(c.front.Cross(c.up).Normalize().Mul(c.speed * delta))
	}
	if state[sdl.GetScancodeFromKey(sdl.K_d)] != 0 {
		c.eye = c.eye.Sub(c.front.Cross(c.up).Normalize().Mul(c.speed * delta))
	}
}

func (c *FPSCamera) right() Vec3 {
	return c.front.Cross(c.up).Normalize()
}
