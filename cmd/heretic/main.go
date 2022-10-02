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
	fb := heretic.NewFrameBuffer(WindowWidth, WindowHeight)
	window := heretic.NewWindow(WindowWidth, WindowHeight)
	defer window.Destroy()

	engine := heretic.NewEngine(window, fb)
	// engine.LoadMesh("assets/f22.obj")

	iso := fft.NewISOReader("/home/adam/tmp/emu/fft.iso")
	defer iso.Close()

	engine.MeshReader = fft.NewMeshReader(iso)
	engine.NextMap()

	engine.Setup()
	for engine.IsRunning {
		engine.ProcessInput()
		engine.Update()
		engine.Render()
	}
}
