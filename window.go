package heretic

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

func NewWindow(width, height int) *Window {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"Heretic",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(width),
		int32(height),
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
		width:    width,
		height:   height,
		window:   window,
		renderer: renderer,
		texture:  texture,
	}
}

type Window struct {
	height, width int
	window        *sdl.Window
	renderer      *sdl.Renderer
	texture       *sdl.Texture
}

// Update takes a color buffer, updates the SDL Texture, copies the texture into
// the SDL Renderer and then updates the screen.
func (w *Window) Update(frameBuffer FrameBuffer) {
	w.texture.Update(nil, unsafe.Pointer(&frameBuffer.cbuf[0]), w.pitch())
	w.renderer.Copy(w.texture, nil, nil)
	w.renderer.Present()
}

func (w *Window) Destroy() {
	w.renderer.Destroy()
	w.window.Destroy()
	sdl.Quit()
}

// pitch returns the byte length of a row (width * sizeof(uint32)).
func (w *Window) pitch() int { return w.width * 4 }
