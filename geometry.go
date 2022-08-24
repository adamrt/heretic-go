package main

func generateTriCube() []Triangle {
	return []Triangle{
		// front
		{points: [3]Vec3{{-1, -1, -1}, {-1, 1, -1}, {1, 1, -1}}},
		{points: [3]Vec3{{-1, -1, -1}, {1, 1, -1}, {1, -1, -1}}},
		// right
		{points: [3]Vec3{{1, -1, -1}, {1, 1, -1}, {1, 1, 1}}},
		{points: [3]Vec3{{1, -1, -1}, {1, 1, 1}, {1, -1, 1}}},
		// back
		{points: [3]Vec3{{1, -1, 1}, {1, 1, 1}, {-1, 1, 1}}},
		{points: [3]Vec3{{1, -1, 1}, {-1, 1, 1}, {-1, -1, 1}}},
		// left
		{points: [3]Vec3{{-1, -1, 1}, {-1, 1, 1}, {-1, 1, -1}}},
		{points: [3]Vec3{{-1, -1, 1}, {-1, 1, -1}, {-1, -1, -1}}},
		// top
		{points: [3]Vec3{{-1, 1, -1}, {-1, 1, 1}, {1, 1, 1}}},
		{points: [3]Vec3{{-1, 1, -1}, {1, 1, 1}, {1, 1, -1}}},
		// bottom
		{points: [3]Vec3{{1, -1, 1}, {-1, -1, 1}, {-1, -1, -1}}},
		{points: [3]Vec3{{1, -1, 1}, {-1, -1, -1}, {1, -1, -1}}},
	}
}

func generateDotCube() []Vec3 {
	cube := []Vec3{}
	for x := -1.0; x <= 1.0; x += .25 {
		for y := -1.0; y <= 1.0; y += .25 {
			for z := -1.0; z <= 1.0; z += .25 {
				point := Vec3{x, y, z}
				cube = append(cube, point)
			}
		}
	}
	return cube
}
