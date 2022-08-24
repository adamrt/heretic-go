package main

import (
	"log"
	"sort"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS = 60
	// Number of milliseconds per frame
	TargetFrameTime = (1000 / FPS)
)

type CullMode int
type RenderMode int

const (
	CullModeNone     CullMode = 0
	CullModeBackFace CullMode = 1

	RenderModeWire       RenderMode = 1
	RenderModeWireVertex RenderMode = 2
	RenderModeWireFill   RenderMode = 3
	RenderModeFill       RenderMode = 4
)

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{
		window:         window,
		renderer:       renderer,
		cameraPosition: Vec3{0, 0, 0},
		isRunning:      true,

		cullMode:   CullModeBackFace,
		renderMode: RenderModeWireFill,
	}
}

type Engine struct {
	window         *Window
	renderer       *Renderer
	cameraPosition Vec3

	// Timing
	previous  uint32
	isRunning bool

	// Rendering
	cullMode   CullMode
	renderMode RenderMode

	// Model
	mesh              *Mesh
	trianglesToRender []Triangle
}

func (e *Engine) Setup() {
	if e.mesh == nil {
		log.Fatalln("no mesh specified")
	}

	e.previous = sdl.GetTicks()

}

func (e *Engine) ProcessInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.isRunning = false
			break
		case *sdl.KeyboardEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				e.isRunning = false
				break
			case sdl.K_1:
				e.renderMode = RenderModeWire
			case sdl.K_2:
				e.renderMode = RenderModeWireVertex
			case sdl.K_3:
				e.renderMode = RenderModeWireFill
			case sdl.K_4:
				e.renderMode = RenderModeFill
			case sdl.K_c:
				e.cullMode = CullModeNone
			case sdl.K_b:
				e.cullMode = CullModeBackFace
			}
		}
	}
}

func (e *Engine) Update() {
	// Target the specified FPS
	wait := TargetFrameTime - (sdl.GetTicks() - e.previous)
	if wait > 0 && wait <= TargetFrameTime {
		sdl.Delay(wait)
	}
	e.previous = sdl.GetTicks()

	// Increase the rotation each frame
	e.mesh.rotation.x += 0.01
	e.mesh.rotation.y += 0.01
	e.mesh.rotation.z += 0.005

	// Project each into 2D
	for _, face := range e.mesh.faces {
		vertices := face.points
		transformedVertices := e.transform(vertices)

		if e.shouldCull(transformedVertices) {
			continue
		}

		projectedTri := e.project(transformedVertices)
		avgDepth := (transformedVertices[0].z + transformedVertices[1].z + transformedVertices[2].z) / 3.0
		projectedTri.averageDepth = avgDepth
		projectedTri.color = face.color

		e.trianglesToRender = append(e.trianglesToRender, projectedTri)
	}

	// Painters algorithm
	sort.Slice(e.trianglesToRender, func(i, j int) bool {
		a := e.trianglesToRender[i]
		b := e.trianglesToRender[j]
		return a.averageDepth > b.averageDepth
	})

}

func (e *Engine) Render() {
	e.renderer.Clear(ColorBlack)
	e.renderer.DrawGrid(ColorGrey)

	for _, tri := range e.trianglesToRender {
		a := tri.points[0]
		b := tri.points[1]
		c := tri.points[2]

		if e.renderMode == RenderModeFill || e.renderMode == RenderModeWireFill {
			e.renderer.DrawFilledTriangle(int(a.x), int(a.y), int(b.x), int(b.y), int(c.x), int(c.y), tri.color)
		}

		if e.renderMode == RenderModeWire || e.renderMode == RenderModeWireVertex || e.renderMode == RenderModeWireFill {
			e.renderer.DrawTriangle(int(a.x), int(a.y), int(b.x), int(b.y), int(c.x), int(c.y), ColorWhite)
		}

		if e.renderMode == RenderModeWireVertex {
			for _, point := range tri.points {
				e.renderer.DrawRectangle(int(point.x), int(point.y), 4, 4, ColorRed)
			}
		}
	}

	// Clear the slice while retaining memory
	e.trianglesToRender = e.trianglesToRender[:0]

	// Render ColorBuffer
	e.window.Update(e.renderer.colorBuffer)
}

func (e *Engine) transform(vertices [3]Vec3) [3]Vec3 {
	var transformedVertices [3]Vec3
	for i, point := range vertices {
		transformedPoint := point
		// Rotate
		transformedPoint = transformedPoint.RotateX(e.mesh.rotation.x)
		transformedPoint = transformedPoint.RotateY(e.mesh.rotation.y)
		transformedPoint = transformedPoint.RotateZ(e.mesh.rotation.z)

		// Translate (away from the camera)
		transformedPoint.z += 5

		transformedVertices[i] = transformedPoint
	}
	return transformedVertices
}

func (e *Engine) shouldCull(tri [3]Vec3) bool {
	if e.cullMode == CullModeNone {
		return false
	}

	a := tri[0]
	b := tri[1]
	c := tri[2]

	vectorAB := b.Sub(a)
	vectorAC := c.Sub(a)
	normal := vectorAB.Cross(vectorAC).Normalize() // Left handed system

	// Find the vector between a point in the triangle and the camera origin
	cameraRay := e.cameraPosition.Sub(a)

	// Use dot product to determine the alignment of the camera ray and the normal
	visibility := normal.Dot(cameraRay)

	// Bypass triangles that are not facing the camera
	return visibility < 0
}

func (e *Engine) project(vertices [3]Vec3) Triangle {
	var projectedTri Triangle
	for i, point := range vertices {
		projectedPoint := point.Project()

		// Scale the projected point to the middle of the screen
		projectedPoint.x += (float64(e.window.width) / 2)
		projectedPoint.y += (float64(e.window.height) / 2)

		// Append the projected 2D point to the projected points
		projectedTri.points[i] = projectedPoint
	}
	return projectedTri
}

// LoadCubeMesh loads the cube geometry into the Engine.mesh
func (e *Engine) LoadMesh(filename string) {
	// Temporary spot for vertices
	e.mesh = NewMesh(filename)
}

// LoadCubeMesh loads the cube geometry into the Engine.mesh
func (e *Engine) LoadCubeMesh() {
	// Temporary spot for vertices
	triangles := generateTriCube()
	e.mesh = &Mesh{faces: triangles}
	e.trianglesToRender = make([]Triangle, len(triangles)*3)
}
