package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS = 60
	// Number of milliseconds per frame
	TargetFrameTime = (1000 / FPS)
)

var (
	// Hacky solution to get current time on app start
	previous = sdl.GetTicks()

	// Temporary spot for vertices
	triangles = generateTriCube()

	// Allocate for all the points
	pointsToRender = make([]Vec2, len(triangles)*3)
)

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

func (e *Engine) Setup() {}

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
	// Target the specified FPS
	wait := TargetFrameTime - (sdl.GetTicks() - previous)
	if wait > 0 && wait <= TargetFrameTime {
		sdl.Delay(wait)
	}
	previous = sdl.GetTicks()

	// Clear the slice while retaining memory
	pointsToRender = pointsToRender[:0]

	// Increase the rotation
	e.rotation.y += 0.01
	e.rotation.x += 0.005

	// Project each into 2D
	for _, tri := range triangles {
		for _, point := range tri.points {
			transformedPoint := point
			// Rotate point on Y axis
			transformedPoint = transformedPoint.RotateY(e.rotation.y)
			transformedPoint = transformedPoint.RotateX(e.rotation.x)

			// Move the point away from the camera
			transformedPoint.z -= e.cameraPosition.z

			// Project the point
			projectedPoint := transformedPoint.Project()

			// Append the projected 2D point to the projected points
			pointsToRender = append(pointsToRender, projectedPoint)
		}
	}
}

func (e *Engine) Render() {
	e.renderer.Clear(ColorBlack)
	e.renderer.DrawGrid(ColorGrey)

	for _, point := range pointsToRender {
		// Move the point towards the center of the window.
		centeredX := point.x + (float64(e.window.width) / 2)
		centeredY := point.y + (float64(e.window.height) / 2)
		e.renderer.DrawRectangle(
			int(centeredX),
			int(centeredY),
			3,
			3,
			ColorYellow,
		)
	}

	// Render ColorBuffer
	e.window.Update(e.renderer.colorBuffer)
}
