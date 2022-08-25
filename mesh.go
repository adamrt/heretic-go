package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Mesh struct {
	// This is just a slice of a slice, but for naming purposes, triangles
	// makes more sense, since that is what it represents.
	faces []Face

	rotation Vec3
	scale    Vec3
	trans    Vec3
}

func NewMesh(filename string) *Mesh {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var mesh Mesh
	mesh.scale = Vec3{1.0, 1.0, 1.0}

	var vertices []Vec3

	scanner := bufio.NewScanner(f)
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
		case strings.HasPrefix(line, "f "):
			var f [3]int // face
			var t int    // trash
			matches, err := fmt.Fscanf(strings.NewReader(line), "f %d/%d/%d %d/%d/%d %d/%d/%d", &f[0], &t, &t, &f[1], &t, &t, &f[2], &t, &t)
			if err != nil || matches != 9 {
				log.Fatalf("face: only %d matches on line %q\n", matches, line)
			}
			a := vertices[f[0]-1]
			b := vertices[f[1]-1]
			c := vertices[f[2]-1]
			mesh.faces = append(mesh.faces, Face{points: [3]Vec3{a, b, c}, color: ColorWhite})
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &mesh
}
