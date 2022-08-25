package main

type Camera struct {
	position  Vec3
	direction Vec3
	up        Vec3
}

func NewCamera(position, direction, up Vec3) Camera {
	return Camera{position: position, direction: direction, up: up}
}
