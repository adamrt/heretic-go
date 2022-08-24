package main

import "math"

// FOV is just matching the width of the screen currently.  This is used to
// scale up the positions of the points from
const FOV float64 = WindowWidth

type Vec2 struct{ x, y float64 }

type Vec3 struct{ x, y, z float64 }

// Project returns a 2D project point from a 3D point using Perspective Divide.
// This currently scales the points from (-1 to 1) to (-800 to 800). This would
// put the points off screen, but the 'perspective divide' division scales it
// back down.
func (v Vec3) Project() Vec2 {
	projectedPoint := Vec2{
		x: (FOV * v.x) / v.z,
		y: (FOV * v.y) / v.z,
	}
	return projectedPoint
}

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
