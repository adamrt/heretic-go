// This file contains a way to read binary data from the FFT ISO.
// It should be expanded to also read the FFT bin file.
//
// It contains the low level methods for different sized ints/uints as well has
// some simple geometry parsing. The higher level iso parsing happens in map.go.
// The split is somewhat arbitrary.
package fft

import (
	"encoding/binary"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/adamrt/heretic"
)

const sectorSize int64 = 2048

func NewISOReader(filename string) ISOReader {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("open iso: %v", err)
	}
	return ISOReader{f}
}

type ISOReader struct {
	file *os.File
}

func (r ISOReader) Close() {
	r.file.Close()
}

// seekSector will seek to the specified sector of the iso file.
func (r ISOReader) seekSector(sector int64) {
	to := sector * sectorSize
	_, err := r.file.Seek(to, 0)
	if err != nil {
		log.Fatalf("seek to sector: %v", err)
	}
}

// seekPointer will seek to the specified sector, plus a little more, of the iso
// file. This is useful when using MeshFileHeader intra-file pointers.
func (r ISOReader) seekPointer(sector int64, ptr int64) {
	to := sector*sectorSize + ptr
	_, err := r.file.Seek(to, 0)
	if err != nil {
		log.Fatalf("seek to pointer: %v", err)
	}
}

func (r ISOReader) readUint8() uint8 {
	size := 1
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return data[0]
}

func (r ISOReader) readUint16() uint16 {
	size := 2
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return binary.LittleEndian.Uint16(data)
}

func (r ISOReader) readUint32() uint32 {
	size := 4
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return binary.LittleEndian.Uint32(data)
}

func (r ISOReader) readInt8() int8   { return int8(r.readUint8()) }
func (r ISOReader) readInt16() int16 { return int16(r.readUint16()) }
func (r ISOReader) readInt32() int32 { return int32(r.readUint32()) }

func (r ISOReader) readRGB8() color.NRGBA {
	return color.NRGBA{
		R: r.readUint8(),
		G: r.readUint8(),
		B: r.readUint8(),
		A: 255,
	}
}

func (mr ISOReader) readRGB15() color.NRGBA {
	val := mr.readUint16()
	var a uint8
	if val == 0 {
		a = 0x00
	} else {
		a = 0xFF
	}

	b := uint8((val & 0b01111100_00000000) >> 7)
	g := uint8((val & 0b00000011_11100000) >> 2)
	r := uint8((val & 0b00000000_00011111) << 3)
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

func (r ISOReader) readVertex() heretic.Vec3 {
	x := float64(r.readInt16())
	y := float64(r.readInt16())
	z := float64(r.readInt16())
	return heretic.Vec3{X: x, Y: -y, Z: z}
}

func (r ISOReader) readTriangle() heretic.Triangle {
	a := r.readVertex()
	b := r.readVertex()
	c := r.readVertex()
	// This allocation is only currently necessary because of the
	// clipping.go code that indexes the len(coords)-1. We could do the
	// check over there and remove this. Not sure if worth it.
	// We also do this in quad.split().
	empty := make([]heretic.Tex, 3)
	return heretic.Triangle{Points: []heretic.Vec3{a, b, c}, Texcoords: empty}
}

func (r ISOReader) readQuad() quad {
	a := r.readVertex()
	b := r.readVertex()
	c := r.readVertex()
	d := r.readVertex()
	return quad{a, b, c, d}
}

func (r ISOReader) readF1x3x12() float64 {
	return float64(r.readInt16()) / 4096.0
}

func (r ISOReader) readNormal() heretic.Vec3 {
	x := r.readF1x3x12()
	y := r.readF1x3x12()
	z := r.readF1x3x12()
	return heretic.Vec3{X: x, Y: -y, Z: z}
}

func (r ISOReader) readTriNormal() []heretic.Vec3 {
	a := r.readNormal()
	b := r.readNormal()
	c := r.readNormal()
	return []heretic.Vec3{a, b, c}
}

func (r ISOReader) readQuadNormal() []heretic.Vec3 {
	a := r.readNormal()
	b := r.readNormal()
	c := r.readNormal()
	d := r.readNormal()
	return []heretic.Vec3{a, b, c, d}
}

func (r ISOReader) readUV() heretic.Tex {
	x := float64(r.readUint8())
	y := float64(r.readUint8())
	return heretic.Tex{U: x, V: y}
}

func (r ISOReader) readTriUV() textureData {
	a := r.readUV()
	palette := int(r.readUint8() & 0b1111)
	r.readUint8() // padding
	b := r.readUV()
	page := int(r.readUint8() & 0b11) // only 2 bits
	r.readUint8()                     // padding
	c := r.readUV()

	a = processTexCoords(a, page)
	b = processTexCoords(b, page)
	c = processTexCoords(c, page)

	return textureData{texCoords: []heretic.Tex{a, b, c}, palette: palette}
}

func (r ISOReader) readQuadUV() textureData {
	a := r.readUV()
	palette := int(r.readUint8() & 0b1111)
	r.readUint8() // padding
	b := r.readUV()
	page := int(r.readUint8() & 0b11) // only 2 bits
	r.readUint8()                     // padding
	c := r.readUV()
	d := r.readUV()

	a = processTexCoords(a, page)
	b = processTexCoords(b, page)
	c = processTexCoords(c, page)
	d = processTexCoords(d, page)

	return textureData{texCoords: []heretic.Tex{a, b, c, d}, palette: palette}
}

func (r ISOReader) readLightColor() uint8 {
	val := r.readF1x3x12()
	return uint8(255 * math.Min(math.Max(0.0, val), 1.0))
}

func (r ISOReader) readDirectionalLights() []heretic.DirectionalLight {
	l1r, l2r, l3r := r.readLightColor(), r.readLightColor(), r.readLightColor()
	l1g, l2g, l3g := r.readLightColor(), r.readLightColor(), r.readLightColor()
	l1b, l2b, l3b := r.readLightColor(), r.readLightColor(), r.readLightColor()

	l1color := color.NRGBA{R: l1r, G: l1g, B: l1b, A: 255}
	l2color := color.NRGBA{R: l2r, G: l2g, B: l2b, A: 255}
	l3color := color.NRGBA{R: l3r, G: l3g, B: l3b, A: 255}

	l1pos, l2pos, l3pos := r.readVertex(), r.readVertex(), r.readVertex()

	return []heretic.DirectionalLight{
		{Position: l1pos, Color: l1color},
		{Position: l2pos, Color: l2color},
		{Position: l3pos, Color: l3color},
	}
}

func (r ISOReader) readAmbientLight() heretic.AmbientLight {
	color := r.readRGB8()
	return heretic.AmbientLight{Color: color}

}

func (r ISOReader) readBackground() heretic.Background {
	top := r.readRGB8()
	bottom := r.readRGB8()
	return heretic.Background{Top: top, Bottom: bottom}
}

// processTexCoords has two functions:
//
// 1. Update the V coordinate to the specific page of the texture. FFT Textures
// have 4 pages (256x1024) and the original V specifies the pixel on one of the
// 4 pages. Multiple the page by the height of a single page (256).
//
// 2. Normalize the coordinates that can be U: 0-255 and V: 0-1023. Just divide
// them by their max to get a 0.0-1.0 value.
func processTexCoords(uv heretic.Tex, page int) heretic.Tex {
	v := float64(int(uv.V) + page*256)
	return heretic.Tex{U: uv.U / 255, V: v / 1023.0}
}
