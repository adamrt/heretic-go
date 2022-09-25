package heretic

import "math"

type Matrix struct {
	m [4][4]float64
}

func (a Matrix) Mul(b Matrix) Matrix {
	var m Matrix
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			m.m[i][j] = a.m[i][0]*b.m[0][j] + a.m[i][1]*b.m[1][j] + a.m[i][2]*b.m[2][j] + a.m[i][3]*b.m[3][j]
		}
	}
	return m
}

func (m Matrix) MulVec4(v Vec4) Vec4 {
	return Vec4{
		m.m[0][0]*v.X + m.m[0][1]*v.Y + m.m[0][2]*v.Z + m.m[0][3]*v.W,
		m.m[1][0]*v.X + m.m[1][1]*v.Y + m.m[1][2]*v.Z + m.m[1][3]*v.W,
		m.m[2][0]*v.X + m.m[2][1]*v.Y + m.m[2][2]*v.Z + m.m[2][3]*v.W,
		m.m[3][0]*v.X + m.m[3][1]*v.Y + m.m[3][2]*v.Z + m.m[3][3]*v.W,
	}
}

func (m Matrix) MulVec4Proj(v Vec4) Vec4 {
	// Multiply the original projection matrix by the vector
	result := m.MulVec4(v)

	// Perspective Divide with original z value (result.w).  The result.w is
	// populated during MulVec4() because of the projection matrix 3/2==1.
	if result.W != 0.0 {
		result.X /= result.W
		result.Y /= result.W
		result.Z /= result.W
	}
	return result
}

// Return an Identity Matrix
// | 1  0  0  0 |
// | 0  1  0  0 |
// | 0  0  1  0 |
// | 0  0  0  0 |
func MatrixIdentity() Matrix {
	return Matrix{[4][4]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}}
}

// Return a Scale Matrix
// | sx  0  0  0 |
// |  0 sy  0  0 |
// |  0  0 sx  0 |
// |  0  0  0  1 |
func MatrixMakeScale(v Vec3) Matrix {
	m := MatrixIdentity()
	m.m[0][0] = v.X
	m.m[1][1] = v.Y
	m.m[2][2] = v.Z
	return m
}

// Return a Translation Matrix
// | 1  0  0  tx |
// | 0  1  0  ty |
// | 0  0  1  tz |
// | 0  0  0   1 |
func MatrixMakeTrans(v Vec3) Matrix {
	m := MatrixIdentity()
	m.m[0][3] = v.X
	m.m[1][3] = v.Y
	m.m[2][3] = v.Z
	return m
}

// Sugar function to run x, y and z rotation matrix functions.
func MatrixMakeRotation(v Vec3) Matrix {
	x := MatrixMakeRotX(v.X)
	y := MatrixMakeRotY(v.Y)
	z := MatrixMakeRotZ(v.Z)
	return x.Mul(y).Mul(z)
}

// Return a Rotation Matrix for x axis
// | 1  0  0  0 |
// | 0  c -s  0 |
// | 0  s  c  0 |
// | 0  0  0  1 |
func MatrixMakeRotX(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)

	m := MatrixIdentity()
	m.m[1][1] = c
	m.m[1][2] = -s
	m.m[2][1] = s
	m.m[2][2] = c
	return m
}

// Return a Rotation Matrix for y axis
// | c  0  s  0 |
// | 0  1  0  0 |
// |-s  0  c  0 |
// | 0  0  0  1 |
func MatrixMakeRotY(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)

	m := MatrixIdentity()
	m.m[0][0] = c
	m.m[0][2] = s
	m.m[2][0] = -s
	m.m[2][2] = c
	return m
}

// Return a Rotation Matrix for z axis
// | c -s  0  0 |
// | s  c  0  0 |
// | 0  0  1  0 |
// | 0  0  0  1 |
func MatrixMakeRotZ(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)

	m := MatrixIdentity()
	m.m[0][0] = c
	m.m[0][1] = -s
	m.m[1][0] = s
	m.m[1][1] = c
	return m
}

// Return a Perspective Projection Matrix
//
// The 3/2==1 stores the original z value for use in MulProjection so we can do
// perspective divide in MulVec4Proj().
func MatrixMakePerspective(fov, aspect, znear, zfar float64) Matrix {
	m := Matrix{}
	m.m[0][0] = aspect * (1 / math.Tan(fov/2))
	m.m[1][1] = 1 / math.Tan(fov/2)
	m.m[2][2] = zfar / (zfar - znear)
	m.m[2][3] = (-zfar * znear) / (zfar - znear)
	m.m[3][2] = 1.0
	return m
}
