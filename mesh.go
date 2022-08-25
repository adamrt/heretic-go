package main

import (
	"bufio"
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"
)

type Mesh struct {
	// This is just a slice of a slice, but for naming purposes, triangles
	// makes more sense, since that is what it represents.
	faces   []Face
	texture Texture

	rotation Vec3
	scale    Vec3
	trans    Vec3
}

func NewMesh(objFile, pngFile string) *Mesh {
	objF, err := os.Open(objFile)
	if err != nil {
		panic(err)
	}
	defer objF.Close()

	pngF, err := os.Open(pngFile)
	if err != nil {
		panic(err)
	}
	defer pngF.Close()

	image, err := png.Decode(pngF)
	if err != nil {
		panic(err)
	}

	mesh := Mesh{
		scale:   Vec3{1.0, 1.0, 1.0},
		texture: NewTexture(image),
	}

	var vertices []Vec3
	var vts []Tex

	scanner := bufio.NewScanner(objF)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "v "):
			var v Vec3
			matches, err := fmt.Fscanf(strings.NewReader(line), "v %f %f %f", &v.x, &v.y, &v.z)
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
			matches, err := fmt.Fscanf(f, "f %d/%d/%d %d/%d/%d %d/%d/%d",
				&vertex_indices[0], &texture_indices[0], &normal_indices[0],
				&vertex_indices[1], &texture_indices[1], &normal_indices[1],
				&vertex_indices[2], &texture_indices[2], &normal_indices[2],
			)
			if err != nil || matches != 9 {
				log.Fatalf("face: only %d matches on line %q\n", matches, line)
			}
			mesh.faces = append(mesh.faces, Face{
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

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &mesh
}
