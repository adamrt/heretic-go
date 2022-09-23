package main

import (
	"github.com/adamrt/heretic"
	"github.com/adamrt/heretic/fft"
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
	// engine.LoadMesh("assets/f22.obj")

	iso := fft.NewISOReader("/home/adam/tmp/emu/fft.iso")
	defer iso.Close()

	meshReader := fft.NewMapReader(iso)
	mesh := meshReader.ReadMesh(3)
	engine.SetMesh(mesh)

	engine.Setup()
	for engine.IsRunning() {
		engine.ProcessInput()
		engine.Update()
		engine.Render()
	}
}
