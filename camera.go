// This file contains a perspective camera. It's very basic currently.
package heretic

type FPSCamera struct {
	eye     Vec3
	front   Vec3
	worldUp Vec3

	yaw   float64
	pitch float64

	speed float64

	rightButtonPressed bool
	leftButtonPressed  bool
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

func (c *FPSCamera) MoveForward(deltaTime float64) {
	velocity := c.front.Mul(c.speed * deltaTime)
	c.eye = c.eye.Add(velocity)
}

func (c *FPSCamera) MoveBackward(deltaTime float64) {
	velocity := c.front.Mul(c.speed * deltaTime)
	c.eye = c.eye.Sub(velocity)
}

func (c *FPSCamera) MoveLeft(deltaTime float64) {
	velocity := c.right().Mul(c.speed * deltaTime)
	c.eye = c.eye.Add(velocity)
}

func (c *FPSCamera) MoveRight(deltaTime float64) {
	velocity := c.right().Mul(c.speed * deltaTime)
	c.eye = c.eye.Sub(velocity)
}

func (c *FPSCamera) Look(xrel, yrel int32) {
	c.yaw += float64(xrel) / 200
	c.pitch += float64(yrel) / 200
}

func (c *FPSCamera) Pan(xrel, yrel int32) {
	// X
	velocity := c.right().Mul(float64(xrel) / 400.0)
	c.eye = c.eye.Add(velocity)

	// Y
	velocity = c.up().Mul(float64(yrel) / 500.0)
	c.eye = c.eye.Add(velocity)
}

func (c *FPSCamera) right() Vec3 {
	return c.front.Cross(c.worldUp).Normalize()
}

func (c *FPSCamera) up() Vec3 {
	return c.right().Cross(c.front).Normalize()
}
