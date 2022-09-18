package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
)

const (
	FFTTextureWidth  int = 256
	FFTTextureHeight int = 1024
	FFTTextureSize   int = FFTTextureWidth * FFTTextureHeight
	RawTextureSize   int = FFTTextureSize / 2
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

func (t Texture) WritePPM(filename string) {
	f, err := os.Create(fmt.Sprintf("./%s", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Buffered writer for performance
	bw := bufio.NewWriter(f)

	header := fmt.Sprintf("P3\n%d %d\n16\n", FFTTextureWidth, FFTTextureHeight)
	_, err = bw.WriteString(header)
	if err != nil {
		log.Fatal(err)
	}

	if len(t.data) != FFTTextureSize {
		log.Fatal("wrong size")
	}
	for _, pixel := range t.data {
		line := fmt.Sprintf("%d %d %d\n", pixel.r, pixel.g, pixel.b)
		_, err := bw.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}
	}
	bw.Flush()
}

func NewTextureFFT(raw []byte) Texture {
	data := splitPixels(raw)
	return Texture{
		width:  FFTTextureWidth,
		height: FFTTextureHeight,
		data:   data,
	}
}

// splitPixels takes the ISO's raw bytes and splits each of them into two
// bytes. The ISO has two pixels per byte to save space. We want each pixel
// independent, so we split them here.
func splitPixels(buf []byte) []Color {
	data := make([]Color, 0)
	for i := 0; i < RawTextureSize; i++ {
		colorA := uint8((buf[i] & 0xF0) >> 4)
		colorB := uint8(buf[i] & 0x0F)
		data = append(data,
			Color{colorA, colorA, colorA, 255},
			Color{colorB, colorB, colorB, 255},
		)

	}
	return data
}
