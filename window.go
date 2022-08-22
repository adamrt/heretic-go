package main

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var projectedPoints []Vec2

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
		sdl.WINDOW_BORDERLESS,
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

		colorBuffer: NewColorBuffer(width, height),
	}
}

type Window struct {
	height, width int
	window        *sdl.Window
	renderer      *sdl.Renderer
	texture       *sdl.Texture

	isRunning bool

	colorBuffer ColorBuffer
}

func (w *Window) Setup() {

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
	cube := getCube()
	for _, point := range cube {
		projectedPoints = append(projectedPoints, project(point))
	}
}

func (w *Window) Render() {
	w.clear(ColorBlack)
	w.drawGrid(ColorGrey)

	for _, point := range projectedPoints {
		w.drawRectangle(
			int(point.x+(float32(w.width)/2)),
			int(point.y+(float32(w.height)/2)),
			4,
			4,
			ColorYellow,
		)
	}

	// Render ColorBuffer
	w.texture.Update(nil, unsafe.Pointer(&w.colorBuffer.buf[0]), w.width*4)
	w.renderer.Copy(w.texture, nil, nil)

	w.renderer.Present()
}

func (w *Window) Destroy() {
	w.renderer.Destroy()
	w.window.Destroy()
	sdl.Quit()
}

func (w *Window) clear(color Color) {
	// Clear color buffer
	w.colorBuffer.Clear(ColorBlack)

	// Clear SDL renderer
	w.renderer.SetDrawColor(color.r, color.g, color.b, color.a)
	err := w.renderer.Clear()
	if err != nil {
		panic(err)
	}
}

func (w *Window) drawPixel(x, y int, color Color) {
	if x > 0 && x < int(w.width) && x > 0 && y < int(w.height) {
		w.colorBuffer.Set(x, y, color)
	}
}

func (w *Window) drawGrid(color Color) {
	for y := 0; y < w.height; y += 10 {
		for x := 0; x < w.width; x += 10 {
			w.drawPixel(x, y, color)
		}
	}
}

func (w *Window) drawRectangle(x, y, width, height int, color Color) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			current_x := x + i
			current_y := y + j
			w.drawPixel(current_x, current_y, color)
		}
	}
}
