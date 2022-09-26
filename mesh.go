package heretic

import (
	"math"
)

func NewMesh(triangles []Triangle, texture Texture) Mesh {
	return Mesh{Triangles: triangles, Texture: texture}
}

type Mesh struct {
	// This is just a slice of a slice, but for naming purposes, triangles
	// makes more sense, since that is what it represents.
	Triangles  []Triangle
	Texture    Texture
	Background *Background

	Rotation    Vec3
	Scale       Vec3
	Translation Vec3

	trianglesToRender []Triangle
}

// NormalizeCoordinates normalizes all vertex coordinates between 0 and 1. This
// scales down large models during import.  This is primary used for loading FFT
// maps since they have very large coordinates.  The min/max values should be
// the min and max of
func (m *Mesh) NormalizeCoordinates() {
	min, max := m.coordMinMax()
	for i := 0; i < len(m.Triangles); i++ {
		for j := 0; j < 3; j++ {
			m.Triangles[i].Points[j].X = normalize(m.Triangles[i].Points[j].X, min, max)
			m.Triangles[i].Points[j].Y = normalize(m.Triangles[i].Points[j].Y, min, max)
			m.Triangles[i].Points[j].Z = normalize(m.Triangles[i].Points[j].Z, min, max)
		}
	}
}

// CenterCoordinates transforms all coordinates so the center of the model is at
// the origin point.
func (m *Mesh) CenterCoordinates() {
	vec3 := m.coordCenter()
	matrix := NewTranslationMatrix(vec3)
	for i := 0; i < len(m.Triangles); i++ {
		for j := 0; j < 3; j++ {
			transformed := matrix.MulVec4(m.Triangles[i].Points[j].Vec4()).Vec3()
			m.Triangles[i].Points[j] = transformed
		}
	}
}

// coordMinMax returns the minimum and maximum value for all vertex coordinates.
// This is useful for normalization.
func (m *Mesh) coordMinMax() (float64, float64) {
	var min float64 = math.MaxInt16
	var max float64 = math.MinInt16

	for _, t := range m.Triangles {
		for i := 0; i < 3; i++ {
			// Min
			if t.Points[i].X < min {
				min = t.Points[i].X
			}
			if t.Points[i].Y < min {
				min = t.Points[i].Y
			}
			if t.Points[i].Z < min {
				min = t.Points[i].Z
			}
			// Max
			if t.Points[i].X > max {
				max = t.Points[i].X
			}
			if t.Points[i].Y > max {
				max = t.Points[i].Y
			}
			if t.Points[i].Z > max {
				max = t.Points[i].Z
			}
		}
	}
	return min, max
}

// centerTranslation returns a translation vector that will center the mesh.
func (m *Mesh) coordCenter() Vec3 {
	var minx float64 = math.MaxInt16
	var maxx float64 = math.MinInt16
	var miny float64 = math.MaxInt16
	var maxy float64 = math.MinInt16
	var minz float64 = math.MaxInt16
	var maxz float64 = math.MinInt16

	for _, t := range m.Triangles {
		for i := 0; i < 3; i++ {
			// Min
			if t.Points[i].X < minx {
				minx = t.Points[i].X
			}
			if t.Points[i].Y < miny {
				miny = t.Points[i].Y
			}
			if t.Points[i].Z < minz {
				minz = t.Points[i].Z
			}
			// Max
			if t.Points[i].X > maxx {
				maxx = t.Points[i].X
			}
			if t.Points[i].Y > maxy {
				maxy = t.Points[i].Y
			}
			if t.Points[i].Z > maxz {
				maxz = t.Points[i].Z
			}
		}
	}

	// Not using the Y coord since FFT maps already sit on the floor. Adding
	// the Y translation would put the floor at the models 1/2 height point.
	x := -(maxx + minx) / 2.0
	y := 0.0 // -(maxy + miny) / 2.0
	z := -(maxz + minz) / 2.0

	return Vec3{x, y, z}
}
