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
	triangles := []triangle{}
	for _, record := range records {
		if record.Type() == RecordTypeTexture {
			texture := r.parseTexture(record)
			textures = append(textures, texture)
		} else if record.Type() == RecordTypeMeshPrimary {
			triangles = r.parsePrimaryMesh(record)
		}
	}

	// Convert fft Triangles to engine Faces
	faces := make([]heretic.Face, len(triangles))
	for i, tri := range triangles {
		faces[i] = tri.face()
	}

	mesh := heretic.NewMesh(faces, textures[0])
	mesh.SetScale(heretic.NewVec3(0.1, 0.1, 0.1))

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
	textureFlipVertical(pixels)
	return heretic.NewTexture(textureWidth, textureHeight, pixels)
}

// parseTexture reads and returns an FFT mesh as an engine Mesh.
func (r MeshReader) parsePrimaryMesh(record GNSRecord) []triangle {
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

	return triangles
}

func (r MeshReader) vertex() vertex {
	x := r.iso.int16()
	y := r.iso.int16()
	z := r.iso.int16()
	return vertex{x, -y, z}
}

func (r MeshReader) triangle() triangle {
	a := r.vertex()
	b := r.vertex()
	c := r.vertex()
	return triangle{a: a, b: b, c: c}
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
