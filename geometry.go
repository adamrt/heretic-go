package main

func generateTriCube() []Face {
	return []Face{
		// front
		{points: [3]Vec3{{-1, -1, -1}, {-1, 1, -1}, {1, 1, -1}}, color: ColorRed},
		{points: [3]Vec3{{-1, -1, -1}, {1, 1, -1}, {1, -1, -1}}, color: ColorRed},
		// right
		{points: [3]Vec3{{1, -1, -1}, {1, 1, -1}, {1, 1, 1}}, color: ColorGreen},
		{points: [3]Vec3{{1, -1, -1}, {1, 1, 1}, {1, -1, 1}}, color: ColorGreen},
		// back
		{points: [3]Vec3{{1, -1, 1}, {1, 1, 1}, {-1, 1, 1}}, color: ColorBlue},
		{points: [3]Vec3{{1, -1, 1}, {-1, 1, 1}, {-1, -1, 1}}, color: ColorBlue},
		// left
		{points: [3]Vec3{{-1, -1, 1}, {-1, 1, 1}, {-1, 1, -1}}, color: ColorYellow},
		{points: [3]Vec3{{-1, -1, 1}, {-1, 1, -1}, {-1, -1, -1}}, color: ColorYellow},
		// top
		{points: [3]Vec3{{-1, 1, -1}, {-1, 1, 1}, {1, 1, 1}}, color: ColorCyan},
		{points: [3]Vec3{{-1, 1, -1}, {1, 1, 1}, {1, 1, -1}}, color: ColorCyan},
		// bottom
		{points: [3]Vec3{{1, -1, 1}, {-1, -1, 1}, {-1, -1, -1}}, color: ColorMagenta},
		{points: [3]Vec3{{1, -1, 1}, {-1, -1, -1}, {1, -1, -1}}, color: ColorMagenta},
	}
}
