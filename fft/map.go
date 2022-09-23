// This file contains a way to read a particular maps data (texture, mesh data, etc).
//
// It might make sense to rename MeshReader as GNSFileReader or MapReader since
// it more accurately represents what it does. The reason it is named MeshReader
// is because the primary output is an engine Mesh (heretic.Mesh).
package fft

import (
	"log"

	"github.com/adamrt/heretic"
)

func NewMeshReader(iso ISOReader) MeshReader {
	return MeshReader{iso}
}

type MeshReader struct {
	iso ISOReader
}

func (r MeshReader) ReadMesh(mapNum int) heretic.Mesh {
	records := r.readGNSRecords(mapNum)

	textures := []heretic.Texture{}
	var mesh heretic.Mesh
	for _, record := range records {
		if record.Type() == RecordTypeTexture {
			texture := r.parseTexture(record)
			textures = append(textures, texture)
		} else if record.Type() == RecordTypeMeshPrimary {
			mesh = r.parsePrimaryMesh(record)
		}
	}

	mesh.SetScale(heretic.NewVec3(0.05, 0.05, 0.05))

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
	pixels := splitPixels(data)
	return heretic.NewTexture(textureWidth, textureHeight, pixels)
}

// parseTexture reads and returns an FFT mesh as an engine Mesh.
func (r MeshReader) parsePrimaryMesh(record GNSRecord) heretic.Mesh {
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

	// Convert fft Triangles to engine Faces
	faces := make([]heretic.Face, len(triangles))
	for i, tri := range triangles {
		faces[i] = heretic.NewFace(tri.points(), heretic.ColorWhite)
	}

	return heretic.NewMesh(faces)
}

func (r MeshReader) vertex() vertex {
	x := r.iso.int16()
	y := -r.iso.int16()
	z := r.iso.int16()
	return vertex{x, y, z}
}

func (r MeshReader) triangle() triangle {
	a := r.vertex()
	b := r.vertex()
	c := r.vertex()
	return triangle{a, b, c}
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
	a := r.f1x3x12()
	b := r.f1x3x12()
	c := r.f1x3x12()
	return normal{a, b, c}
}

func (r MeshReader) uv() polygonTexData {
	u := r.iso.uint8()
	v := r.iso.uint8()
	return polygonTexData{u: u, v: v}
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

func (r MeshReader) triUV() []polygonTexData {
	a := r.uv()
	palette := r.iso.uint8() & 0b1111
	r.iso.uint8() // padding
	b := r.uv()
	page := r.iso.uint8() & 0b11
	r.iso.uint8() // padding
	c := r.uv()

	colorPoint := (page << 4) + palette

	return []polygonTexData{{a.u, a.v, colorPoint}, {b.u, b.v, colorPoint}, {c.u, c.v, colorPoint}}
}

func (r MeshReader) quadUV() []polygonTexData {
	a := r.uv()
	palette := r.iso.uint8() & 0b1111
	r.iso.uint8() // padding
	b := r.uv()
	page := r.iso.uint8() & 0b11
	r.iso.uint8() // padding
	c := r.uv()
	d := r.uv()

	colorPoint := (page << 4) + palette

	return []polygonTexData{
		{a.u, a.v, colorPoint},
		{b.u, b.v, colorPoint},
		{c.u, c.v, colorPoint},
		{d.u, d.v, colorPoint},
	}
}
