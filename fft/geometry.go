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
		t.textureData.a,
		t.textureData.b,
		t.textureData.c,
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

func texWithPage(uv heretic.Tex, page int) heretic.Tex {
	v := float64(int(uv.V) + page*256)
	return heretic.Tex{U: uv.U / 255, V: v / 1023.0}
}

type triangleTexData struct {
	a, b, c heretic.Tex
	palette int
}

type quadTexData struct {
	a, b, c, d heretic.Tex
	palette    int
}

func (q quadTexData) split() []triangleTexData {
	return []triangleTexData{
		{a: q.a, b: q.b, c: q.c, palette: q.palette},
		{a: q.b, b: q.d, c: q.c, palette: q.palette},
	}
}
