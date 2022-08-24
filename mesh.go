package main

type Mesh struct {
	// This is just a slice of a slice, but for naming purposes, triangles
	// makes more sense, since that is what it represents.
	triangles [][3]Vec3

	rotation Vec3
}
