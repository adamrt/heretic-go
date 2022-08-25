package main

import "math"

type Mat4 struct {
	m [4][4]float64
}

func (m Mat4) MulVec4(v Vec4) Vec4 {
	return Vec4{
		m.m[0][0]*v.x + m.m[0][1]*v.y + m.m[0][2]*v.z + m.m[0][3]*v.w,
		m.m[1][0]*v.x + m.m[1][1]*v.y + m.m[1][2]*v.z + m.m[1][3]*v.w,
		m.m[2][0]*v.x + m.m[2][1]*v.y + m.m[2][2]*v.z + m.m[2][3]*v.w,
		m.m[3][0]*v.x + m.m[3][1]*v.y + m.m[3][2]*v.z + m.m[3][3]*v.w,
	}
}

// Return an Identity Matrix
// | 1  0  0  0 |
// | 0  1  0  0 |
// | 0  0  1  0 |
// | 0  0  0  0 |
func Mat4Identity() Mat4 {
	return Mat4{[4][4]float64{
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
func Mat4MakeScale(sx, sy, sz float64) Mat4 {
	m := Mat4Identity()
	m.m[0][0] = sx
	m.m[1][1] = sy
	m.m[2][2] = sz
	return m
}

// Return a Translation Matrix
// | 1  0  0  tx |
// | 0  1  0  ty |
// | 0  0  1  tz |
// | 0  0  0   1 |
func Mat4MakeTrans(tx, ty, tz float64) Mat4 {
	m := Mat4Identity()
	m.m[0][3] = tx
	m.m[1][3] = ty
	m.m[2][3] = tz
	return m
}

// Return a Rotation Matrix for x axis
// | 1  0  0  0 |
// | 0  c -s  0 |
// | 0  s  c  0 |
// | 0  0  0  1 |
func Mat4MakeRotX(angle float64) Mat4 {
	c := math.Cos(angle)
	s := math.Sin(angle)

	m := Mat4Identity()
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
func Mat4MakeRotY(angle float64) Mat4 {
	c := math.Cos(angle)
	s := math.Sin(angle)

	m := Mat4Identity()
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
func Mat4MakeRotZ(angle float64) Mat4 {
	c := math.Cos(angle)
	s := math.Sin(angle)

	m := Mat4Identity()
	m.m[0][0] = c
	m.m[0][1] = -s
	m.m[1][0] = s
	m.m[1][1] = c
	return m
}
