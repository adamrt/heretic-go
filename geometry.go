package main

var cube = [9 * 9 * 9]Vec3{}

func getCube() []Vec3 {
	cube := []Vec3{}
	for x := float32(-1.0); x <= 1.0; x += .25 {
		for y := float32(-1.0); y <= 1.0; y += .25 {
			for z := float32(-1.0); z <= 1.0; z += .25 {
				point := Vec3{x, y, z}
				cube = append(cube, point)
			}
		}
	}
	return cube
}
