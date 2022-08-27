package main

import (
	"log"
	"math"

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
	RenderModeTexture    RenderMode = 5
)

var globalLight = Light{direction: Vec3{0, 0, 1}}
var buttonPressed = false

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{
		window:    window,
		renderer:  renderer,
		isRunning: true,

		cullMode:   CullModeBackFace,
		renderMode: RenderModeTexture,
	}
}

type Engine struct {
	window   *Window
	renderer *Renderer

	// Timing
	previous  uint32
	deltaTime float64

	isRunning bool

	// Rendering
	cullMode   CullMode
	renderMode RenderMode
	projMatrix Matrix
	camera     Camera
	frustrum   Frustrum

	// Model
	mesh              *Mesh
	trianglesToRender []Triangle
}

func (e *Engine) Setup() {
	if e.mesh == nil {
		log.Fatalln("no mesh specified")
	}

	// If there is no texture, change the RenderMode to filled
	if len(e.mesh.texture.data) == 0 {
		e.renderMode = RenderModeWireFill
	}

	// Projection matrix. We only need this calculate this once.
	aspectX := float64(WindowWidth) / float64(WindowHeight)
	aspectY := float64(WindowHeight) / float64(WindowWidth)
	fovY := math.Pi / 3.0 // Same as 180/3 or 60deg
	fovX := math.Atan(math.Tan(fovY/2.0)*aspectX) * 2.0
	znear := 1.0
	zfar := 20.0

	e.projMatrix = MatrixMakePerspective(fovY, aspectY, znear, zfar)
	e.frustrum = NewFrustrum(fovX, fovY, znear, zfar)
	e.camera = NewCamera(Vec3{0, 0, 0}, Vec3{0, 0, 1})

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
			case sdl.K_5:
				e.renderMode = RenderModeTexture
			case sdl.K_c:
				e.cullMode = CullModeNone
			case sdl.K_b:
				e.cullMode = CullModeBackFace
			case sdl.K_w:
				e.camera.velocity = e.camera.direction.Mul(10.0 * e.deltaTime * 2)
				e.camera.position = e.camera.position.Add(e.camera.velocity)
			case sdl.K_s:
				e.camera.velocity = e.camera.direction.Mul(10.0 * e.deltaTime)
				e.camera.position = e.camera.position.Sub(e.camera.velocity)
			case sdl.K_a:
				e.camera.yaw += 1.0 * e.deltaTime
			case sdl.K_d:
				e.camera.yaw -= 1.0 * e.deltaTime
			case sdl.K_q:
				e.camera.position.y += 3.0 * e.deltaTime
			case sdl.K_e:
				e.camera.position.y -= 3.0 * e.deltaTime

			}
		case *sdl.MouseWheelEvent:
			e.camera.velocity = e.camera.direction.Mul(float64(t.PreciseY) * e.deltaTime * 2)
			e.camera.position = e.camera.position.Add(e.camera.velocity)
		case *sdl.MouseButtonEvent:
			buttonPressed = t.Type == sdl.MOUSEBUTTONDOWN
		case *sdl.MouseMotionEvent:
			if buttonPressed {
				e.camera.yaw += float64(t.XRel) * e.deltaTime / 10
				e.camera.pitch += float64(t.YRel) * e.deltaTime / 10

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

	// Getting the deltaTime and multiplying the transformation keep
	// animation speed consistent regardless of FPS. It basically changes it
	// from tranforms per second instead of transforms per frame.
	e.deltaTime = float64(sdl.GetTicks()-e.previous) / 1000.0
	e.previous = sdl.GetTicks()

	// Increase the rotation/scale each frame
	e.mesh.rotation.x -= 0.5 * e.deltaTime
	e.mesh.rotation.z = math.Pi / 2
	// e.mesh.rotation.y += 0.5 * e.deltaTime
	// e.mesh.rotation.z += 0.3 * e.deltaTime

	// e.mesh.scale.x += 0.002 * e.deltaTime
	// e.mesh.scale.y += 0.001 * e.deltaTime

	// e.mesh.trans.x += 0.01
	e.mesh.trans.z = 4.0 // constant

	// e.camera.position.x += 0.02 * e.deltaTime
	// e.camera.position.y += 0.01 * e.deltaTime
	// e.camera.position.z += 0.3 * e.deltaTime

	// World matrix. Combination of scale, rotation and translation
	worldMatrix := MatrixIdentity()
	scaleMatrix := MatrixMakeScale(e.mesh.scale.x, e.mesh.scale.y, e.mesh.scale.z)
	rotXMatrix := MatrixMakeRotX(e.mesh.rotation.x)
	rotYMatrix := MatrixMakeRotY(e.mesh.rotation.y)
	rotZMatrix := MatrixMakeRotZ(e.mesh.rotation.z)
	transMatrix := MatrixMakeTrans(e.mesh.trans.x, e.mesh.trans.y, e.mesh.trans.z)

	worldMatrix = scaleMatrix.Mul(worldMatrix)
	worldMatrix = rotXMatrix.Mul(worldMatrix)
	worldMatrix = rotYMatrix.Mul(worldMatrix)
	worldMatrix = rotZMatrix.Mul(worldMatrix)
	worldMatrix = transMatrix.Mul(worldMatrix)

	// Camera
	up := Vec3{0, 1, 0}
	target := e.camera.LookAtTarget()
	viewMatrix := e.camera.LookAtMatrix(target, up)

	// Project each into 2D
	for _, face := range e.mesh.faces {
		//
		// Transformation
		//

		var transformedTri Triangle
		for i := 0; i < 3; i++ {
			transformedPoint := worldMatrix.MulVec4(face.points[i].Vec4())
			transformedPoint = viewMatrix.MulVec4(transformedPoint)
			transformedTri.points[i] = transformedPoint
		}

		//
		// Backface Culling
		//

		a := transformedTri.points[0].Vec3()
		b := transformedTri.points[1].Vec3()
		c := transformedTri.points[2].Vec3()
		vectorAB := b.Sub(a).Normalize()
		vectorAC := c.Sub(a).Normalize()
		normal := vectorAB.Cross(vectorAC).Normalize() // Left handed system

		// Find the vector between a point in the triangle and the camera origin
		origin := Vec3{0, 0, 0}
		cameraRay := origin.Sub(a)
		// Use dot product to determine the alignment of the camera ray and the normal
		visibility := normal.Dot(cameraRay)
		// Bypass triangles that are not facing the camera
		if e.cullMode == CullModeBackFace {
			if visibility < 0 {
				continue
			}
		}

		// Clip Polygons
		clippedTriangles := e.frustrum.Clip(transformedTri, face.texcoords)

		for _, tri := range clippedTriangles {

			//
			// Projection
			//

			var projectedTri Triangle
			for i, point := range tri.points {
				// Project the current vertex
				projectedPoint := e.projMatrix.MulVec4Proj(point)

				// Scale
				projectedPoint.x *= (float64(e.window.width) / 2.0)
				projectedPoint.y *= (float64(e.window.height) / 2.0)

				// Invert Y to deal with obj coordinates system.  FIXME:
				// I don't like this being here. I would rather it be
				// during obj parsing, but its not as simple as just
				// inverting Y (culling and lighting need to be inverted
				// as well).
				projectedPoint.y *= -1

				// Translate the projected points to the middle of the screen.
				projectedPoint.x += (float64(e.window.width) / 2.0)
				projectedPoint.y += (float64(e.window.height) / 2.0)

				// Append the projected 2D point to the projected points
				projectedTri.points[i] = projectedPoint
			}

			// Calculate shade intensity based on the normal and light direction
			lightIntensity := -normal.Dot(globalLight.direction)
			// Calculate color based on light
			projectedTri.color = applyLightIntensity(face.color, lightIntensity)
			projectedTri.lightIntensity = lightIntensity
			projectedTri.texcoords = tri.texcoords

			e.trianglesToRender = append(e.trianglesToRender, projectedTri)
		}
	}
}

func (e *Engine) Render() {
	e.renderer.ColorBufferClear(ColorBlack)
	e.renderer.ZBufferClear()
	e.renderer.DrawGrid(ColorGrey)

	for _, tri := range e.trianglesToRender {
		a := tri.points[0]
		b := tri.points[1]
		c := tri.points[2]

		if e.renderMode == RenderModeTexture {
			e.renderer.DrawTexturedTriangle(
				int(tri.points[0].x), int(tri.points[0].y), tri.points[0].z, tri.points[0].w, tri.texcoords[0],
				int(tri.points[1].x), int(tri.points[1].y), tri.points[1].z, tri.points[1].w, tri.texcoords[1],
				int(tri.points[2].x), int(tri.points[2].y), tri.points[2].z, tri.points[2].w, tri.texcoords[2],
				tri.lightIntensity,
				e.mesh.texture,
			)
		}
		if e.renderMode == RenderModeFill || e.renderMode == RenderModeWireFill {
			e.renderer.DrawFilledTriangle(
				int(a.x), int(a.y), a.z, a.w,
				int(b.x), int(b.y), b.z, b.w,
				int(c.x), int(c.y), c.z, c.w,
				tri.color)
		}

		if e.renderMode == RenderModeWire || e.renderMode == RenderModeWireVertex || e.renderMode == RenderModeWireFill {
			e.renderer.DrawTriangle(int(a.x), int(a.y), int(b.x), int(b.y), int(c.x), int(c.y), ColorWhite)
		}

		if e.renderMode == RenderModeWireVertex {
			for _, point := range tri.points {
				e.renderer.DrawRectangle(int(point.x-2), int(point.y-2), 4, 4, ColorRed)
			}
		}
	}

	// Clear the slice while retaining memory
	e.trianglesToRender = e.trianglesToRender[:0]

	// Render ColorBuffer
	e.window.Update(e.renderer.colorBuffer)
}

// LoadCubeMesh loads the cube geometry into the Engine.mesh
func (e *Engine) LoadMesh(objFile string) {
	// Temporary spot for vertices
	e.mesh = NewMesh(objFile)
}
