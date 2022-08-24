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
// Vec3D
//

type Vec3 struct{ x, y, z float64 }

func (v Vec3) Length() float64         { return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z) }
func (v Vec3) Normalize() Vec3         { l := v.Length(); return Vec3{v.x / l, v.y / l, v.z / l} }
func (v Vec3) Add(u Vec3) Vec3         { return Vec3{v.x + u.x, v.y + u.y, v.z + u.z} }
func (v Vec3) Sub(u Vec3) Vec3         { return Vec3{v.x - u.x, v.y - u.y, v.z - u.z} }
func (v Vec3) Mul(factor float64) Vec3 { return Vec3{v.x * factor, v.y * factor, v.z * factor} }
func (v Vec3) Div(factor float64) Vec3 { return Vec3{v.x / factor, v.y / factor, v.z / factor} }
func (v Vec3) Dot(u Vec3) float64      { return v.x*u.x + v.y*u.y + v.z*u.z }
func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{v.y*u.z - v.z*u.y, v.z*u.x - v.x*u.z, v.x*u.y - v.y*u.x}
}

// Project returns a 2D project point from a 3D point using Perspective Divide.
// This currently scales the points from (-1 to 1) to (-800 to 800). This would
// put the points off screen, but the 'perspective divide' division scales it
// back down.
func (v Vec3) Project() Vec2 { return Vec2{x: (FOV * v.x) / v.z, y: (FOV * v.y) / v.z} }

func (v Vec3) RotateX(angle float64) Vec3 {
	return Vec3{
		x: v.x,
		y: v.y*math.Cos(angle) - v.z*math.Sin(angle),
		z: v.y*math.Sin(angle) + v.z*math.Cos(angle),
	}
}
func (v Vec3) RotateY(angle float64) Vec3 {
	return Vec3{
		x: v.x*math.Cos(angle) - v.z*math.Sin(angle),
		y: v.y,
		z: v.x*math.Sin(angle) + v.z*math.Cos(angle),
	}

}
func (v Vec3) RotateZ(angle float64) Vec3 {
	return Vec3{
		x: v.x*math.Cos(angle) - v.y*math.Sin(angle),
		y: v.x*math.Sin(angle) + v.y*math.Cos(angle),
		z: v.z,
	}
}
