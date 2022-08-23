package main

func getCube() []Vec3 {
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
