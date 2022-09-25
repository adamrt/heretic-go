package fft

import (
	"github.com/adamrt/heretic"
)

const (
	textureWidth  int = 256
	textureHeight int = 1024
	textureLen    int = textureWidth * textureHeight
	textureRawLen int = textureLen / 2
)

// textureSplitPixels takes the ISO's raw bytes and splits each of them into two
// bytes. The ISO has two pixels per byte to save space. We want each pixel
// independent, so we split them here.
func textureSplitPixels(buf []byte) []heretic.Color {
	data := make([]heretic.Color, 0)
	for i := 0; i < textureRawLen; i++ {
		colorA := uint8(buf[i] & 0x0F)
		colorB := uint8((buf[i] & 0xF0) >> 4)

		// We dont care about RGB here.
		// This is just an index to the palette.
		data = append(data,
			heretic.Color{R: colorA, G: colorA, B: colorA, A: 255},
			heretic.Color{R: colorB, G: colorB, B: colorB, A: 255},
		)
	}
	return data
}
