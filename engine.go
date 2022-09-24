package heretic

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

	RenderModeWire        RenderMode = 1
	RenderModeWireVertex  RenderMode = 2
	RenderModeWireFill    RenderMode = 3
	RenderModeFill        RenderMode = 4
	RenderModeTexture     RenderMode = 5
	RenderModeTextureWire RenderMode = 6
)

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{
		window:    window,
		renderer:  renderer,
		isRunning: true,

		ambientLight: Light{direction: Vec3{0, 0, 1}},

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

	ambientLight Light
}

func (e *Engine) IsRunning() bool {
	return e.isRunning
}

func (e *Engine) Setup() {
	if e.mesh == nil {
		log.Fatalln("no mesh specified")
	}

	// If there is no texture, change the RenderMode to filled
	if len(e.mesh.Texture.data) == 0 {
		e.renderMode = RenderModeWireFill
	}
	// Projection matrix. We only need this calculate this once.
	aspectX := float64(e.window.width) / float64(e.window.height)
	aspectY := float64(e.window.height) / float64(e.window.width)
	fovY := math.Pi / 3.0 // Same as 180/3 or 60deg
	fovX := math.Atan(math.Tan(fovY/2.0)*aspectX) * 2.0
	znear := 1.0
	zfar := 1000.0

	e.projMatrix = MatrixMakePerspective(fovY, aspectY, znear, zfar)
	e.frustrum = NewFrustrum(fovX, fovY, znear, zfar)
	e.camera = NewCamera(Vec3{0, 0, 0}, Vec3{0, 0, 1})

	e.previous = sdl.GetTicks()

}

func (e *Engine) ProcessInput() {
	state := sdl.GetKeyboardState()
	//state[sdl.GetScancodeFromKey(sdl.K_UP)] != 0
	switch {
	case state[sdl.GetScancodeFromKey(sdl.K_ESCAPE)] != 0:
		e.isRunning = false
		break
	case state[sdl.GetScancodeFromKey(sdl.K_1)] != 0:
		e.renderMode = RenderModeWire
	case state[sdl.GetScancodeFromKey(sdl.K_2)] != 0:
		e.renderMode = RenderModeWireVertex
	case state[sdl.GetScancodeFromKey(sdl.K_3)] != 0:
		e.renderMode = RenderModeWireFill
	case state[sdl.GetScancodeFromKey(sdl.K_4)] != 0:
		e.renderMode = RenderModeFill
	case state[sdl.GetScancodeFromKey(sdl.K_5)] != 0:
		e.renderMode = RenderModeTexture
	case state[sdl.GetScancodeFromKey(sdl.K_6)] != 0:
		e.renderMode = RenderModeTextureWire
	case state[sdl.GetScancodeFromKey(sdl.K_c)] != 0:
		e.cullMode = CullModeNone
	case state[sdl.GetScancodeFromKey(sdl.K_b)] != 0:
		e.cullMode = CullModeBackFace
	case state[sdl.GetScancodeFromKey(sdl.K_w)] != 0:
		e.camera.velocity = e.camera.direction.Mul(15.0 * e.deltaTime)
		e.camera.position = e.camera.position.Add(e.camera.velocity)
	case state[sdl.GetScancodeFromKey(sdl.K_s)] != 0:
		e.camera.velocity = e.camera.direction.Mul(15.0 * e.deltaTime)
		e.camera.position = e.camera.position.Sub(e.camera.velocity)
	case state[sdl.GetScancodeFromKey(sdl.K_a)] != 0:
		e.camera.velocity = e.camera.right.Mul(15.0 * e.deltaTime)
		e.camera.position = e.camera.position.Add(e.camera.velocity)
	case state[sdl.GetScancodeFromKey(sdl.K_d)] != 0:
		e.camera.velocity = e.camera.right.Mul(15.0 * e.deltaTime)
		e.camera.position = e.camera.position.Sub(e.camera.velocity)
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.isRunning = false
			break
		case *sdl.MouseWheelEvent:
			e.camera.velocity = e.camera.direction.Mul(float64(t.PreciseY) * e.deltaTime * 15)
			e.camera.position = e.camera.position.Add(e.camera.velocity)
		case *sdl.MouseButtonEvent:
			if t.Button == sdl.BUTTON_RIGHT {
				e.camera.rightButtonPressed = t.Type == sdl.MOUSEBUTTONDOWN
			}
			if t.Button == sdl.BUTTON_LEFT {
				e.camera.leftButtonPressed = t.Type == sdl.MOUSEBUTTONDOWN
			}
		case *sdl.MouseMotionEvent:
			if e.camera.leftButtonPressed {
				e.camera.pitch += float64(t.YRel) / 200
				e.camera.yaw += float64(t.XRel) / 200
			}
			if e.camera.rightButtonPressed {
				e.camera.velocity = e.camera.right.Mul(float64(t.XRel) / 50.0)
				e.camera.position = e.camera.position.Add(e.camera.velocity)
				e.camera.velocity = e.camera.up.Mul(float64(t.YRel) / 50.0)
				e.camera.position = e.camera.position.Add(e.camera.velocity)
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
	// e.mesh.rotation.x -= 0.5 * e.deltaTime
	// e.mesh.rotation.z = math.Pi / 2
	// e.mesh.rotation.y += 0.5 * e.deltaTime
	// e.mesh.rotation.z += 0.3 * e.deltaTime

	// e.mesh.scale.x += 0.002 * e.deltaTime
	// e.mesh.scale.y += 0.001 * e.deltaTime

	// e.mesh.trans.x += 0.01
	// e.mesh.trans.z = 4.0 // constant

	// e.camera.position.x += 0.02 * e.deltaTime
	// e.camera.position.y += 0.01 * e.deltaTime
	// e.camera.position.z += 0.3 * e.deltaTime

	// World matrix. Combination of scale, rotation and translation
	worldMatrix := MatrixIdentity()
	scaleMatrix := MatrixMakeScale(e.mesh.Scale.x, e.mesh.Scale.y, e.mesh.Scale.z)
	rotXMatrix := MatrixMakeRotX(e.mesh.Rotation.x)
	rotYMatrix := MatrixMakeRotY(e.mesh.Rotation.y)
	rotZMatrix := MatrixMakeRotZ(e.mesh.Rotation.z)
	transMatrix := MatrixMakeTrans(e.mesh.Trans.x, e.mesh.Trans.y, e.mesh.Trans.z)

	worldMatrix = scaleMatrix.Mul(worldMatrix)
	worldMatrix = rotXMatrix.Mul(worldMatrix)
	worldMatrix = rotYMatrix.Mul(worldMatrix)
	worldMatrix = rotZMatrix.Mul(worldMatrix)
	worldMatrix = transMatrix.Mul(worldMatrix)

	// Camera
	up := Vec3{0, 1, 0}
	target := e.camera.LookAtTarget(Vec3{0, -0.2, 1})
	viewMatrix := e.camera.LookAtMatrix(target, up)

	// Project each into 2D
	for _, face := range e.mesh.Faces {
		// Transformation
		var transformedTri Triangle
		for i := 0; i < 3; i++ {
			transformedPoint := worldMatrix.MulVec4(face.points[i].Vec4())
			transformedPoint = viewMatrix.MulVec4(transformedPoint)
			transformedTri.points[i] = transformedPoint
		}

		normal := transformedTri.Normal()

		// Backface Culling
		//
		// 1. Find the vector between a point in the triangle and the camera origin.
		// 2. Determine the alignment of the ray and the normal
		if e.cullMode == CullModeBackFace {
			origin := Vec3{0, 0, 0}
			cameraRay := origin.Sub(transformedTri.points[0].Vec3())
			visibility := normal.Dot(cameraRay)
			if visibility < 0 {
				continue
			}
		}

		// Clip Polygons against the frustrum
		clippedTriangles := e.frustrum.Clip(transformedTri, face.texcoords)

		lightIntensity := -normal.Dot(e.ambientLight.Direction)

		// Projection
		for _, tri := range clippedTriangles {

			// The final triangle we will render
			triangleToRender := Triangle{
				// This is for filled triangles
				color: applyLightIntensity(face.color, lightIntensity),
				// This is for textured triangles
				lightIntensity: lightIntensity,
				texcoords:      tri.texcoords,
				palette:        face.palette,
			}

			for i, point := range tri.points {
				projected := e.projMatrix.MulVec4Proj(point)
				// FIXME: Invert Y to deal with obj coordinates system.
				projected.y *= -1

				// Scale into view (tiny otherwise)
				projected.x *= (float64(e.window.width) / 2.0)
				projected.y *= (float64(e.window.height) / 2.0)

				// Translate the projected points to the middle of the screen.
				projected.x += (float64(e.window.width) / 2.0)
				projected.y += (float64(e.window.height) / 2.0)

				triangleToRender.points[i] = projected
			}

			e.trianglesToRender = append(e.trianglesToRender, triangleToRender)
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

		if e.renderMode == RenderModeTexture || e.renderMode == RenderModeTextureWire {
			if tri.HasTexture() {
				e.renderer.DrawTexturedTriangle(
					int(tri.points[0].x), int(tri.points[0].y), tri.points[0].z, tri.points[0].w, tri.texcoords[0],
					int(tri.points[1].x), int(tri.points[1].y), tri.points[1].z, tri.points[1].w, tri.texcoords[1],
					int(tri.points[2].x), int(tri.points[2].y), tri.points[2].z, tri.points[2].w, tri.texcoords[2],
					tri.lightIntensity,
					e.mesh.Texture,
					tri.palette,
				)
			} else {
				e.renderer.DrawFilledTriangle(
					int(tri.points[0].x), int(tri.points[0].y), tri.points[0].z, tri.points[0].w,
					int(tri.points[1].x), int(tri.points[1].y), tri.points[1].z, tri.points[1].w,
					int(tri.points[2].x), int(tri.points[2].y), tri.points[2].z, tri.points[2].w,
					ColorBlack,
				)
			}

		}
		if e.renderMode == RenderModeFill || e.renderMode == RenderModeWireFill {
			e.renderer.DrawFilledTriangle(
				int(a.x), int(a.y), a.z, a.w,
				int(b.x), int(b.y), b.z, b.w,
				int(c.x), int(c.y), c.z, c.w,
				tri.color)
		}

		if e.renderMode == RenderModeWire || e.renderMode == RenderModeWireVertex || e.renderMode == RenderModeWireFill || e.renderMode == RenderModeTextureWire {
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
	e.mesh = NewMeshFromFile(objFile)
}

func (e *Engine) SetMesh(mesh Mesh) {
	e.mesh = &mesh
}
