// This file contains types and functions for clipping triangles against a
// frustrum.
package heretic

import (
	"math"
)

func NewFrustrum(fovX, fovY, znear, zfar float64) Frustrum {
	sinHalfFovX := math.Sin(fovX / 2.0)
	cosHalfFovX := math.Cos(fovX / 2.0)
	sinHalfFovY := math.Sin(fovY / 2.0)
	cosHalfFovY := math.Cos(fovY / 2.0)
	return Frustrum{
		planes: []Plane{
			// Left Plane
			{point: Vec3{0, 0, 0}, normal: Vec3{cosHalfFovX, 0, sinHalfFovX}},
			// Right Plane
			{point: Vec3{0, 0, 0}, normal: Vec3{-cosHalfFovX, 0, sinHalfFovX}},
			// Top Plane
			{point: Vec3{0, 0, 0}, normal: Vec3{0, -cosHalfFovY, sinHalfFovY}},
			// Bottom Plane
			{point: Vec3{0, 0, 0}, normal: Vec3{0, cosHalfFovY, sinHalfFovY}},
			// Near Plane
			{point: Vec3{0, 0, znear}, normal: Vec3{0, 0, 1}},
			// Far Plane
			{point: Vec3{0, 0, zfar}, normal: Vec3{0, 0, -1}},
		},
	}
}

// Frustrum is typically a 6 plane (front, back, right, left, top, bottom) geometry.
type Frustrum struct {
	planes []Plane
}

// Clip clips trianlges against each plane and returns 1 or more triangles.
func (f Frustrum) Clip(triangle Triangle) []Triangle {
	for _, plane := range f.planes {
		if len(triangle.Projected) > 0 {
			triangle = f.clipAgainstPlane(triangle, plane)
		}
	}
	return splitTriangle(triangle)
}

// clipAgainstPlane tries to clip a triangle against a plane. Instead of
// returning multiple triangles, it just returns one triangle with extra
// vertices, if clipped.
func (f Frustrum) clipAgainstPlane(triangle Triangle, plane Plane) Triangle {
	insideVertices := []Vec3{}
	insideTexcoords := []Tex{}

	previousVertex := triangle.Projected[len(triangle.Projected)-1].Vec3()
	previousTexcoord := triangle.Texcoords[len(triangle.Texcoords)-1]

	previousDot := previousVertex.Sub(plane.point).Dot(plane.normal)

	for i := 0; i < len(triangle.Projected); i++ {
		currentVertex := triangle.Projected[i].Vec3()
		currentTexcoord := triangle.Texcoords[i]

		currentDot := currentVertex.Sub(plane.point).Dot(plane.normal)

		if currentDot*previousDot < 0 {
			t := previousDot / (previousDot - currentDot)
			intersectionPoint := Vec3{
				X: lerp(previousVertex.X, currentVertex.X, t),
				Y: lerp(previousVertex.Y, currentVertex.Y, t),
				Z: lerp(previousVertex.Z, currentVertex.Z, t),
			}
			insideVertices = append(insideVertices, intersectionPoint)

			interpolatedTexcoord := Tex{
				U: lerp(previousTexcoord.U, currentTexcoord.U, t),
				V: lerp(previousTexcoord.V, currentTexcoord.V, t),
			}
			insideTexcoords = append(insideTexcoords, interpolatedTexcoord)
		}

		if currentDot > 0 {
			insideVertices = append(insideVertices, currentVertex)
			insideTexcoords = append(insideTexcoords, currentTexcoord)
		}

		previousDot = currentDot
		previousVertex = currentVertex
		previousTexcoord = currentTexcoord

	}

	// Convert these back to Vec4's
	projected := make([]Vec4, 0, len(insideVertices))
	for _, v := range insideVertices {
		projected = append(projected, v.Vec4())
	}

	// Update the original triangle so we retain the other fields of the
	// triangle (Palette, Color, etc).
	triangle.Projected = projected
	triangle.Texcoords = insideTexcoords

	return triangle
}

type Plane struct {
	point  Vec3
	normal Vec3
}

// splitTriangle splits a triangle into 1 or more triangles depending on how
// many projected points it has after being clipped.
func splitTriangle(triangle Triangle) []Triangle {
	triangles := []Triangle{}
	for i := 0; i < len(triangle.Projected)-2; i++ {
		index0 := 0
		index1 := i + 1
		index2 := i + 2

		// Copy the original so we retain the properties of the triangle.
		t := triangle

		t.Projected = []Vec4{
			triangle.Projected[index0],
			triangle.Projected[index1],
			triangle.Projected[index2],
		}
		t.Texcoords = []Tex{
			triangle.Texcoords[index0],
			triangle.Texcoords[index1],
			triangle.Texcoords[index2],
		}

		triangles = append(triangles, t)
	}
	return triangles
}
