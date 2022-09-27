package heretic

import "math"

//
// Vec2
//

type Vec2 struct{ X, Y float64 }

func (v Vec2) Add(u Vec2) Vec2    { return Vec2{v.X + u.X, v.Y + u.Y} }
func (v Vec2) Sub(u Vec2) Vec2    { return Vec2{v.X - u.X, v.Y - u.Y} }
func (v Vec2) Mul(f float64) Vec2 { return Vec2{v.X * f, v.Y * f} }
func (v Vec2) Div(f float64) Vec2 { return Vec2{v.X / f, v.Y / f} }
func (v Vec2) Dot(u Vec2) float64 { return v.X*u.X + v.Y*u.Y }
func (v Vec2) Length() float64    { return math.Sqrt(v.X*v.X + v.Y*v.Y) }
func (v Vec2) Normalize() Vec2    { l := v.Length(); return Vec2{v.X / l, v.Y / l} }

//
// Vec3
//

type Vec3 struct{ X, Y, Z float64 }

func (v Vec3) Add(u Vec3) Vec3    { return Vec3{v.X + u.X, v.Y + u.Y, v.Z + u.Z} }
func (v Vec3) Sub(u Vec3) Vec3    { return Vec3{v.X - u.X, v.Y - u.Y, v.Z - u.Z} }
func (v Vec3) Mul(f float64) Vec3 { return Vec3{v.X * f, v.Y * f, v.Z * f} }
func (v Vec3) Div(f float64) Vec3 { return Vec3{v.X / f, v.Y / f, v.Z / f} }
func (v Vec3) Length() float64    { return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z) }
func (v Vec3) Normalize() Vec3    { l := v.Length(); return Vec3{v.X / l, v.Y / l, v.Z / l} }
func (v Vec3) Vec4() Vec4         { return Vec4{v.X, v.Y, v.Z, 1} }
func (v Vec3) Dot(u Vec3) float64 { return v.X*u.X + v.Y*u.Y + v.Z*u.Z }
func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{v.Y*u.Z - v.Z*u.Y, v.Z*u.X - v.X*u.Z, v.X*u.Y - v.Y*u.X}
}

//
// Vec4
//

type Vec4 struct{ X, Y, Z, W float64 }

func (v Vec4) Add(u Vec4) Vec4    { return Vec4{v.X + u.X, v.Y + u.Y, v.Z + u.Z, v.W + v.W} }
func (v Vec4) Sub(u Vec4) Vec4    { return Vec4{v.X - u.X, v.Y - u.Y, v.Z - u.Z, v.W - v.W} }
func (v Vec4) Mul(f float64) Vec4 { return Vec4{v.X * f, v.Y * f, v.Z * f, v.W * f} }
func (v Vec4) Vec2() Vec2         { return Vec2{v.X, v.Y} }
func (v Vec4) Vec3() Vec3         { return Vec3{v.X, v.Y, v.Z} }
