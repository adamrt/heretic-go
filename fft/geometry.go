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

type triangle struct {
	a, b, c     vertex
	textureData triangleTexData
	palette     *heretic.Palette
}

func (t triangle) face() heretic.Face {
	return heretic.NewFace(
		t.points(),
		t.texcoords(),
		t.palette,
		heretic.ColorWhite,
	)
}

func (t triangle) points() [3]heretic.Vec3 {
	return [3]heretic.Vec3{t.a.vec3(), t.b.vec3(), t.c.vec3()}
}

func (t triangle) texcoords() [3]heretic.Tex {
	return [3]heretic.Tex{
		t.textureData.a.tex(t.textureData.page),
		t.textureData.b.tex(t.textureData.page),
		t.textureData.c.tex(t.textureData.page),
	}
}

type quad struct {
	a, b, c, d vertex
}

func (q quad) split() []triangle {
	return []triangle{
		{a: q.a, b: q.b, c: q.c},
		{a: q.b, b: q.d, c: q.c},
	}
}

type uv struct {
	x, y uint8
}

func (t uv) tex(page int) heretic.Tex {
	y := int(t.y) + page*256
	return heretic.NewTex(float64(t.x)/255, float64(y)/1023.0)
}

type triangleTexData struct {
	a, b, c uv
	palette int
	page    int
}

type quadTexData struct {
	a, b, c, d uv
	palette    int
	page       int
}

func (q quadTexData) split() []triangleTexData {
	return []triangleTexData{
		{a: q.a, b: q.b, c: q.c, palette: q.palette, page: q.page},
		{a: q.b, b: q.d, c: q.c, palette: q.palette, page: q.page},
	}
}
