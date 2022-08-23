package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

var points = getCube()
var projectedPoints []Vec2

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{
		window:         window,
		renderer:       renderer,
		cameraPosition: Vec3{0, 0, -5},
		rotation:       Vec3{0, 0, 0},
		isRunning:      true,
	}
}

type Engine struct {
	window         *Window
	renderer       *Renderer
	cameraPosition Vec3
	rotation       Vec3
	isRunning      bool
}

func (e *Engine) Setup() {
	// Allocate for all the points
	projectedPoints = make([]Vec2, len(points))
}

func (e *Engine) ProcessInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.isRunning = false
			break
		case *sdl.KeyboardEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				e.isRunning = false
				break
			}
		}
	}
}

func (e *Engine) Update() {
	// Clear the slice while retaining memory
	projectedPoints = projectedPoints[:0]

	// Increase the rotation
	e.rotation.y += 0.001
	e.rotation.x += 0.0005

	// Project each point into 2D
	for _, point := range points {
		// Rotate point on Y axis
		transformedPoint := point
		transformedPoint = transformedPoint.RotateY(e.rotation.y)
		transformedPoint = transformedPoint.RotateX(e.rotation.x)

		// Move the point away from the camera
		transformedPoint.z -= e.cameraPosition.z

		// Project the point
		projectedPoint := transformedPoint.Project()

		// Append the projected 2D point to the projected points
		projectedPoints = append(projectedPoints, projectedPoint)
	}
}

func (e *Engine) Render() {
	e.renderer.Clear(ColorBlack)
	e.renderer.DrawGrid(ColorGrey)

	for _, point := range projectedPoints {
		// Move the point towards the center of the window.
		centeredX := point.x + (float64(e.window.width) / 2)
		centeredY := point.y + (float64(e.window.height) / 2)
		e.renderer.DrawRectangle(
			int(centeredX),
			int(centeredY),
			4,
			4,
			ColorYellow,
		)
	}

	// Render ColorBuffer
	e.window.Update(e.renderer.colorBuffer)
}
