package main

import "math"

//
// Vec2
//

type Vec2 struct{ x, y float64 }

func (v Vec2) Length() float64         { return math.Sqrt(v.x*v.x + v.y*v.y) }
func (v Vec2) Normalize() Vec2         { l := v.Length(); return Vec2{v.x / l, v.y / l} }
func (v Vec2) Add(u Vec2) Vec2         { return Vec2{v.x + u.x, v.y + u.y} }
func (v Vec2) Sub(u Vec2) Vec2         { return Vec2{v.x - u.x, v.y - u.y} }
func (v Vec2) Mul(factor float64) Vec2 { return Vec2{v.x * factor, v.y * factor} }
func (v Vec2) Div(factor float64) Vec2 { return Vec2{v.x / factor, v.y / factor} }
func (v Vec2) Dot(u Vec2) float64      { return v.x*u.x + v.y*u.y }

//
// Vec3
//

type Vec3 struct{ x, y, z float64 }

func (v Vec3) Length() float64         { return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z) }
func (v Vec3) Normalize() Vec3         { l := v.Length(); return Vec3{v.x / l, v.y / l, v.z / l} }
func (v Vec3) Add(u Vec3) Vec3         { return Vec3{v.x + u.x, v.y + u.y, v.z + u.z} }
func (v Vec3) Sub(u Vec3) Vec3         { return Vec3{v.x - u.x, v.y - u.y, v.z - u.z} }
func (v Vec3) Mul(factor float64) Vec3 { return Vec3{v.x * factor, v.y * factor, v.z * factor} }
func (v Vec3) Div(factor float64) Vec3 { return Vec3{v.x / factor, v.y / factor, v.z / factor} }
func (v Vec3) Dot(u Vec3) float64      { return v.x*u.x + v.y*u.y + v.z*u.z }
func (v Vec3) Vec4() Vec4              { return Vec4{v.x, v.y, v.z, 1} }
func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{v.y*u.z - v.z*u.y, v.z*u.x - v.x*u.z, v.x*u.y - v.y*u.x}
}

//
// Vec4
//

type Vec4 struct{ x, y, z, w float64 }

func (v Vec4) Add(u Vec4) Vec4 { return Vec4{v.x + u.x, v.y + u.y, v.z + u.z, v.w + v.w} }
func (v Vec4) Sub(u Vec4) Vec4 { return Vec4{v.x - u.x, v.y - u.y, v.z - u.z, v.w - v.w} }
func (v Vec4) Mul(factor float64) Vec4 {
	return Vec4{v.x * factor, v.y * factor, v.z * factor, v.w * factor}
}
func (v Vec4) Vec3() Vec3 { return Vec3{v.x, v.y, v.z} }
func (v Vec4) Vec2() Vec2 { return Vec2{v.x, v.y} }
