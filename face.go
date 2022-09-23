package heretic

// Face represents a triangle before rasterization.
type Face struct {
	points    [3]Vec3
	texcoords [3]Tex
	color     Color
}

type Tex struct {
	u, v float64
}
