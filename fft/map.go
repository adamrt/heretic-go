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

	// Normalize all coordinates to 0.0 - 1.0.
	min, max := minMaxTriangles(m.triangles)
	for i := 0; i < len(m.triangles); i++ {
		m.triangles[i] = normalizeTriangle(m.triangles[i], min, max)
	}

	// Center all coordinates so the center of the model is the model is
	// the origin point.
	translate := centerTraslation(m.triangles)
	tmx := heretic.NewTranslationMatrix(translate)
	for i := 0; i < len(m.triangles); i++ {
		for j := 0; j < 3; j++ {
			m.triangles[i].points[j] = tmx.MulVec4(m.triangles[i].points[j].Vec4()).Vec3()
		}
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
	fileHeader := make(meshFileHeader, meshFileHeaderLen)
	n, err := r.iso.file.Read(fileHeader)
	if err != nil || int64(n) != meshFileHeaderLen {
		log.Fatalf("read mesh file header: %v", err)
	}

	// Primary mesh pointer tells us where the primary mesh data is.  I
	// think this is always 196 as it starts directly after the header,
	// which has a size of 196. Keep dynamic here as pointer access will
	// grow and this keeps it consistent.
	primaryMeshPointer := fileHeader.PrimaryMesh()
	if primaryMeshPointer == 0 || primaryMeshPointer != 196 {
		log.Fatal("missing primary mesh pointer")
	}

	// Seek to the primary mesh data.
	r.iso.seekPointer(record.Sector(), primaryMeshPointer)

	// Mesh header contains the number of triangles and quads that exist.
	header := make(meshHeader, meshHeaderLen)
	n, err = r.iso.file.Read(header)
	if err != nil || int64(n) != meshHeaderLen {
		log.Fatalf("read mesh file header: %v", err)
	}

	triangles := make([]triangle, 0)
	for i := 0; i < header.N(); i++ {
		triangles = append(triangles, r.iso.readTriangle())
	}
	for i := 0; i < header.P(); i++ {
		triangles = append(triangles, r.iso.readQuad().split()...)
	}
	for i := 0; i < header.Q(); i++ {
		triangles = append(triangles, r.iso.readTriangle())
	}
	for i := 0; i < header.R(); i++ {
		triangles = append(triangles, r.iso.readQuad().split()...)
	}

	// Normals
	// Nothing is actually collected. They are just read here so the
	// iso read position moves forward, so we can read polygon texture data
	// next.  This could be cleaned up as a seek, but we may eventually use
	// the normal data here.
	for i := 0; i < header.N(); i++ {
		r.iso.readTriNormal()
	}
	for i := 0; i < header.P(); i++ {
		r.iso.readQuadNormal()
	}

	// Polygon texture data
	for i := 0; i < header.N(); i++ {
		triangles[i].textureData = r.iso.readTriUV()
	}
	for i := header.N(); i < header.TT(); i = i + 2 {
		textureData := r.iso.readQuadUV().split()
		triangles[i].textureData = textureData[0]
		triangles[i+1].textureData = textureData[1]
	}

	//
	// We skip an unknown chunk and the polygon tile locations for now
	//

	// Skip ahead to color palettes
	r.iso.seekPointer(record.Sector(), fileHeader.TexturePalettesColor())

	palettes := [16]*heretic.Palette{}
	for i := 0; i < 16; i++ {
		palette := &heretic.Palette{}
		for j := 0; j < 16; j++ {
			palette[j] = r.iso.readRGB15()
		}
		palettes[i] = palette
	}

	for i := 0; i < len(triangles); i++ {
		triangles[i].palette = palettes[triangles[i].textureData.palette]
	}

	// Skip ahead to lights
	r.iso.seekPointer(record.Sector(), fileHeader.LightsAndBackground())

	directionalLights := r.iso.readDirectionalLights()
	ambientLight := r.iso.readAmbientLight()
	background := r.iso.readBackground()

	return mesh{
		triangles:         triangles,
		directionalLights: directionalLights,
		ambientLight:      ambientLight,
		background:        background,
	}
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
