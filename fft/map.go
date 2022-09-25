// This file contains a way to read a particular maps data (texture, mesh data, etc).
//
// It might make sense to rename MeshReader as GNSFileReader or MapReader since
// it more accurately represents what it does. The reason it is named MeshReader
// is because the primary output is an engine Mesh (heretic.Mesh).
package fft

import (
	"log"
	"math"

	"github.com/adamrt/heretic"
)

func NewMeshReader(iso ISOReader) MeshReader {
	return MeshReader{iso}
}

type mesh struct {
	triangles         []triangle
	directionalLights []heretic.DirectionalLight
	ambientLight      heretic.AmbientLight
	background        heretic.Background
}
type MeshReader struct {
	iso ISOReader
}

func (r MeshReader) ReadMesh(mapNum int) heretic.Mesh {
	records := r.readGNSRecords(mapNum)

	textures := []heretic.Texture{}
	m := mesh{}
	for _, record := range records {
		if record.Type() == RecordTypeTexture {
			texture := r.parseTexture(record)
			textures = append(textures, texture)
		} else if record.Type() == RecordTypeMeshPrimary {
			m = r.parsePrimaryMesh(record)
		}
	}

	// Normalize all coordinates to -1.0 - 1.0.
	min, max := minMaxTriangles(m.triangles)
	for i := 0; i < len(m.triangles); i++ {
		m.triangles[i] = normalizeTriangle(m.triangles[i], min, max)
	}

	// Convert fft Triangles to engine Faces
	faces := make([]heretic.Face, len(m.triangles))
	for i, tri := range m.triangles {
		faces[i] = tri.face()
	}

	mesh := heretic.Mesh{
		Faces:      faces,
		Texture:    textures[0],
		Background: &m.background,
		Scale:      heretic.Vec3{X: 1, Y: 1, Z: 1},
	}

	return mesh
}

func (r MeshReader) readGNSRecords(mapNum int) []GNSRecord {
	sector := GNSSectors[mapNum]
	r.iso.seekSector(sector)

	records := []GNSRecord{}
	for {
		record := make(GNSRecord, GNSRecordLen)
		n, err := r.iso.file.Read(record)
		if err != nil || n != GNSRecordLen {
			log.Fatalf("read gns record: %v", err)
		}
		if record.Type() == RecordTypeEnd {
			break
		}
		records = append(records, record)
	}
	return records
}

// parseTexture reads and returns an FFT texture as an engine Texture.
func (r MeshReader) parseTexture(record GNSRecord) heretic.Texture {
	r.iso.seekSector(record.Sector())
	data := make([]byte, record.Len())
	n, err := r.iso.file.Read(data)
	if err != nil || int64(n) != record.Len() {
		log.Fatalf("read texture data: %v", err)
	}
	pixels := textureSplitPixels(data)
	return heretic.NewTexture(textureWidth, textureHeight, pixels)
}

// parseTexture reads and returns an FFT mesh as an engine Mesh.
func (r MeshReader) parsePrimaryMesh(record GNSRecord) mesh {
	r.iso.seekSector(record.Sector())

	// File header contains intra-file pointers to areas of mesh data.
	meshFileHeader := make(meshFileHeader, meshFileHeaderLen)
	n, err := r.iso.file.Read(meshFileHeader)
	if err != nil || int64(n) != meshFileHeaderLen {
		log.Fatalf("read mesh file header: %v", err)
	}

	// Primary mesh pointer tells us where the primary mesh data is.  I
	// think this is always 196 as it starts directly after the header,
	// which has a size of 196. Keep dynamic here as pointer access will
	// grow and this keeps it consistent.
	primaryMeshPointer := meshFileHeader.PrimaryMesh()
	if primaryMeshPointer == 0 || primaryMeshPointer != 196 {
		log.Fatal("missing primary mesh pointer")
	}

	// Seek to the primary mesh data.
	r.iso.seekPointer(record.Sector(), primaryMeshPointer)

	// Mesh header contains the number of triangles and quads that exist.
	meshHeader := make(meshHeader, meshHeaderLen)
	n, err = r.iso.file.Read(meshHeader)
	if err != nil || int64(n) != meshHeaderLen {
		log.Fatalf("read mesh file header: %v", err)
	}

	triangles := make([]triangle, 0)
	for i := 0; i < meshHeader.numTexturedTriangles(); i++ {
		triangles = append(triangles, r.triangle())
	}
	for i := 0; i < meshHeader.numTexturedQuads(); i++ {
		triangles = append(triangles, r.quad().split()...)
	}
	for i := 0; i < meshHeader.numUntexturedTriangles(); i++ {
		triangles = append(triangles, r.triangle())
	}
	for i := 0; i < meshHeader.numUntexturedQuads(); i++ {
		triangles = append(triangles, r.quad().split()...)
	}

	// Normals
	// Nothing is actually collected. They are just read here so the
	// iso read position moves forward, so we can read polygon texture data
	// next.  This could be cleaned up as a seek, but we may eventually use
	// the normal data here.
	for i := 0; i < meshHeader.numTexturedTriangles(); i++ {
		r.triNormal()
	}
	for i := 0; i < meshHeader.numTexturedQuads(); i++ {
		r.quadNormal()
	}

	// Polygon texture data
	for i := 0; i < meshHeader.numTexturedTriangles(); i++ {
		triangles[i].textureData = r.triUV()
	}
	for i := meshHeader.numTexturedTriangles(); i < meshHeader.numTexturedTriangles()+(meshHeader.numTexturedQuads()*2); i = i + 2 {
		textureData := r.quadUV().split()
		triangles[i].textureData = textureData[0]
		triangles[i+1].textureData = textureData[1]
	}

	//
	// We skip an unknown chunk and the polygon tile locations for now
	//

	// Skip ahead to color palettes
	r.iso.seekPointer(record.Sector(), meshFileHeader.TexturePalettesColor())

	palettes := [16]*heretic.Palette{}
	for i := 0; i < 16; i++ {
		palette := &heretic.Palette{}
		for j := 0; j < 16; j++ {
			palette[j] = r.rgb15()
		}
		palettes[i] = palette
	}

	for i := 0; i < len(triangles); i++ {
		triangles[i].palette = palettes[triangles[i].textureData.palette]
	}

	// Skip ahead to lights
	r.iso.seekPointer(record.Sector(), meshFileHeader.LightsAndBackground())

	directionalLights := r.directionalLights()
	ambientLight := r.ambientLight()
	background := r.background()

	return mesh{
		triangles:         triangles,
		directionalLights: directionalLights,
		ambientLight:      ambientLight,
		background:        background,
	}
}

func (r MeshReader) rgb8() heretic.Color {
	return heretic.Color{
		R: r.iso.uint8(),
		G: r.iso.uint8(),
		B: r.iso.uint8(),
		A: 255,
	}
}

func (mr MeshReader) rgb15() heretic.Color {
	val := mr.iso.uint16()
	var a uint8
	if val == 0 {
		a = 0x00
	} else {
		a = 0xFF
	}

	b := uint8((val & 0b01111100_00000000) >> 7)
	g := uint8((val & 0b00000011_11100000) >> 2)
	r := uint8((val & 0b00000000_00011111) << 3)
	return heretic.Color{R: r, G: g, B: b, A: a}
}

func (r MeshReader) vertex() heretic.Vec3 {
	x := float64(r.iso.int16())
	y := float64(r.iso.int16())
	z := float64(r.iso.int16())
	return heretic.Vec3{x, -y, z}
}

func (r MeshReader) triangle() triangle {
	a := r.vertex()
	b := r.vertex()
	c := r.vertex()
	return triangle{points: [3]heretic.Vec3{a, b, c}}
}

func (r MeshReader) quad() quad {
	a := r.vertex()
	b := r.vertex()
	c := r.vertex()
	d := r.vertex()
	return quad{a, b, c, d}
}

func (r MeshReader) f1x3x12() float64 {
	return float64(r.iso.int16()) / 4096.0
}

func (r MeshReader) normal() normal {
	x := r.f1x3x12()
	y := r.f1x3x12()
	z := r.f1x3x12()
	return normal{x, -y, z}
}

func (r MeshReader) triNormal() []normal {
	a := r.normal()
	b := r.normal()
	c := r.normal()
	return []normal{a, b, c}
}

func (r MeshReader) quadNormal() []normal {
	a := r.normal()
	b := r.normal()
	c := r.normal()
	d := r.normal()
	return []normal{a, b, c, d}
}

func (r MeshReader) uv() uv {
	x := r.iso.uint8()
	y := r.iso.uint8()
	return uv{x: x, y: y}
}

func (r MeshReader) triUV() triangleTexData {
	a := r.uv()
	palette := int(r.iso.uint8() & 0b1111)
	r.iso.uint8() // padding
	b := r.uv()
	page := int(r.iso.uint8() & 0b11) // only 2 bits
	r.iso.uint8()                     // padding
	c := r.uv()
	return triangleTexData{a: a, b: b, c: c, palette: palette, page: page}
}

func (r MeshReader) quadUV() quadTexData {
	a := r.uv()
	palette := int(r.iso.uint8() & 0b1111)
	r.iso.uint8() // padding
	b := r.uv()
	page := int(r.iso.uint8() & 0b11) // only 2 bits
	r.iso.uint8()                     // padding
	c := r.uv()
	d := r.uv()
	return quadTexData{a: a, b: b, c: c, d: d, palette: palette, page: page}
}

func (r MeshReader) lightColor() uint8 {
	val := r.f1x3x12()
	return uint8(255 * math.Min(math.Max(0.0, val), 1.0))
}

func (r MeshReader) directionalLights() []heretic.DirectionalLight {
	l1r, l2r, l3r := r.lightColor(), r.lightColor(), r.lightColor()
	l1g, l2g, l3g := r.lightColor(), r.lightColor(), r.lightColor()
	l1b, l2b, l3b := r.lightColor(), r.lightColor(), r.lightColor()

	l1color := heretic.Color{R: l1r, G: l1g, B: l1b, A: 255}
	l2color := heretic.Color{R: l2r, G: l2g, B: l2b, A: 255}
	l3color := heretic.Color{R: l3r, G: l3g, B: l3b, A: 255}

	l1pos, l2pos, l3pos := r.vertex(), r.vertex(), r.vertex()

	return []heretic.DirectionalLight{
		{Position: l1pos, Color: l1color},
		{Position: l2pos, Color: l2color},
		{Position: l3pos, Color: l3color},
	}
}

func (r MeshReader) ambientLight() heretic.AmbientLight {
	color := r.rgb8()
	return heretic.AmbientLight{Color: color}

}

func (r MeshReader) background() heretic.Background {
	top := r.rgb8()
	bottom := r.rgb8()
	return heretic.Background{Top: top, Bottom: bottom}
}

func minMaxTriangles(triangles []triangle) (float64, float64) {
	var min float64 = math.MaxInt16
	var max float64 = math.MinInt16

	for _, t := range triangles {

		// Each point for max
		for i := 0; i < 3; i++ {
			// Max
			if t.points[i].X > max {
				max = t.points[i].X
			}
			if t.points[i].Y > max {
				max = t.points[i].Y
			}
			if t.points[i].Z > max {
				max = t.points[i].Z
			}

			// Min
			if t.points[i].X < min {
				min = t.points[i].X
			}
			if t.points[i].Y < min {
				min = t.points[i].Y
			}
			if t.points[i].Z < min {
				min = t.points[i].Z
			}
		}
	}
	return min, max
}
