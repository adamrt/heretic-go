package main

import "github.com/veandco/go-sdl2/sdl"

func NewWindow(height, width int) *Window {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"Heretic",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(height),
		int32(width),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}

	return &Window{
		window:    window,
		renderer:  renderer,
		width:     width,
		height:    height,
		isRunning: true,
	}
}

type Window struct {
	height, width int
	window        *sdl.Window
	renderer      *sdl.Renderer

	isRunning bool
}

func (w *Window) ProcessInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			w.isRunning = false
			break
		case *sdl.KeyboardEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				w.isRunning = false
				break
			}
		}
	}
}

func (w *Window) Update() {
}

func (w *Window) Render() {
	w.clear()
	w.renderer.Present()
}

func (w *Window) Destroy() {
	w.renderer.Destroy()
	w.window.Destroy()
	sdl.Quit()
}

func (w *Window) clear() {
	w.renderer.SetDrawColor(255, 0, 0, 255)
	err := w.renderer.Clear()
	if err != nil {
		panic(err)
	}
}
