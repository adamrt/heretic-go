package main

func main() {
	window := NewWindow(640, 480)
	defer window.Destroy()

	window.Setup()

	for window.isRunning {
		window.ProcessInput()
		window.Update()
		window.Render()
	}
}
