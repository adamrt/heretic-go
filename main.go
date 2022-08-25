package main

const (
	WindowWidth  = 800
	WindowHeight = 800
)

func main() {
	renderer := NewRenderer(WindowWidth, WindowHeight)
	window := NewWindow(WindowWidth, WindowHeight)
	defer window.Destroy()

	engine := NewEngine(window, renderer)
	engine.LoadMesh("assets/f22.obj", "assets/f22.png")
	engine.Setup()

	for engine.isRunning {
		engine.ProcessInput()
		engine.Update()
		engine.Render()
	}
}
