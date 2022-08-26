package main

type Camera struct {
	position  Vec3
	direction Vec3
	velocity  Vec3
	yaw       float64
	pitch     float64
}

func NewCamera(position, direction Vec3) Camera {
	return Camera{
		position:  position,
		direction: direction,
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
	target = c.position.Add(c.direction)
	return target
}

func (c *Camera) LookAtMatrix(target, up Vec3) Matrix {
	z := target.Sub(c.position).Normalize()
	x := up.Cross(z).Normalize()
	y := z.Cross(x)

	// View Matrix
	return Matrix{m: [4][4]float64{
		{x.x, x.y, x.z, -x.Dot(c.position)},
		{y.x, y.y, y.z, -y.Dot(c.position)},
		{z.x, z.y, z.z, -z.Dot(c.position)},
		{0, 0, 0, 1},
	}}
}
