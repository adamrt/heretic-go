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
	"math"

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

// normalizeTriangle accepts a face and the min/max of all face coordinates in a
// mesh and normalizeds them between 0 and 1. This scales down large models
// during import.  This is primary used for loading FFT maps since they have
// very large coordinates.
func normalizeTriangle(t triangle, min, max float64) triangle {
	normalized := make([]heretic.Vec3, 3)
	for i := 0; i < 3; i++ {
		normalized[i].X = normalize(t.points[i].X, min, max)
		normalized[i].Y = normalize(t.points[i].Y, min, max)
		normalized[i].Z = normalize(t.points[i].Z, min, max)
	}

	t.points = normalized
	return t
}

// Normalize takes the min/max of face coordinates in a mesh and then normalizes
// them to the 0.0-1.0 space.
func normalize(x, min, max float64) float64 {
	return (x - min) / (max - min)
}

// centerTranslation returns a translation vector that will center the mesh.
func centerTraslation(triangles []triangle) heretic.Vec3 {
	var minx float64 = math.MaxInt16
	var maxx float64 = math.MinInt16
	var miny float64 = math.MaxInt16
	var maxy float64 = math.MinInt16
	var minz float64 = math.MaxInt16
	var maxz float64 = math.MinInt16

	for _, t := range triangles {

		// Each point for max
		for i := 0; i < 3; i++ {
			// Max
			if t.points[i].X > maxx {
				maxx = t.points[i].X
			}
			if t.points[i].Y > maxy {
				maxy = t.points[i].Y
			}
			if t.points[i].Z > maxz {
				maxz = t.points[i].Z
			}

			// Min
			if t.points[i].X < minx {
				minx = t.points[i].X
			}
			if t.points[i].Y < miny {
				miny = t.points[i].Y
			}
			if t.points[i].Z < minz {
				minz = t.points[i].Z
			}
		}
	}

	// Not applying the Y coordinate since FFT maps already sit on the
	// floor. Adding the Y tranlation would put the floor at the mad 1/2
	// height midway point.
	return heretic.Vec3{
		X: -(maxx + minx) / 2.0,
		// Y: -(maxy + miny) / 2.0,
		Z: -(maxz + minz) / 2.0,
	}
}
