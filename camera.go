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
		speed:     15.0,
	}
}

func (c *Camera) LookAtTarget(target Vec3) Vec3 {
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

func (c *Camera) LookAtMatrix(target, up Vec3) Matrix {
	z := target.Sub(c.position).Normalize()
	x := up.Cross(z).Normalize()
	y := z.Cross(x)

	// View Matrix
	return Matrix{m: [4][4]float64{
		{x.X, x.Y, x.Z, -x.Dot(c.position)},
		{y.X, y.Y, y.Z, -y.Dot(c.position)},
		{z.X, z.Y, z.Z, -z.Dot(c.position)},
		{0, 0, 0, 1},
	}}
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

func (c *Camera) Pan(xrel, yrel int32) {
	// X
	velocity := c.right.Mul(float64(xrel) / 50.0)
	c.position = c.position.Add(velocity)

	// Y
	velocity = c.up.Mul(float64(yrel) / 50.0)
	c.position = c.position.Add(velocity)
}
