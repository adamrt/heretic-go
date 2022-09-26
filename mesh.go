package heretic

import (
	"bufio"
	"errors"
	"fmt"
	"image/png"
	"io/fs"
	"log"
	"math"
	"os"
	"strings"
)

func NewMesh(triangles []Triangle, texture Texture) Mesh {
	return Mesh{Triangles: triangles, Texture: texture}
}

type Mesh struct {
	// This is just a slice of a slice, but for naming purposes, triangles
	// makes more sense, since that is what it represents.
	Triangles  []Triangle
	Texture    Texture
	Background *Background

	Rotation    Vec3
	Scale       Vec3
	Translation Vec3

	trianglesToRender []Triangle
}

// NormalizeCoordinates normalizes all vertex coordinates between 0 and 1. This
// scales down large models during import.  This is primary used for loading FFT
// maps since they have very large coordinates.  The min/max values should be
// the min and max of
func (m *Mesh) NormalizeCoordinates() {
	min, max := m.coordMinMax()
	for i := 0; i < len(m.Triangles); i++ {
		for j := 0; j < 3; j++ {
			m.Triangles[i].Points[j].X = normalize(m.Triangles[i].Points[j].X, min, max)
			m.Triangles[i].Points[j].Y = normalize(m.Triangles[i].Points[j].Y, min, max)
			m.Triangles[i].Points[j].Z = normalize(m.Triangles[i].Points[j].Z, min, max)
		}
	}
}

// CenterCoordinates transforms all coordinates so the center of the model is at
// the origin point.
func (m *Mesh) CenterCoordinates() {
	vec3 := m.coordCenter()
	matrix := NewTranslationMatrix(vec3)
	for i := 0; i < len(m.Triangles); i++ {
		for j := 0; j < 3; j++ {
			transformed := matrix.MulVec4(m.Triangles[i].Points[j].Vec4()).Vec3()
			m.Triangles[i].Points[j] = transformed
		}
	}
}

// coordMinMax returns the minimum and maximum value for all vertex coordinates.
// This is useful for normalization.
func (m *Mesh) coordMinMax() (float64, float64) {
	var min float64 = math.MaxInt16
	var max float64 = math.MinInt16

	for _, t := range m.Triangles {
		for i := 0; i < 3; i++ {
			// Min
			if t.Points[i].X < min {
				min = t.Points[i].X
			}
			if t.Points[i].Y < min {
				min = t.Points[i].Y
			}
			if t.Points[i].Z < min {
				min = t.Points[i].Z
			}
			// Max
			if t.Points[i].X > max {
				max = t.Points[i].X
			}
			if t.Points[i].Y > max {
				max = t.Points[i].Y
			}
			if t.Points[i].Z > max {
				max = t.Points[i].Z
			}
		}
	}
	return min, max
}

// centerTranslation returns a translation vector that will center the mesh.
func (m *Mesh) coordCenter() Vec3 {
	var minx float64 = math.MaxInt16
	var maxx float64 = math.MinInt16
	var miny float64 = math.MaxInt16
	var maxy float64 = math.MinInt16
	var minz float64 = math.MaxInt16
	var maxz float64 = math.MinInt16

	for _, t := range m.Triangles {
		for i := 0; i < 3; i++ {
			// Min
			if t.Points[i].X < minx {
				minx = t.Points[i].X
			}
			if t.Points[i].Y < miny {
				miny = t.Points[i].Y
			}
			if t.Points[i].Z < minz {
				minz = t.Points[i].Z
			}
			// Max
			if t.Points[i].X > maxx {
				maxx = t.Points[i].X
			}
			if t.Points[i].Y > maxy {
				maxy = t.Points[i].Y
			}
			if t.Points[i].Z > maxz {
				maxz = t.Points[i].Z
			}
		}
	}

	// Not using the Y coord since FFT maps already sit on the floor. Adding
	// the Y translation would put the floor at the models 1/2 height point.
	x := -(maxx + minx) / 2.0
	y := 0.0 // -(maxy + miny) / 2.0
	z := -(maxz + minz) / 2.0

	return Vec3{x, y, z}
}

func NewMeshFromFile(objFilename string) Mesh {
	objFile, err := os.Open(objFilename)
	if err != nil {
		panic(err)
	}
	defer objFile.Close()

	mesh := Mesh{
		Scale: Vec3{1.0, 1.0, 1.0},
	}

	pngFilename := strings.Split(objFilename, ".")[0] + ".png"
	pngFile, err := os.Open(pngFilename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Println(pngFilename, "does not exist.")
		} else {
			panic(err)
		}
	} else {
		defer pngFile.Close()

		image, err := png.Decode(pngFile)
		if err != nil {
			panic(err)
		}

		mesh.Texture = NewTextureFromImage(image)
	}

	var vertices []Vec3
	var vts []Tex

	scanner := bufio.NewScanner(objFile)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "v "):
			var v Vec3
			matches, err := fmt.Fscanf(strings.NewReader(line), "v %f %f %f", &v.X, &v.Y, &v.Z)
			if err != nil || matches != 3 {
				log.Fatalf("vertex: only %d matches on line %q\n", matches, line)
			}
			vertices = append(vertices, v)
		case strings.HasPrefix(line, "vt "):
			var vt Tex
			matches, err := fmt.Fscanf(strings.NewReader(line), "vt %f %f", &vt.U, &vt.V)
			if err != nil || matches != 2 {
				log.Fatalf("vertex: only %d matches on line %q\n", matches, line)
			}
			vt.V = 1 - vt.V
			vts = append(vts, vt)
		case strings.HasPrefix(line, "f "):
			var vertex_indices [3]int
			var texture_indices [3]int
			var normal_indices [3]int

			f := strings.NewReader(line)

			if !strings.Contains(line, "/") {
				matches, err := fmt.Fscanf(f, "f %d %d %d", &vertex_indices[0], &vertex_indices[1], &vertex_indices[2])
				if err != nil || matches != 3 {
					log.Fatalf("face: only %d matches on line %q\n", matches, line)
				}
				mesh.Triangles = append(mesh.Triangles, Triangle{
					Points: []Vec3{
						vertices[vertex_indices[0]-1],
						vertices[vertex_indices[1]-1],
						vertices[vertex_indices[2]-1],
					},
					Color: ColorWhite,
				})
			} else {
				matches, err := fmt.Fscanf(f, "f %d/%d/%d %d/%d/%d %d/%d/%d",
					&vertex_indices[0], &texture_indices[0], &normal_indices[0],
					&vertex_indices[1], &texture_indices[1], &normal_indices[1],
					&vertex_indices[2], &texture_indices[2], &normal_indices[2],
				)
				if err != nil || matches != 9 {
					log.Fatalf("face: only %d matches on line %q\n", matches, line)
				}
				mesh.Triangles = append(mesh.Triangles, Triangle{
					Points: []Vec3{
						vertices[vertex_indices[0]-1],
						vertices[vertex_indices[1]-1],
						vertices[vertex_indices[2]-1],
					},
					Texcoords: []Tex{
						vts[texture_indices[0]-1],
						vts[texture_indices[1]-1],
						vts[texture_indices[2]-1],
					},
					Color: ColorWhite,
				})
			}

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	mesh.NormalizeCoordinates()
	mesh.CenterCoordinates()

	return mesh
}
