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

// meshReader is a temporary interface to avoid circular imports with the fft
// package. It will be removed once the project is better organized.
type meshReader interface {
	ReadMesh(mapNum int) Mesh
}

// Engine is the top level object that contains windows, renderers, etc.
// It has a basic game loop and in the typical for is run like so:
//
// engine.Setup()
// for engine.IsRunning() {
//   engine.ProcessInput()
//   engine.Update()
//   engine.Render()
// }
//
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
	camera     FPSCamera
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
	znear := 0.3
	zfar := 100.0

	e.projMatrix = MatrixMakePerspective(fovY, aspectY, znear, zfar)
	e.frustrum = NewFrustrum(fovX, fovY, znear, zfar)
	e.camera = NewFPSCamera(Vec3{0, 0.5, -1}, Vec3{0, 0, 0})

	e.previous = sdl.GetTicks()
}

func (e *Engine) ProcessInput() {
	// WASD keys are polled with keyboard state instead of polling for
	// events to create much smoother movement and to disregard key-repeat
	// functionality. Example: If we hold left, just smoothly move left.
	state := sdl.GetKeyboardState()
	e.camera.processKeyboardInput(state, e.deltaTime)

	// The other mouse/keyboard functionality can be handled via polling.
	// Moving the keys into keyboard state polling (above) will appear to be
	// multiple presses in a row when we just want one.
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.IsRunning = false
			break
		case *sdl.MouseWheelEvent:
			e.camera.processMouseWheel(float64(t.PreciseY), e.deltaTime)
		case *sdl.MouseButtonEvent:
			down := t.Type == sdl.MOUSEBUTTONDOWN
			if t.Button == sdl.BUTTON_RIGHT {
				e.camera.processMouseButton(MouseButtonRight, down)
			}
			if t.Button == sdl.BUTTON_LEFT {
				e.camera.processMouseButton(MouseButtonLeft, down)
			}
		case *sdl.MouseMotionEvent:
			e.camera.processMouseMovement(float64(t.XRel), float64(t.YRel), e.deltaTime)
		case *sdl.KeyboardEvent:
			if t.Type != sdl.KEYDOWN {
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
		}
	}
}

func (e *Engine) Update() {
	// Target the specified FPS
	wait := TargetFrameTime - (sdl.GetTicks() - e.previous)
	if wait > 0 && wait <= TargetFrameTime {
		sdl.Delay(wait)
	}

	// Using a deltaTime for transformations keeps animation speed
	// consistent regardless of FPS. It basically changes the engine
	// transforms-per-frame to tranforms-per-second.
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
		// I don't really like this. This should be computed in the camera itself.
		cameraRotation := MatrixIdentity().Mul(MatrixMakeRotY(e.camera.yaw)).Mul(MatrixMakeRotX(e.camera.pitch))

		// Start by looking down z-axis (left handed)
		target := Vec3{0, 0, 1}
		e.camera.front = cameraRotation.MulVec4(target.Vec4()).Vec3()
		target = e.camera.eye.Add(e.camera.front)
		viewMatrix := LookAt(e.camera.eye, target, e.camera.worldUp)

		// Project each into 2D
		for _, triangle := range mesh.Triangles {
			triangle.Projected = make([]Vec4, 3)

			// Transformation
			for i := 0; i < 3; i++ {
				transformedVertex := triangle.Points[i].Vec4()
				// Transform to world space
				transformedVertex = worldMatrix.MulVec4(transformedVertex)
				// Transform to view space
				transformedVertex = viewMatrix.MulVec4(transformedVertex)
				triangle.Projected[i] = transformedVertex
			}

			// Backface Culling
			//
			// 1. Find the vector between a point in the triangle and the camera origin.
			// 2. Determine the alignment of the ray and the normal
			//
			// Origin is always {0,0,0} since the camera is at in
			// view space. It shouldn't be eye/position.
			if e.cullMode == CullModeBackFace {
				origin := Vec3{0, 0, 0}
				cameraRay := origin.Sub(triangle.Projected[0].Vec3())
				visibility := triangle.Normal().Dot(cameraRay)
				if visibility < 0 {
					continue
				}
			}

			// Currently unused until we improve our lighting.
			// triangle.LightIntensity = -triangle.Normal().Dot(e.ambientLight.Direction)

			// Clip Polygons against the frustrum
			clippedTriangles := e.frustrum.Clip(triangle)

			// Projection
			for _, triangleToRender := range clippedTriangles {
				for i, point := range triangleToRender.Projected {
					projected := e.projMatrix.MulVec4Proj(point)
					// FIXME: Invert Y to deal with obj coordinates
					// system.  I'd like to get rid of this but its
					// more complex than it seems. I think it has to
					// do with the handedness rules.
					projected.Y *= -1

					// Scale into view (tiny otherwise)
					projected.X *= (float64(e.window.width) / 2.0)
					projected.Y *= (float64(e.window.height) / 2.0)

					// Translate the projected points to the
					// middle of the screen.  FIXME: If this
					// is removed, the viewport is the top
					// left only.  I understand the model
					// would be in top left, but I don't
					// understand why the viewport/frustrum
					// is changed.
					projected.X += (float64(e.window.width) / 2.0)
					projected.Y += (float64(e.window.height) / 2.0)

					triangleToRender.Projected[i] = projected
				}

				mesh.trianglesToRender = append(mesh.trianglesToRender, triangleToRender)
			}
		}
	}
}

func (e *Engine) Render() {
	// Draw a nice gradient background if we have one (typically from a FFT
	// Map) or fallback to just a black background with a grid.
	if e.scene.Background() != nil {
		e.renderer.colorBuffer.SetBackground(*e.scene.Background())
	} else {
		e.renderer.colorBuffer.Clear(ColorBlack)
		e.renderer.DrawGrid(ColorGrey)
	}

	e.renderer.zBuffer.Clear()

	for _, mesh := range e.scene.Meshes {
		for _, triangle := range mesh.trianglesToRender {
			a := triangle.Projected[0]
			b := triangle.Projected[1]
			c := triangle.Projected[2]

			//
			// Draw the triangles depending on the rendering mode.
			//

			if e.renderMode == RenderModeTexture || e.renderMode == RenderModeTextureWire {
				if triangle.HasTexture() {
					e.renderer.DrawTexturedTriangle(
						int(a.X), int(a.Y), a.Z, a.W, triangle.Texcoords[0],
						int(b.X), int(b.Y), b.Z, b.W, triangle.Texcoords[1],
						int(c.X), int(c.Y), c.Z, c.W, triangle.Texcoords[2],
						triangle.LightIntensity,
						triangle.Palette,
						mesh.Texture,
					)
				} else {
					e.renderer.DrawFilledTriangle(
						int(a.X), int(a.Y), a.Z, a.W,
						int(b.X), int(b.Y), b.Z, b.W,
						int(c.X), int(c.Y), c.Z, c.W,
						triangle.Color,
					)
				}

			}
			if e.renderMode == RenderModeFill || e.renderMode == RenderModeWireFill {
				e.renderer.DrawFilledTriangle(
					int(a.X), int(a.Y), a.Z, a.W,
					int(b.X), int(b.Y), b.Z, b.W,
					int(c.X), int(c.Y), c.Z, c.W,
					triangle.Color)
			}

			if e.renderMode == RenderModeWire || e.renderMode == RenderModeWireVertex || e.renderMode == RenderModeWireFill || e.renderMode == RenderModeTextureWire {
				e.renderer.DrawTriangle(int(a.X), int(a.Y), int(b.X), int(b.Y), int(c.X), int(c.Y), ColorWhite)
			}

			if e.renderMode == RenderModeWireVertex {
				for _, point := range triangle.Projected {
					e.renderer.DrawRectangle(int(point.X-2), int(point.Y-2), 4, 4, ColorRed)
				}
			}
		}

		// Clear the slice while retaining capacity so we don't have to
		// keep allocating each frame. The number of triangles can
		// change due to frustrum clipping and backface culling, but
		// keeping the capacity at the maximum seems reasonable.
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

// Move to the next FFT map. This is pretty hacky.
func (e *Engine) NextMap() {
	if e.currentMap < 125 {
		e.currentMap++
		mesh := e.MeshReader.ReadMesh(e.currentMap)
		e.SetMesh(mesh)
		e.Setup()
	}
}

// Move to the previous FFT map. This is pretty hacky.
func (e *Engine) PrevMap() {
	if e.currentMap > 1 {
		e.currentMap--
		mesh := e.MeshReader.ReadMesh(e.currentMap)
		e.SetMesh(mesh)
		e.Setup()
	}
}
