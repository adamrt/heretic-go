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
		isRunning:      true,
	}
}

type Engine struct {
	window         *Window
	renderer       *Renderer
	cameraPosition Vec3
	isRunning      bool
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
	// Project each point into 2D
	for _, point := range points {
		// Move the points away from the camera
		point.z -= e.cameraPosition.z
		projectedPoints = append(projectedPoints, point.Project())
	}
}

func (e *Engine) Render() {
	e.renderer.Clear(ColorBlack)
	e.renderer.DrawGrid(ColorGrey)

	for _, point := range projectedPoints {
		e.renderer.DrawRectangle(
			int(point.x+(float32(e.window.width)/2)),
			int(point.y+(float32(e.window.height)/2)),
			4,
			4,
			ColorYellow,
		)
	}

	// Render ColorBuffer
	e.window.Update(e.renderer.colorBuffer)
}
