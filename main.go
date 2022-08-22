package main

func main() {
	window := NewWindow(800, 800)
	defer window.Destroy()

	window.Setup()

	for window.isRunning {
		window.ProcessInput()
		window.Update()
		window.Render()
	}
}
