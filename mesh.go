package heretic

import (
	"bufio"
	"errors"
	"fmt"
	"image/png"
	"io/fs"
	"log"
	"os"
	"strings"
)

type Mesh struct {
	// This is just a slice of a slice, but for naming purposes, triangles
	// makes more sense, since that is what it represents.
	Faces   []Face
	Texture Texture

	Background Background

	Rotation Vec3
	Scale    Vec3
	Trans    Vec3
}

func NewMesh(faces []Face, texture Texture) Mesh {
	return Mesh{Faces: faces, Texture: texture}
}

func NewMeshFromFile(objFilename string) *Mesh {
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
			matches, err := fmt.Fscanf(strings.NewReader(line), "vt %f %f", &vt.u, &vt.v)
			if err != nil || matches != 2 {
				log.Fatalf("vertex: only %d matches on line %q\n", matches, line)
			}
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
				mesh.Faces = append(mesh.Faces, Face{
					points: [3]Vec3{
						vertices[vertex_indices[0]-1],
						vertices[vertex_indices[1]-1],
						vertices[vertex_indices[2]-1],
					},
					color: ColorWhite,
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
				mesh.Faces = append(mesh.Faces, Face{
					points: [3]Vec3{
						vertices[vertex_indices[0]-1],
						vertices[vertex_indices[1]-1],
						vertices[vertex_indices[2]-1],
					},
					texcoords: [3]Tex{
						vts[texture_indices[0]-1],
						vts[texture_indices[1]-1],
						vts[texture_indices[2]-1],
					},
					color: ColorWhite,
				})
			}

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &mesh
}
