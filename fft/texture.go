package fft

import "image/color"

const (
	textureWidth  int = 256
	textureHeight int = 1024
	textureLen    int = textureWidth * textureHeight
	textureRawLen int = textureLen / 2
)

// textureSplitPixels takes the ISO's raw bytes and splits each of them into two
// bytes. The ISO has two pixels per byte to save space. We want each pixel
// independent, so we split them here.
func textureSplitPixels(buf []byte) []color.NRGBA {
	data := make([]color.NRGBA, 0)
	for i := 0; i < textureRawLen; i++ {
		colorA := uint8(buf[i] & 0x0F)
		colorB := uint8((buf[i] & 0xF0) >> 4)

		// We dont care about RGB here.
		// This is just an index to the palette.
		data = append(data,
			color.NRGBA{R: colorA, G: colorA, B: colorA, A: 255},
			color.NRGBA{R: colorB, G: colorB, B: colorB, A: 255},
		)
	}
	return data
}
