package heretic

// Face represents a triangle before rasterization.

func NewFace(points [3]Vec3, color Color) Face {
	return Face{points: points, color: color}
}

type Face struct {
	points    [3]Vec3
	texcoords [3]Tex
	color     Color
}

func NewTex(u, v float64) Tex {
	return Tex{u, v}
}

type Tex struct {
	u, v float64
}
