package main

func generateTriCube() []Face {
	return []Face{
		// front
		{points: [3]Vec3{{-1, -1, -1}, {-1, 1, -1}, {1, 1, -1}}, texcoords: [3]Tex{{0, 0}, {0, 1}, {1, 1}}, color: ColorWhite},
		{points: [3]Vec3{{-1, -1, -1}, {1, 1, -1}, {1, -1, -1}}, texcoords: [3]Tex{{0, 0}, {1, 1}, {1, 0}}, color: ColorWhite},
		// right
		{points: [3]Vec3{{1, -1, -1}, {1, 1, -1}, {1, 1, 1}}, texcoords: [3]Tex{{0, 0}, {0, 1}, {1, 1}}, color: ColorWhite},
		{points: [3]Vec3{{1, -1, -1}, {1, 1, 1}, {1, -1, 1}}, texcoords: [3]Tex{{0, 0}, {1, 1}, {1, 0}}, color: ColorWhite},
		// back
		{points: [3]Vec3{{1, -1, 1}, {1, 1, 1}, {-1, 1, 1}}, texcoords: [3]Tex{{0, 0}, {0, 1}, {1, 1}}, color: ColorWhite},
		{points: [3]Vec3{{1, -1, 1}, {-1, 1, 1}, {-1, -1, 1}}, texcoords: [3]Tex{{0, 0}, {1, 1}, {1, 0}}, color: ColorWhite},
		// left
		{points: [3]Vec3{{-1, -1, 1}, {-1, 1, 1}, {-1, 1, -1}}, texcoords: [3]Tex{{0, 0}, {0, 1}, {1, 1}}, color: ColorWhite},
		{points: [3]Vec3{{-1, -1, 1}, {-1, 1, -1}, {-1, -1, -1}}, texcoords: [3]Tex{{0, 0}, {1, 1}, {1, 0}}, color: ColorWhite},
		// top
		{points: [3]Vec3{{-1, 1, -1}, {-1, 1, 1}, {1, 1, 1}}, texcoords: [3]Tex{{0, 0}, {0, 1}, {1, 1}}, color: ColorWhite},
		{points: [3]Vec3{{-1, 1, -1}, {1, 1, 1}, {1, 1, -1}}, texcoords: [3]Tex{{0, 0}, {1, 1}, {1, 0}}, color: ColorWhite},
		// bottom
		{points: [3]Vec3{{1, -1, 1}, {-1, -1, 1}, {-1, -1, -1}}, texcoords: [3]Tex{{0, 0}, {0, 1}, {1, 1}}, color: ColorWhite},
		{points: [3]Vec3{{1, -1, 1}, {-1, -1, -1}, {1, -1, -1}}, texcoords: [3]Tex{{0, 0}, {1, 1}, {1, 0}}, color: ColorWhite},
	}
}
