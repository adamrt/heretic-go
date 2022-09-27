// This file contains a perspective camera. It's very basic currently.
package heretic

type Camera struct {
	position  Vec3
	direction Vec3
	right     Vec3
	up        Vec3
	worldUp   Vec3

	yaw   float64
	pitch float64

	speed float64

	rightButtonPressed bool
	leftButtonPressed  bool
}

func NewCamera(position, direction Vec3) Camera {
	return Camera{
		direction: direction,
		position:  position,
		worldUp:   Vec3{0, 1, 0},
		speed:     2.0,
	}
}

func (c *Camera) LookAtTarget() Vec3 {
	target := Vec3{0, 0, 1}
	yawRotation := MatrixMakeRotY(c.yaw)
	pitchRotation := MatrixMakeRotX(c.pitch)

	cameraRotation := MatrixIdentity()
	cameraRotation = pitchRotation.Mul(cameraRotation)
	cameraRotation = yawRotation.Mul(cameraRotation)

	c.direction = cameraRotation.MulVec4(target.Vec4()).Vec3()
	c.right = c.direction.Cross(c.worldUp).Normalize()
	c.up = c.right.Cross(c.direction).Normalize()

	target = c.position.Add(c.direction)
	return target
}

func (c *Camera) LookAt(eye, target, up Vec3) Matrix {
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

func (c *Camera) MoveForward(deltaTime float64) {
	velocity := c.direction.Mul(c.speed * deltaTime)
	c.position = c.position.Add(velocity)
}

func (c *Camera) MoveBackward(deltaTime float64) {
	velocity := c.direction.Mul(c.speed * deltaTime)
	c.position = c.position.Sub(velocity)
}

func (c *Camera) MoveLeft(deltaTime float64) {
	velocity := c.right.Mul(c.speed * deltaTime)
	c.position = c.position.Add(velocity)
}

func (c *Camera) MoveRight(deltaTime float64) {
	velocity := c.right.Mul(c.speed * deltaTime)
	c.position = c.position.Sub(velocity)
}

func (c *Camera) Look(xrel, yrel int32) {
	c.yaw += float64(xrel) / 200
	c.pitch += float64(yrel) / 200
}

func (c *Camera) Pan(xrel, yrel int32) {
	// X
	velocity := c.right.Mul(float64(xrel) / 400.0)
	c.position = c.position.Add(velocity)

	// Y
	velocity = c.up.Mul(float64(yrel) / 500.0)
	c.position = c.position.Add(velocity)
}
