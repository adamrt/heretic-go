// This file contains the structs that will be read from the FFT ISO.
//
// FFT mesh data is primarily represented by quads, but our home grown engine
// only handles triangles. The quads are read in then split into two triangles.
// This has to be done for geometry, normals and texture coordinates.
//
// There are methods such as vertex.vec3() that will convert the local type to
// the engine type.
package fft

import (
	"github.com/adamrt/heretic"
)

type normal struct {
	x, y, z float64
}

type triangle struct {
	points      []heretic.Vec3
	textureData triangleTexData
	palette     heretic.Palette
}

func (t triangle) triangle() heretic.Triangle {
	return heretic.Triangle{
		Points:    t.points,
		Texcoords: t.texcoords(),
		Palette:   t.palette,
		Color:     heretic.ColorBlack,
	}
}

func (t triangle) texcoords() []heretic.Tex {
	return []heretic.Tex{
		t.textureData.a.tex(t.textureData.page),
		t.textureData.b.tex(t.textureData.page),
		t.textureData.c.tex(t.textureData.page),
	}
}

type quad struct {
	a, b, c, d heretic.Vec3
}

func (q quad) split() []triangle {
	return []triangle{
		{points: []heretic.Vec3{q.a, q.b, q.c}},
		{points: []heretic.Vec3{q.b, q.d, q.c}},
	}
}

type uv struct {
	x, y uint8
}

func (t uv) tex(page int) heretic.Tex {
	y := int(t.y) + page*256
	return heretic.Tex{U: float64(t.x) / 255, V: float64(y) / 1023.0}
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
