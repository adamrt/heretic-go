package main

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
func Mat4Identity() Mat4 {
	return Mat4{[4][4]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}}
}

func Mat4MakeScale(sx, sy, sz float64) Mat4 {
	m := Mat4Identity()
	m.m[0][0] = sx
	m.m[1][1] = sy
	m.m[2][2] = sz
	return m
}
