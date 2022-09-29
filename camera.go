// This file contains a perspective camera. It's very basic currently.
package heretic

import "github.com/veandco/go-sdl2/sdl"

type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonMiddle
	MouseButtonRight
)

type FPSCamera struct {
	eye     Vec3
	front   Vec3
	worldUp Vec3

	yaw   float64
	pitch float64

	speed float64

	rightButtonHeld bool
	leftButtonHeld  bool
}

func NewFPSCamera(eye, front Vec3) FPSCamera {
	return FPSCamera{
		front:   front,
		eye:     eye,
		worldUp: Vec3{0, 1, 0},
		speed:   2.0,
	}
}

func (c *FPSCamera) LookAt(eye, target, up Vec3) Matrix {
	// Forward
	z := target.Sub(eye).Normalize()
	// Right
	x := up.Cross(z).Normalize()
	// Up
	y := z.Cross(x).Normalize()

	viewMatrix := Matrix{
		{x.X, x.Y, x.Z, -x.Dot(eye)},
		{y.X, y.Y, y.Z, -y.Dot(eye)},
		{z.X, z.Y, z.Z, -z.Dot(eye)},
		{0, 0, 0, 1},
	}

	return viewMatrix
}

func (c *FPSCamera) processMouseMovement(xrel, yrel, delta float64) {
	if c.leftButtonHeld {
		c.yaw += xrel * delta / 4
		c.pitch += yrel * delta / 4
	} else if c.rightButtonHeld {
		// X
		velocity := c.right().Mul(float64(xrel) / 400.0)
		c.eye = c.eye.Add(velocity)

		// Y
		velocity = c.up().Mul(float64(yrel) / 500.0)
		c.eye = c.eye.Add(velocity)
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
		velocity := c.front.Mul(c.speed * delta)
		c.eye = c.eye.Add(velocity)
	}
	if state[sdl.GetScancodeFromKey(sdl.K_s)] != 0 {
		velocity := c.front.Mul(c.speed * delta)
		c.eye = c.eye.Sub(velocity)
	}
	if state[sdl.GetScancodeFromKey(sdl.K_a)] != 0 {
		velocity := c.right().Mul(c.speed * delta)
		c.eye = c.eye.Add(velocity)
	}
	if state[sdl.GetScancodeFromKey(sdl.K_d)] != 0 {
		velocity := c.right().Mul(c.speed * delta)
		c.eye = c.eye.Sub(velocity)
	}
}

func (c *FPSCamera) right() Vec3 {
	return c.front.Cross(c.worldUp).Normalize()
}

func (c *FPSCamera) up() Vec3 {
	return c.right().Cross(c.front).Normalize()
}
