package main

type Vec2 struct{ x, y float32 }
type Vec3 struct{ x, y, z float32 }

const FOV float32 = 128

func project(p Vec3) Vec2 {
	projectedPoint := Vec2{
		FOV * p.x,
		FOV * p.y,
	}
	return projectedPoint
}
