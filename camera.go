package main

type Camera struct {
	position  Vec3
	direction Vec3
	velocity  Vec3
	yaw       float64
}

func NewCamera(position, direction Vec3) Camera {
	return Camera{
		position:  position,
		direction: direction,
	}
}
