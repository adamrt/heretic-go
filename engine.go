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

const (
	CullModeNone     CullMode = 0
	CullModeBackFace CullMode = 1
)

type RenderMode int

const (
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
		IsRunning: true,

		ambientLight: DirectionalLight{Direction: Vec3{0, 0, 1}},

		cullMode:   CullModeBackFace,
		renderMode: RenderModeWireFill,
		scene:      NewScene(),

		// Rotation is set so if the user presses spacebar they get some
		// rotation. But autoRotation is off by default. Use
		// SetAutoRotation() to override the rotation value.
		autoRotation: false,
		rotation:     Vec3{0, 0.5, 0},
	}
}

type meshReader interface {
	ReadMesh(mapNum int) Mesh
}

type Engine struct {
	window   *Window
	renderer *Renderer

	IsRunning bool

	// Timing
	previous  uint32
	deltaTime float64

	// Rendering
	cullMode   CullMode
	renderMode RenderMode
	projMatrix Matrix
	camera     Camera
	frustrum   Frustrum

	// Model
	scene *scene

	MeshReader meshReader
	currentMap int

	ambientLight DirectionalLight

	// These two control the mesh rotating on its own.
	// The amount can be set by SetAutoRotation().
	autoRotation bool
	rotation     Vec3
}

func (e *Engine) Setup() {
	if len(e.scene.Meshes) == 0 {
		log.Fatalln("no mesh specified")
	}

	// If there is any texture on any mesh, show it.
	for _, mesh := range e.scene.Meshes {
		if len(mesh.Texture.data) != 0 {
			e.renderMode = RenderModeTexture
		}
	}
	// Projection matrix. We only need this calculate this once.
	aspectX := float64(e.window.width) / float64(e.window.height)
	aspectY := float64(e.window.height) / float64(e.window.width)
	fovY := math.Pi / 3.0 // Same as 180/3 or 60deg
	fovX := math.Atan(math.Tan(fovY/2.0)*aspectX) * 2.0
	znear := 0.1
	zfar := 100.0

	e.projMatrix = MatrixMakePerspective(fovY, aspectY, znear, zfar)
	e.frustrum = NewFrustrum(fovX, fovY, znear, zfar)
	e.camera = NewCamera(Vec3{0, 0.5, -1}, Vec3{0, 0, 0})

	e.previous = sdl.GetTicks()

}

func (e *Engine) ProcessInput() {
	state := sdl.GetKeyboardState()
	switch {
	case state[sdl.GetScancodeFromKey(sdl.K_w)] != 0:
		e.camera.MoveForward(e.deltaTime)
	case state[sdl.GetScancodeFromKey(sdl.K_s)] != 0:
		e.camera.MoveBackward(e.deltaTime)
	case state[sdl.GetScancodeFromKey(sdl.K_a)] != 0:
		e.camera.MoveLeft(e.deltaTime)
	case state[sdl.GetScancodeFromKey(sdl.K_d)] != 0:
		e.camera.MoveRight(e.deltaTime)
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			if event.GetType() != sdl.KEYDOWN {
				continue
			}
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				e.IsRunning = false
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
			case sdl.K_6:
				e.renderMode = RenderModeTextureWire
			case sdl.K_c:
				e.cullMode = CullModeNone
			case sdl.K_b:
				e.cullMode = CullModeBackFace
			case sdl.K_k:
				e.NextMap()
			case sdl.K_j:
				e.PrevMap()
			case sdl.K_SPACE:
				e.autoRotation = !e.autoRotation
			}
		case *sdl.QuitEvent:
			e.IsRunning = false
			break
		case *sdl.MouseWheelEvent:
			e.camera.MoveForward(float64(t.PreciseY) * e.deltaTime)
		case *sdl.MouseButtonEvent:
			if t.Button == sdl.BUTTON_RIGHT {
				e.camera.rightButtonPressed = t.Type == sdl.MOUSEBUTTONDOWN
			}
			if t.Button == sdl.BUTTON_LEFT {
				e.camera.leftButtonPressed = t.Type == sdl.MOUSEBUTTONDOWN
			}
		case *sdl.MouseMotionEvent:
			if e.camera.leftButtonPressed {
				e.camera.Look(t.XRel, t.YRel)
			}
			if e.camera.rightButtonPressed {
				e.camera.Pan(t.XRel, t.YRel)
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

	for _, mesh := range e.scene.Meshes {
		// Apply the engine's rotation vector. This is for automatic rotation.
		if e.autoRotation {
			mesh.Rotation = mesh.Rotation.Add(e.rotation.Mul(e.deltaTime))
		}

		// World matrix. Combination of scale, rotation and translation.
		worldMatrix := MatrixIdentity()
		worldMatrix = worldMatrix.Mul(NewScaleMatrix(mesh.Scale))
		worldMatrix = worldMatrix.Mul(NewRotationMatrix(mesh.Rotation))
		worldMatrix = worldMatrix.Mul(NewTranslationMatrix(mesh.Translation))

		// Setup Camera
		up := Vec3{0, 1, 0}
		target := e.camera.LookAtTarget()
		viewMatrix := e.camera.LookAtMatrix(target, up)

		// Project each into 2D
		for _, face := range mesh.Faces {
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
				// Why is this not the camera.position or
				// camera.direction?  Testing with the f22 gives
				// unexpected results, while Vec3{0,0,0} gives us the
				// expected results, but doesn't seem logical.
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
					lightIntensity: lightIntensity,
					texcoords:      tri.texcoords,
					palette:        face.palette,
					color:          face.color, // This is for filled triangles
				}

				for i, point := range tri.points {
					projected := e.projMatrix.MulVec4Proj(point)
					// FIXME: Invert Y to deal with obj coordinates
					// system.  I'd like to get rid of this but its
					// more complex than it seems. I think it has to
					// do with the handedness rules.
					projected.Y *= -1

					// Scale into view (tiny otherwise)
					projected.X *= (float64(e.window.width) / 2.0)
					projected.Y *= (float64(e.window.height) / 2.0)

					// Translate the projected points to the middle of the screen.
					projected.X += (float64(e.window.width) / 2.0)
					projected.Y += (float64(e.window.height) / 2.0)

					triangleToRender.points[i] = projected
				}

				mesh.trianglesToRender = append(mesh.trianglesToRender, triangleToRender)
			}
		}
	}
}

func (e *Engine) Render() {
	if e.scene.Background() != nil {
		e.renderer.ColorBufferBackground(*e.scene.Background())
	} else {
		e.renderer.ColorBufferColor(ColorBlack)
		e.renderer.DrawGrid(ColorGrey)
	}
	e.renderer.ZBufferClear()

	for _, mesh := range e.scene.Meshes {
		for _, tri := range mesh.trianglesToRender {
			a := tri.points[0]
			b := tri.points[1]
			c := tri.points[2]

			if e.renderMode == RenderModeTexture || e.renderMode == RenderModeTextureWire {
				if tri.HasTexture() {
					e.renderer.DrawTexturedTriangle(
						int(a.X), int(a.Y), a.Z, a.W, tri.texcoords[0],
						int(b.X), int(b.Y), b.Z, b.W, tri.texcoords[1],
						int(c.X), int(c.Y), c.Z, c.W, tri.texcoords[2],
						tri.lightIntensity,
						mesh.Texture,
						tri.palette,
					)
				} else {
					e.renderer.DrawFilledTriangle(
						int(a.X), int(a.Y), a.Z, a.W,
						int(b.X), int(b.Y), b.Z, b.W,
						int(c.X), int(c.Y), c.Z, c.W,
						tri.color,
					)
				}

			}
			if e.renderMode == RenderModeFill || e.renderMode == RenderModeWireFill {
				e.renderer.DrawFilledTriangle(
					int(a.X), int(a.Y), a.Z, a.W,
					int(b.X), int(b.Y), b.Z, b.W,
					int(c.X), int(c.Y), c.Z, c.W,
					tri.color)
			}

			if e.renderMode == RenderModeWire || e.renderMode == RenderModeWireVertex || e.renderMode == RenderModeWireFill || e.renderMode == RenderModeTextureWire {
				e.renderer.DrawTriangle(int(a.X), int(a.Y), int(b.X), int(b.Y), int(c.X), int(c.Y), ColorWhite)
			}

			if e.renderMode == RenderModeWireVertex {
				for _, point := range tri.points {
					e.renderer.DrawRectangle(int(point.X-2), int(point.Y-2), 4, 4, ColorRed)
				}
			}
		}
		// Clear the slice while retaining memory
		mesh.trianglesToRender = mesh.trianglesToRender[:0]
	}

	// Render ColorBuffer
	e.window.Update(e.renderer.colorBuffer)
}

func (e *Engine) SetMesh(mesh Mesh) {
	e.scene.Meshes = []*Mesh{&mesh}
}

func (e *Engine) AppendMesh(mesh Mesh) {
	e.scene.Meshes = append(e.scene.Meshes, &mesh)
}

func (e *Engine) SetAutoRotation(v Vec3) {
	e.rotation = v
	e.autoRotation = true
}

func (e *Engine) NextMap() {
	if e.currentMap < 125 {
		e.currentMap++
		mesh := e.MeshReader.ReadMesh(e.currentMap)
		e.SetMesh(mesh)
		e.Setup()
	}
}

func (e *Engine) PrevMap() {
	if e.currentMap > 1 {
		e.currentMap--
		mesh := e.MeshReader.ReadMesh(e.currentMap)
		e.SetMesh(mesh)
		e.Setup()
	}
}
