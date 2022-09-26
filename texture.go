package heretic

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
)

// Tex is a texture coordinate
type Tex struct {
	U, V float64
}

func (t Tex) IsEmpty() bool {
	return t.U == 0.0 && t.V == 0.0
}

type Texture struct {
	width, height int
	data          []color.NRGBA
}

func NewTexture(width, height int, data []color.NRGBA) Texture {
	return Texture{width, height, data}
}

func NewTextureFromImage(image image.Image) Texture {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()
	data := make([]color.NRGBA, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := image.At(x, y).RGBA()
			color := color.NRGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}
			data[(y*width)+x] = color
		}
	}
	return Texture{width, height, data}
}

func (t Texture) WritePPM(filename string) {
	f, err := os.Create(fmt.Sprintf("./%s", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Buffered writer for performance
	bw := bufio.NewWriter(f)
	defer bw.Flush()

	// Write Header
	header := fmt.Sprintf("P3\n%d %d\n16\n", t.width, t.height)
	_, err = bw.WriteString(header)
	if err != nil {
		log.Fatal(err)
	}

	// Write pixel data
	for _, pixel := range t.data {
		line := fmt.Sprintf("%d %d %d\n", pixel.R, pixel.G, pixel.B)
		_, err := bw.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}
	}
}
