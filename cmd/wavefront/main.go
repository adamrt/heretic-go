package main

import (
	"github.com/adamrt/heretic"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
)

func main() {
	renderer := heretic.NewRenderer(WindowWidth, WindowHeight)
	window := heretic.NewWindow(WindowWidth, WindowHeight)
	defer window.Destroy()

	engine := heretic.NewEngine(window, renderer)
	engine.LoadMesh("assets/drone.obj")

	engine.SetAutoRotation(heretic.Vec3{Y: 0.5})

	engine.Setup()
	for engine.IsRunning() {
		engine.ProcessInput()
		engine.Update()
		engine.Render()
	}
}
