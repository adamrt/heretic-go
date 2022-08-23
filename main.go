package main

const (
	WindowWidth  = 800
	WindowHeight = 800
)

func main() {
	window := NewWindow(WindowWidth, WindowHeight)
	defer window.Destroy()
	renderer := NewRenderer(WindowWidth, WindowHeight)
	engine := NewEngine(window, renderer)

	for engine.isRunning {
		engine.ProcessInput()
		engine.Update()
		engine.Render()
	}
}
