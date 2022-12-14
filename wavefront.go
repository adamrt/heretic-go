// This file is for loading wavefront obj files as meshes.
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

func NewMeshFromObj(objFilename string) Mesh {
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
