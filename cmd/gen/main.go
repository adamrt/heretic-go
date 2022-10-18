package main

import (
	"time"

	"github.com/adamrt/heretic"
	"github.com/adamrt/heretic/fft"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"
)

// NewTorus creates a torus geometry with the specified revolution radius, tube radius,
// number of radial segments, number of tubular segments, and arc length angle in radians.
// TODO instead of 'arc' have thetaStart and thetaLength for consistency with other generators
// TODO then rename this to NewTorusSector and add a NewTorus constructor
func NewMap(mesh heretic.Mesh) *geometry.Geometry {

	t := geometry.NewGeometry()

	// Create buffers
	positions := math32.NewArrayF32(0, 0)
	normals := math32.NewArrayF32(0, 0)
	uvs := math32.NewArrayF32(0, 0)
	// indices := math32.NewArrayU32(0, 0)

	for _, t := range mesh.Triangles {
		for _, v := range t.Points {
			var vertex math32.Vector3
			vertex.X = float32(v.X)
			vertex.Y = float32(v.Y)
			vertex.Z = float32(v.Z)
			positions.AppendVector3(&vertex)
		}

		for _, tc := range t.Texcoords {
			var vertex math32.Vector2
			vertex.X = float32(tc.U)
			vertex.Y = float32(tc.V)
			uvs.AppendVector2(&vertex)
		}

		for _, n := range t.Normals {
			var vertex math32.Vector3
			vertex.X = float32(n.X)
			vertex.Y = float32(n.Y)
			vertex.Z = float32(n.Z)
			normals.AppendVector3(&vertex)
		}
	}

	t.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	t.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	t.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))
	return t
}

func main() {

	iso := fft.NewISOReader("/home/adam/tmp/emu/fft.iso")
	defer iso.Close()

	reader := fft.NewMeshReader(iso)

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager

	gui.Manager().Set(scene)

	width, height := a.GetSize()
	a.Gls().Viewport(0, 0, int32(width), int32(height))
	aspect := float32(width) / float32(height)

	// Create perspective camera
	// cam := camera.NewOrthographic(aspect, 0.1, 1000, 1.0, camera.Vertical)
	cam := camera.New(aspect)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	mapNum := 1
	count := 1
	startx := float32(-5.0)
	startz := float32(-5.0)
	for {
		m, err := reader.ReadMesh(mapNum)
		mapNum++
		if err != nil {
			continue
		}
		if len(m.Texture.Data) == 0 {
			continue
		}
		count++
		geom := NewMap(m)

		mat := material.NewStandard(math32.NewColor("White")) //
		tex := texture.NewTexture2DFromRGBA(m.Texture.RGBA())
		mat.AddTexture(tex)

		// m.Texture.WritePNG("texture2.png")

		mesh := graphic.NewMesh(geom, mat)
		mesh.SetPosition(startx, 0, startz)
		scene.Add(mesh)

		if mapNum > 50 {
			break
		}

		if count%5 == 0 {
			startx = -5.0
			startz += 1.0
		} else {
			startx += 1.0
		}
		// nb := material.NewStandard(math32.NewColor("Red"))
		// line := graphic.NewLines(geom, nb)
		// scene.Add(helper.NewNormals(line, 0.05, &math32.Color{1.0, 1.0, 1.0}, 0.50))
	}
	// mesh.SetScaleVec(&math32.Vector3{1.0, 1.0, -1.0})

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1.0, 1.0, 1.0}, 1.0)
	pointLight.SetPosition(1, 1, 1)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	// scene.Add(helper.NewAxes(0.5))
	// scene.Add(helper.NewGrid(3, 0.25, &math32.Color{0.3, 0.3, 0.3}))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)

	onResize("", nil)

	toggleWireframe := func(evname string, ev interface{}) {
		kev := ev.(*window.KeyEvent)
		if kev.Key == window.KeyW {
			// mat.SetWireframe(!mat.Wireframe())
		}
	}
	a.Subscribe(window.OnKeyDown, toggleWireframe)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
