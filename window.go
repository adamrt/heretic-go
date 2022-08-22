package main

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

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

	// I have no idea why we have to use ABGR8888. I would think it would be
	// RGBA8888 since that is the order of our `Color` struct.
	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(width),
		int32(height),
	)
	if err != nil {
		panic(err)
	}

	return &Window{
		width:     width,
		height:    height,
		window:    window,
		renderer:  renderer,
		texture:   texture,
		isRunning: true,

		colorBuffer: make([]Color, width*height),
	}
}

type Window struct {
	height, width int
	window        *sdl.Window
	renderer      *sdl.Renderer
	texture       *sdl.Texture

	isRunning bool

	colorBuffer []Color
}

func (w *Window) Setup() {
	w.clear()
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
	w.renderColorBuffer()
	w.clearColorBuffer(ColorYellow)
	w.renderer.Present()
}

func (w *Window) Destroy() {
	w.renderer.Destroy()
	w.window.Destroy()
	sdl.Quit()
}

func (w *Window) renderColorBuffer() {
	w.texture.Update(nil, unsafe.Pointer(&w.colorBuffer[0]), w.width*4)
	w.renderer.Copy(w.texture, nil, nil)
}

func (w *Window) clearColorBuffer(color Color) {
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			w.colorBuffer[y*w.width+x] = color
		}
	}
}

func (w *Window) clear() {
	w.renderer.SetDrawColor(0, 0, 0, 255)
	err := w.renderer.Clear()
	if err != nil {
		panic(err)
	}
}
