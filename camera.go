// This file contains a simple Arcball Camera that will rotate around the target.
//
// FIXME: when rotating to the left or right, then dragging up, we should take
//        into consideration the side we are looking at and make that side go up
//        instead of the front going up or down regardless of viewing angle.
// FIXME: Add mouse wheel movement
// FIXME: Add panning
//
package heretic

import "math"

type Camera struct {
	eye   Vec3
	front Vec3
	up    Vec3

	viewMatrix Matrix

	height, width float64
}

func NewCamera(eye, front, up Vec3, width, height int) *Camera {
	c := Camera{
		eye:    eye,
		front:  front,
		up:     up,
		height: float64(height),
		width:  float64(width),
	}
	return &c
}

func (c *Camera) ProcessMouseMovement(xrel, yrel, delta float64) {
	const EPS = 0.0001

	minPolarAngle := 0.0
	maxPolarAngle := math.Pi // 180 degrees as radians
	minAzimuthAngle := math.Inf(-1)
	maxAzimuthAngle := math.Inf(1)

	// Compute direction vector from target to camera
	tcam := c.eye
	tcam.Sub(c.front)

	// Calculate angles based on current camera position plus deltas
	radius := tcam.Length()
	theta := math.Atan2(tcam.X, tcam.Z) + (xrel * delta / 4)
	phi := math.Acos(tcam.Y/radius) + (-yrel * delta / 4)

	// Restrict phi and theta to be between desired limits
	phi = clamp(phi, minPolarAngle, maxPolarAngle)
	phi = clamp(phi, EPS, math.Pi-EPS)
	theta = clamp(theta, minAzimuthAngle, maxAzimuthAngle)

	// Calculate new cartesian coordinates
	tcam.X = radius * math.Sin(phi) * math.Sin(theta)
	tcam.Y = radius * math.Cos(phi)
	tcam.Z = radius * math.Sin(phi) * math.Cos(theta)

	// Update camera position and orientation
	c.eye = c.front.Add(tcam)
}
