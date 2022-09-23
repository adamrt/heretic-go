package fft

import "encoding/binary"

// Table of pointers contained in the meshFileHeader.
const (
	ptrPrimaryMesh          = 0x40
	ptrTexturePalettesColor = 0x44
	ptrUnknown              = 0x4c // Only non-zero in MAP000.5
	ptrLightsAndBackground  = 0x64 // Light colors/positions, bg gradient colors
	ptrTerrain              = 0x68 // Tile heights, slopes, and surface types
	ptrTextureAnimInst      = 0x6c
	ptrPaletteAnimInst      = 0x70
	ptrTexturePalettesGray  = 0x7c
	ptrMeshAnimInst         = 0x8c
	ptrAnimatedMesh1        = 0x90
	ptrAnimatedMesh2        = 0x94
	ptrAnimatedMesh3        = 0x98
	ptrAnimatedMesh4        = 0x9c
	ptrAnimatedMesh5        = 0xa0
	ptrAnimatedMesh6        = 0xa4
	ptrAnimatedMesh7        = 0xa8
	ptrAnimatedMesh8        = 0xac
	ptrVisibilityAngles     = 0xb0
)

// meshFileHeader contains 32-bit unsigned little-endian pointers to an area of
// the mesh data. Zero is returned if there is no pointer.
type meshFileHeader []byte

// meshFileHeaderLen is the length in bytes.
const meshFileHeaderLen = 196

// Return the intra-file pointer for different parts of the mesh data.
// All pointers are converted to int64 since thats what seek functions take
func (h meshFileHeader) ptr(location int32) int64 {
	const ptrLen = 4 // Intra-file pointers are always 32bit
	return int64(binary.LittleEndian.Uint32(h[location : location+ptrLen]))
}

func (h meshFileHeader) PrimaryMesh() int64          { return h.ptr(ptrPrimaryMesh) }
func (h meshFileHeader) TexturePalettesColor() int64 { return h.ptr(ptrTexturePalettesColor) }
func (h meshFileHeader) Unknown() int64              { return h.ptr(ptrUnknown) }
func (h meshFileHeader) LightsAndBackground() int64  { return h.ptr(ptrLightsAndBackground) }
func (h meshFileHeader) Terrain() int64              { return h.ptr(ptrTerrain) }
func (h meshFileHeader) TextureAnimInst() int64      { return h.ptr(ptrTextureAnimInst) }
func (h meshFileHeader) PaletteAnimInst() int64      { return h.ptr(ptrPaletteAnimInst) }
func (h meshFileHeader) TexturePalettesGray() int64  { return h.ptr(ptrTexturePalettesGray) }
func (h meshFileHeader) MeshAnimInst() int64         { return h.ptr(ptrMeshAnimInst) }
func (h meshFileHeader) AnimatedMesh1() int64        { return h.ptr(ptrAnimatedMesh1) }
func (h meshFileHeader) AnimatedMesh2() int64        { return h.ptr(ptrAnimatedMesh2) }
func (h meshFileHeader) AnimatedMesh3() int64        { return h.ptr(ptrAnimatedMesh3) }
func (h meshFileHeader) AnimatedMesh4() int64        { return h.ptr(ptrAnimatedMesh4) }
func (h meshFileHeader) AnimatedMesh5() int64        { return h.ptr(ptrAnimatedMesh5) }
func (h meshFileHeader) AnimatedMesh6() int64        { return h.ptr(ptrAnimatedMesh6) }
func (h meshFileHeader) AnimatedMesh7() int64        { return h.ptr(ptrAnimatedMesh7) }
func (h meshFileHeader) AnimatedMesh8() int64        { return h.ptr(ptrAnimatedMesh8) }
func (h meshFileHeader) VisibilityAngles() int64     { return h.ptr(ptrVisibilityAngles) }

// meshHeader contains the number of triangles and quads, textured and
// untextured. The numbers are represented by 16-bit unsigned integers.
type meshHeader []byte

// meshHeaderLen is the length in bytes.
const meshHeaderLen = 8

func (h meshHeader) numTexturedTriangles() int {
	return int(binary.LittleEndian.Uint16(h[0:2]))
}

func (h meshHeader) numTexturedQuads() int {
	return int(binary.LittleEndian.Uint16(h[2:4]))
}

func (h meshHeader) numUntexturedTriangles() int {
	return int(binary.LittleEndian.Uint16(h[4:6]))
}

func (h meshHeader) numUntexturedQuads() int {
	return int(binary.LittleEndian.Uint16(h[6:8]))
}
