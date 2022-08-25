package main

import (
	"image"
)

type Texture struct {
	width, height int
	data          []Color
}

func NewTexture(image image.Image) Texture {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()
	data := make([]Color, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := image.At(x, y).RGBA()
			color := Color{
				r: uint8(r >> 8),
				g: uint8(g >> 8),
				b: uint8(b >> 8),
				a: uint8(a >> 8),
			}
			data[(y*width)+x] = color
		}
	}
	return Texture{
		width:  width,
		height: height,
		data:   data,
	}
}
