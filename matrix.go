package heretic

import "math"

// Matrix is a 4x4 row-major matrix
type Matrix [4][4]float64

// Mul multiplies a matrix a by matrix b.
func (a Matrix) Mul(b Matrix) Matrix {
	var m Matrix
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			m[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j] + a[i][3]*b[3][j]
		}
	}
	return m
}

// MulVec4 multiples a matrix by a Vec4 and returns a Vec4.
func (m Matrix) MulVec4(v Vec4) Vec4 {
	return Vec4{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]*v.W,
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]*v.W,
		m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]*v.W,
		m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]*v.W,
	}
}

// Return an Identity Matrix
// | 1  0  0  0 |
// | 0  1  0  0 |
// | 0  0  1  0 |
// | 0  0  0  0 |
func MatrixIdentity() Matrix {
	return Matrix{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Return a Scale Matrix
// | sx  0  0  0 |
// |  0 sy  0  0 |
// |  0  0 sx  0 |
// |  0  0  0  1 |
func NewScaleMatrix(v Vec3) Matrix {
	m := MatrixIdentity()
	m[0][0] = v.X
	m[1][1] = v.Y
	m[2][2] = v.Z
	return m
}

// Return a Translation Matrix
// | 1  0  0  tx |
// | 0  1  0  ty |
// | 0  0  1  tz |
// | 0  0  0   1 |
func NewTranslationMatrix(v Vec3) Matrix {
	m := MatrixIdentity()
	m[0][3] = v.X
	m[1][3] = v.Y
	m[2][3] = v.Z
	return m
}

// Sugar function to run x, y and z rotation matrix functions.
func NewRotationMatrix(v Vec3) Matrix {
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
	m[1][1] = c
	m[1][2] = -s
	m[2][1] = s
	m[2][2] = c
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
	m[0][0] = c
	m[0][2] = s
	m[2][0] = -s
	m[2][2] = c
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
	m[0][0] = c
	m[0][1] = -s
	m[1][0] = s
	m[1][1] = c
	return m
}

// Return a Perspective Projection Matrix
//
// The 3/2==1 stores the original z value for use in MulProjection so we can do
// perspective divide in MulVec4Proj().
func MatrixMakePerspective(fov, aspect, znear, zfar float64) Matrix {
	m := Matrix{}
	m[0][0] = aspect * (1 / math.Tan(fov/2))
	m[1][1] = 1 / math.Tan(fov/2)
	m[2][2] = zfar / (zfar - znear)
	m[2][3] = (-zfar * znear) / (zfar - znear)
	m[3][2] = 1.0
	return m
}

func MatrixMakeOrtho(left, right, bottom, top, near, far float64) Matrix {
	m := MatrixIdentity()
	m[0][0] = 2 / (right - left)
	m[1][1] = 2 / (top - bottom)
	m[2][2] = -2 / (far - near)
	m[0][3] = -(right + left) / (right - left)
	m[1][3] = -(top + bottom) / (top - bottom)
	m[2][3] = -(far + near) / (far - near)
	return m
}

func LookAt(eye, target, up Vec3) Matrix {
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
