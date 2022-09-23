// This file contains the structs that will be read from the FFT ISO.
//
// FFT mesh data is primarily represented by quads, but our home grown engine
// only handles triangles. The quads are read in then split into two triangles.
// This has to be done for geometry, normals and texture coordinates.
//
// There are methods such as vertex.vec3() that will convert the local type to
// the engine type.
package fft

import "github.com/adamrt/heretic"

type vertex struct {
	x, y, z int16
}

func (v vertex) vec3() heretic.Vec3 {
	return heretic.NewVec3(float64(v.x), float64(v.y), float64(v.z))
}

type normal struct {
	x, y, z float64
}

type polygonTexData struct {
	u, v       uint8
	colorPoint uint8
}

func (t polygonTexData) tex() heretic.Tex {
	return heretic.NewTex(float64(t.u), float64(t.v))
}

type triangle struct {
	a, b, c vertex
}

func (t triangle) points() [3]heretic.Vec3 {
	return [3]heretic.Vec3{t.a.vec3(), t.b.vec3(), t.c.vec3()}
}

type quad struct{ a, b, c, d vertex }

func (q quad) split() []triangle {
	return []triangle{
		{q.a, q.b, q.c},
		{q.b, q.d, q.c},
	}
}
