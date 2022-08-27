package main

// Triangle represents a triangle after rasterization.
type Triangle struct {
	points    [3]Vec4
	texcoords [3]Tex
	color     Color

	lightIntensity float64
}

func (t Triangle) Normal() Vec3 {
	a := t.points[0].Vec3()
	b := t.points[1].Vec3()
	c := t.points[2].Vec3()
	vectorAB := b.Sub(a).Normalize()
	vectorAC := c.Sub(a).Normalize()
	normal := vectorAB.Cross(vectorAC).Normalize() // Left handed system
	return normal
}
