package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS = 60
	// Number of milliseconds per frame
	TargetFrameTime = (1000 / FPS)
)

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{
		window:         window,
		renderer:       renderer,
		cameraPosition: Vec3{0, 0, -5},
		isRunning:      true,
	}
}

type Engine struct {
	window         *Window
	renderer       *Renderer
	cameraPosition Vec3

	// Timing
	previous  uint32
	isRunning bool

	// Model
	mesh              *Mesh
	trianglesToRender []Triangle
}

func (e *Engine) Setup() {
	e.previous = sdl.GetTicks()
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
	// Target the specified FPS
	wait := TargetFrameTime - (sdl.GetTicks() - e.previous)
	if wait > 0 && wait <= TargetFrameTime {
		sdl.Delay(wait)
	}
	e.previous = sdl.GetTicks()

	// Increase the rotation each frame
	e.mesh.rotation.x += 0.01

	// Project each into 2D
	for _, tri := range e.mesh.triangles {
		var projectedTriangle Triangle

		for i, point := range tri {
			transformedPoint := point
			// Rotate point on X axis
			transformedPoint = transformedPoint.RotateX(e.mesh.rotation.x)

			// Translate the vertex away from the camera
			transformedPoint.z -= e.cameraPosition.z

			// Project the vertex (x/z, y/z)
			projectedPoint := transformedPoint.Project()

			// Scale the projected point to the middle of the screen
			projectedPoint.x += (float64(e.window.width) / 2)
			projectedPoint.y += (float64(e.window.height) / 2)

			// Append the projected 2D point to the projected points
			projectedTriangle.points[i] = projectedPoint
		}
		e.trianglesToRender = append(e.trianglesToRender, projectedTriangle)
	}
}

func (e *Engine) Render() {
	e.renderer.Clear(ColorBlack)
	e.renderer.DrawGrid(ColorGrey)

	for _, tri := range e.trianglesToRender {
		e.renderer.DrawTriangle(tri, ColorYellow)
	}

	// Clear the slice while retaining memory
	e.trianglesToRender = e.trianglesToRender[:0]

	// Render ColorBuffer
	e.window.Update(e.renderer.colorBuffer)
}

// LoadCubeMesh loads the cube geometry into the Engine.mesh
func (e *Engine) LoadMesh(filename string) {
	// Temporary spot for vertices
	e.mesh = NewMesh(filename)
}
