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

func (r MeshReader) ReadMesh(mapNum int) (heretic.Mesh, error) {
	records := r.readGNSRecords(mapNum)

	textures := []heretic.Texture{}
	mesh := heretic.Mesh{}
	var err error
	for _, record := range records {
		if record.Type() == RecordTypeTexture {
			texture := r.parseTexture(record)
			textures = append(textures, texture)
		} else if record.Type() == RecordTypeMeshPrimary {
			mesh, err = r.parseMesh(record)
			if err != nil {
				return heretic.Mesh{}, err
			}
		} else if record.Type() == RecordTypeMeshAlt {
			// Sometimes there is no primary mesh (ie MAP002.GNS),
			// there is only an alternate. I'm not sure why. So we
			// treat this one as the primary, only if the primary
			// hasn't been set. Kinda Hacky until we start treating
			// each GNS Record as a Scenario.
			if len(mesh.Triangles) == 0 {
				mesh, err = r.parseMesh(record)
				if err != nil {
					return heretic.Mesh{}, err
				}
			}
		}
	}

	mesh.Scale = heretic.Vec3{X: 1, Y: 1, Z: 1}
	if len(textures) > 0 {
		mesh.Texture = textures[0]
	}

	mesh.NormalizeCoordinates()
	mesh.CenterCoordinates()
	return mesh, nil
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

// parseMesh reads mesh data for primary and alternate meshes.
//
func (r MeshReader) parseMesh(record GNSRecord) (heretic.Mesh, error) {
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

	// Previously we did these pointer checks on every map. But some maps
	// (ie MAP002.GNS) don't have a primary mesh, only alternative. The location of that mesh
	if record.Type() == RecordTypeMeshPrimary {
		if primaryMeshPointer == 0 || primaryMeshPointer != 196 {
			return heretic.Mesh{}, err
		}
	}

	// Skip ahead to color palettes
	r.iso.seekPointer(record.Sector(), fileHeader.TexturePalettesColor())

	palettes := make([]heretic.Palette, 16)
	for i := 0; i < 16; i++ {
		palette := make(heretic.Palette, 16)
		for j := 0; j < 16; j++ {
			palette[j] = r.iso.readRGB15()
		}
		palettes[i] = palette
	}

	// Seek to the mesh data.
	r.iso.seekPointer(record.Sector(), primaryMeshPointer)

	// Mesh header contains the number of triangles and quads that exist.
	header := make(meshHeader, meshHeaderLen)
	n, err = r.iso.file.Read(header)
	if err != nil || int64(n) != meshHeaderLen {
		log.Fatalf("read mesh file header: %v", err)
	}

	// FIXME: Change capacity from TT to total with untextured.
	triangles := make([]heretic.Triangle, 0, header.TT())
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
		triangles[i].Normals = r.iso.readTriNormal()
	}
	for i := header.N(); i < header.TT(); i = i + 2 {
		normals := splitNormals(r.iso.readQuadNormal())
		triangles[i].Normals = normals[0]
		triangles[i+1].Normals = normals[1]
	}
	for i := header.TT(); i < len(triangles); i++ {
		triangles[i].Normals = []heretic.Vec3{
			{0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0},
			{0.0, 0.0, 0.0},
		}
	}

	// Polygon texture data
	for i := 0; i < header.N(); i++ {
		uvData := r.iso.readTriUV()
		triangles[i].Texcoords = uvData.texCoords
		triangles[i].Palette = palettes[uvData.palette]
	}
	for i := header.N(); i < header.TT(); i = i + 2 {
		uvDatas := r.iso.readQuadUV().split()
		triangles[i].Texcoords = uvDatas[0].texCoords
		triangles[i].Palette = palettes[uvDatas[0].palette]

		triangles[i+1].Texcoords = uvDatas[1].texCoords
		triangles[i+1].Palette = palettes[uvDatas[1].palette]
	}

	// Skip ahead to lights
	r.iso.seekPointer(record.Sector(), fileHeader.LightsAndBackground())

	directionalLights := r.iso.readDirectionalLights()
	ambientLight := r.iso.readAmbientLight()
	background := r.iso.readBackground()

	return heretic.Mesh{
		Triangles:         triangles,
		Background:        &background,
		DirectionalLights: directionalLights,
		AmbientLight:      ambientLight,
	}, nil
}
