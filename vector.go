package main

const FOV float64 = 640

type Vec2 struct{ x, y float64 }

type Vec3 struct{ x, y, z float64 }

// Project returns a 2D project point from a 3D point using Perspective Divide.
func (v Vec3) Project() Vec2 {
	projectedPoint := Vec2{
		x: (FOV * v.x) / v.z,
		y: (FOV * v.y) / v.z,
	}
	return projectedPoint
}
