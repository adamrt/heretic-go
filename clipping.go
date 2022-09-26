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

type Frustrum struct {
	planes []Plane
}

func (f Frustrum) Clip(triangle Triangle) []Triangle {
	polygon := Polygon{
		vertices: []Vec3{
			triangle.Projected[0].Vec3(),
			triangle.Projected[1].Vec3(),
			triangle.Projected[2].Vec3(),
		},
		texcoords: []Tex{
			triangle.Texcoords[0],
			triangle.Texcoords[1],
			triangle.Texcoords[2],
		},
	}
	for _, plane := range f.planes {
		if len(polygon.vertices) > 0 {
			polygon = f.clipAgainstPlane(polygon, plane)
		}
	}
	return polygon.AsTriangles()
}

func (f Frustrum) clipAgainstPlane(polygon Polygon, plane Plane) Polygon {
	insideVertices := []Vec3{}
	insideTexcoords := []Tex{}

	previousVertex := polygon.vertices[len(polygon.vertices)-1]
	previousTexcoord := polygon.texcoords[len(polygon.texcoords)-1]

	previousDot := previousVertex.Sub(plane.point).Dot(plane.normal)

	for i := 0; i < len(polygon.vertices); i++ {
		currentVertex := polygon.vertices[i]
		currentTexcoord := polygon.texcoords[i]

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
	return Polygon{insideVertices, insideTexcoords}
}

type Plane struct {
	point  Vec3
	normal Vec3
}

type Polygon struct {
	vertices  []Vec3
	texcoords []Tex
}

func (p Polygon) AsTriangles() []Triangle {
	tt := []Triangle{}
	for i := 0; i < len(p.vertices)-2; i++ {
		index0 := 0
		index1 := i + 1
		index2 := i + 2
		t := Triangle{
			Projected: []Vec4{
				p.vertices[index0].Vec4(),
				p.vertices[index1].Vec4(),
				p.vertices[index2].Vec4(),
			},
			Texcoords: []Tex{
				p.texcoords[index0],
				p.texcoords[index1],
				p.texcoords[index2],
			},
		}
		tt = append(tt, t)
	}
	return tt
}
