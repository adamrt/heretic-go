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

	trianglesToRender = make([]Triangle, len(triangles)*3)
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

	// Increase the rotation
	e.rotation.y += 0.01
	e.rotation.x += 0.005
	e.rotation.z += 0.0025

	// Project each into 2D
	for _, tri := range triangles {
		var projectedTriangle Triangle

		for i, point := range tri {
			transformedPoint := point
			// Rotate point on Y axis
			transformedPoint = transformedPoint.RotateX(e.rotation.x)
			transformedPoint = transformedPoint.RotateY(e.rotation.y)
			transformedPoint = transformedPoint.RotateZ(e.rotation.z)

			// Translate the vertex away from the camera
			transformedPoint.z -= e.cameraPosition.z

			// Project the vertex
			projectedPoint := ProjectPoint(transformedPoint)

			// Scale the projected point to the middle of the screen
			projectedPoint.x += (float64(e.window.width) / 2)
			projectedPoint.y += (float64(e.window.height) / 2)

			// Append the projected 2D point to the projected points
			projectedTriangle.points[i] = projectedPoint
		}
		trianglesToRender = append(trianglesToRender, projectedTriangle)
	}
}

func (e *Engine) Render() {
	e.renderer.Clear(ColorBlack)
	e.renderer.DrawGrid(ColorGrey)

	for _, tri := range trianglesToRender {
		e.renderer.DrawTriangle(tri, ColorYellow)
	}

	// Clear the slice while retaining memory
	trianglesToRender = trianglesToRender[:0]

	// Render ColorBuffer
	e.window.Update(e.renderer.colorBuffer)
}
